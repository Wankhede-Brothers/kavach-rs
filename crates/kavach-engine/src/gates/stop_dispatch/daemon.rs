//! RPC daemon self-heal: launchd respawn (`spawn`) + bounded poll-retry probe.
mod spawn;

use spawn::try_spawn_rpc_daemon;

/// RPC ok+task -> `Ok(Some(json))` · ok+empty -> `Ok(None)` ·
/// transport error -> `Err(())` (caller fails closed via sentinel).
///
/// SELF-HEAL: on the first transport error, attempt a one-shot daemon respawn
/// and poll-retry (bounded 20×500ms) while `RocksDB` cold-opens, then honest
/// `Err(())` — the bounded escape valve for the fail-closed sentinel.
#[expect(
    clippy::print_stderr,
    reason = "hook engine has no tracing dep; stderr is the hook log channel"
)]
pub(super) fn rpc_next(method: &str, project_slug: &str) -> Result<Option<serde_json::Value>, ()> {
    let dbg = std::env::var("KAVACH_RPC_CLIENT_DEBUG").is_ok();
    let params = serde_json::json!({ "project": project_slug });
    let classify = |v: serde_json::Value| {
        if v.is_object() && v.get("key").is_some() {
            Ok(Some(v))
        } else {
            Ok(None)
        }
    };
    let first = kavach_rpc::client::call::<_, serde_json::Value>(method, Some(params.clone()));
    if dbg {
        match &first {
            Ok(v) => eprintln!("[rpc_next] {method}: first OK = {v}"),
            Err(e) => eprintln!("[rpc_next] {method}: first ERR = {e:?}"),
        }
    }
    if let Ok(v) = first {
        return classify(v);
    }
    // Daemon unreachable — on Unix self-heal: spawn it, then poll-retry while
    // RocksDB cold-opens (~2-4s, bounded 20×500ms). `try_spawn_rpc_daemon` is a
    // const `false` on non-Unix (the sync UDS daemon cannot run there), so the
    // retry below is skipped and we fall straight through to the direct fallback.
    let spawned = try_spawn_rpc_daemon();
    if dbg {
        eprintln!("[rpc_next] try_spawn_rpc_daemon -> {spawned}");
    }
    if spawned {
        for attempt in 0..20 {
            std::thread::sleep(std::time::Duration::from_millis(500));
            match kavach_rpc::client::call::<_, serde_json::Value>(method, Some(params.clone())) {
                Ok(v) => {
                    if dbg {
                        eprintln!("[rpc_next] retry {attempt} OK after spawn");
                    }
                    return classify(v);
                }
                Err(e) if dbg => eprintln!("[rpc_next] retry {attempt} ERR = {e:?}"),
                Err(_) => {}
            }
        }
    }
    // Still unreachable. On Unix a spawn/retry miss is a genuine daemon outage →
    // fail closed (a wrong "no work" answer is worse than an error on the
    // platform where the daemon is the single writer). On non-Unix the daemon
    // CANNOT run, so a client miss is the normal case: fall back to a direct
    // in-process DB query that reuses the daemon's own roadmap method. A failure
    // there still surfaces as the fail-closed SOURCE_DOWN block to the caller.
    direct_next_fallback(method, project_slug)
}

/// Unix: a daemon outage after spawn+retry is anomalous — fail closed so the
/// stop gate never reads a phantom "no work" from an unreachable source.
#[cfg(unix)]
fn direct_next_fallback(
    _method: &str,
    _project_slug: &str,
) -> Result<Option<serde_json::Value>, ()> {
    Err(())
}

/// Non-Unix direct-DB fallback: the sync UDS daemon never runs here, so open
/// `SurrealDB` locally and call the SAME roadmap RPC method the daemon would,
/// against a one-shot in-process `AppState`. Reuses the daemon's exact query
/// logic (no duplication). Errors collapse to `Err(())` (explicit `.ok()`
/// discard) so the caller fails closed for that single pass.
#[cfg(not(unix))]
fn direct_next_fallback(method: &str, project_slug: &str) -> Result<Option<serde_json::Value>, ()> {
    let (rt, state) = shared_db_state().ok_or(())?;
    rt.block_on(async {
        // `NextOpenTaskParams` is #[non_exhaustive] (cross-crate) — build it by
        // deserializing the same `{ "project": ... }` shape the RPC uses.
        let params: kavach_rpc::methods::roadmap::NextOpenTaskParams =
            serde_json::from_value(serde_json::json!({ "project": project_slug }))
                .ok()
                .ok_or(())?;
        let outcome = match method {
            "roadmap.next_open_task" => {
                kavach_rpc::methods::roadmap::next_open_task(state, params).await
            }
            "roadmap.next_open_hunt" => {
                kavach_rpc::methods::roadmap::next_open_hunt(state, params).await
            }
            "roadmap.promote_next_backlog" => {
                kavach_rpc::methods::roadmap::promote_next_backlog(state, params).await
            }
            _ => return Err(()),
        };
        match outcome {
            Ok(Some(card)) => serde_json::to_value(card).map(Some).map_err(|_| ()),
            Ok(None) => Ok(None),
            Err(_) => Err(()),
        }
    })
}

/// Process-wide cached `(Runtime, AppState)` for the direct-DB fallback. The
/// store is opened ONCE per process and reused across every `rpc_next` call.
/// This is critical on Windows: re-opening `RocksDB` per call self-contends,
/// because `SurrealDB` releases the single-writer LOCK *asynchronously* on drop,
/// so each new open races the previous handle's close. One persistent handle
/// (and one runtime to keep it alive) removes the per-call open entirely.
#[cfg(not(unix))]
#[expect(
    clippy::print_stderr,
    reason = "hook engine has no tracing dep; stderr is the hook log channel"
)]
fn shared_db_state() -> Option<&'static (tokio::runtime::Runtime, kavach_rpc::state::AppState)> {
    static CELL: std::sync::OnceLock<
        Option<(tokio::runtime::Runtime, kavach_rpc::state::AppState)>,
    > = std::sync::OnceLock::new();
    CELL.get_or_init(|| {
        let dbg = std::env::var("KAVACH_RPC_CLIENT_DEBUG").is_ok();
        let rt = tokio::runtime::Runtime::new().ok()?;
        let db = match rt.block_on(kavach_surreal::open_default_daemon()) {
            Ok(d) => d,
            Err(e) => {
                if dbg {
                    eprintln!("[direct_next] open ERR = {e}");
                }
                return None;
            }
        };
        Some((rt, kavach_rpc::state::AppState::new(db)))
    })
    .as_ref()
}
