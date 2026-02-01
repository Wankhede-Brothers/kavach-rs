//! STM (Short-Term Memory) updater: manages scratchpad and hot-context.
//! Reads current task state and updates STM files for session continuity.

use std::io::Write;
use std::path::PathBuf;

use crate::session;

pub fn run() -> anyhow::Result<()> {
    let sess = session::get_or_create_session();
    let out = std::io::stdout();
    let mut w = out.lock();

    let stm_dir = stm_base_dir();
    let project_dir = stm_dir.join("projects").join(&sess.project);

    if !project_dir.exists() {
        std::fs::create_dir_all(&project_dir)?;
    }

    // Update scratchpad with current task state
    let scratchpad_path = project_dir.join("scratchpad.toon");
    update_scratchpad(&scratchpad_path, &sess)?;

    // Update hot-context with modified files
    let hot_ctx_path = project_dir.join("hot-context.json");
    update_hot_context(&hot_ctx_path, &sess)?;

    writeln!(w, "[STM:UPDATED]")?;
    writeln!(w, "project: {}", sess.project)?;
    writeln!(w, "scratchpad: {}", scratchpad_path.display())?;
    writeln!(w, "hot_context: {}", hot_ctx_path.display())?;
    writeln!(w, "files_tracked: {}", sess.files_modified.len())?;

    Ok(())
}

fn stm_base_dir() -> PathBuf {
    let home = dirs::home_dir().unwrap_or_default();
    home.join(".local")
        .join("shared")
        .join("shared-ai")
        .join("stm")
}

fn update_scratchpad(path: &std::path::Path, sess: &session::SessionState) -> anyhow::Result<()> {
    let mut f = std::fs::File::create(path)?;
    writeln!(f, "# Scratchpad - {}", sess.today)?;
    writeln!(f, "session: {}", sess.id)?;
    writeln!(f, "project: {}", sess.project)?;
    writeln!(f)?;
    writeln!(
        f,
        "intent: {}",
        if sess.current_task.is_empty() {
            "null"
        } else {
            &sess.current_task
        }
    )?;
    writeln!(
        f,
        "status: {}",
        if sess.task_status.is_empty() {
            "idle"
        } else {
            &sess.task_status
        }
    )?;
    writeln!(f, "turn: {}", sess.turn_count)?;
    writeln!(f, "research: {}", sess.research_done)?;
    Ok(())
}

fn update_hot_context(path: &std::path::Path, sess: &session::SessionState) -> anyhow::Result<()> {
    let files: Vec<serde_json::Value> = sess
        .files_modified
        .iter()
        .map(|f| serde_json::json!({"path": f}))
        .collect();
    let json = serde_json::json!({
        "session": sess.id,
        "project": sess.project,
        "date": sess.today,
        "files": files,
    });
    let content = serde_json::to_string_pretty(&json)?;
    std::fs::write(path, content)?;
    Ok(())
}
