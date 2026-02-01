//! Session resume: restore state from Memory Bank after compaction.
//! DACE: ~80 tokens output, pointers to commands not data.

use std::io::Write;
use std::path::PathBuf;

use crate::session;

pub fn run() -> anyhow::Result<()> {
    let mut sess = session::get_or_create_session();
    let out = std::io::stdout();
    let mut w = out.lock();

    let was_post_compact = sess.is_post_compact();
    if was_post_compact {
        sess.clear_post_compact();
    }

    writeln!(w, "[RESUME:DACE]")?;
    writeln!(w, "date: {}", sess.today)?;
    writeln!(w, "session: {}", sess.id)?;
    writeln!(w, "project: {}", sess.project)?;
    if was_post_compact {
        writeln!(w, "compact_recovered: true")?;
    }
    writeln!(w)?;

    writeln!(w, "[STATE]")?;
    writeln!(
        w,
        "research_done: {} | memory: {} | ceo: {}",
        bool_str(sess.research_done),
        bool_str(sess.memory_queried),
        bool_str(sess.ceo_invoked)
    )?;
    writeln!(w)?;

    writeln!(w, "[ENFORCE]")?;
    writeln!(
        w,
        "TABULA_RASA: cutoff={} | WebSearch BEFORE code",
        sess.training_cutoff
    )?;
    writeln!(w, "NO_AMNESIA: query memory bank")?;
    writeln!(w)?;

    // Check for task to continue
    let mut task_found = false;
    if !sess.current_task.is_empty() {
        writeln!(
            w,
            "[TASK] {} | status: {}",
            sess.current_task, sess.task_status
        )?;
        writeln!(w)?;
        task_found = true;
    } else {
        // Try loading from scratchpad
        if let Some((intent, status)) = load_scratchpad_task(&sess.project) {
            writeln!(w, "[TASK] {intent} | status: {status}")?;
            writeln!(w)?;
            task_found = true;
        }
    }

    if !task_found {
        writeln!(w, "[TASK] none | Ask user for next task")?;
        writeln!(w)?;
    }

    // DACE: pointer only
    let mem_dir = memory_dir();
    let total = count_toon_recursive(&mem_dir);
    writeln!(w, "[MEMORY] {total} docs | query: kavach memory bank")?;

    sess.mark_memory_queried();

    Ok(())
}

fn load_scratchpad_task(project: &str) -> Option<(String, String)> {
    if project.is_empty() {
        return None;
    }

    let stm = stm_dir();
    let path = stm.join("projects").join(project).join("scratchpad.toon");
    let content = std::fs::read_to_string(&path).ok()?;

    let mut intent = String::new();
    let mut status = String::new();

    for line in content.lines() {
        let line = line.trim();
        if let Some(rest) = line.strip_prefix("intent:") {
            intent = rest.trim().to_string();
        } else if let Some(rest) = line.strip_prefix("status:") {
            status = rest.trim().to_string();
        }
    }

    if !intent.is_empty() && intent != "null" {
        Some((intent, status))
    } else {
        None
    }
}

fn bool_str(b: bool) -> &'static str {
    if b {
        "done"
    } else {
        "pending"
    }
}

fn memory_dir() -> PathBuf {
    let home = dirs::home_dir().unwrap_or_default();
    home.join(".local")
        .join("shared")
        .join("shared-ai")
        .join("memory")
}

fn stm_dir() -> PathBuf {
    let home = dirs::home_dir().unwrap_or_default();
    home.join(".local")
        .join("shared")
        .join("shared-ai")
        .join("stm")
}

fn count_toon_recursive(dir: &std::path::Path) -> usize {
    if !dir.exists() {
        return 0;
    }
    let mut count = 0;
    let Ok(entries) = std::fs::read_dir(dir) else {
        return 0;
    };
    for entry in entries.flatten() {
        let path = entry.path();
        if path.is_dir() {
            count += count_toon_recursive(&path);
        } else if path.extension().map(|e| e == "toon").unwrap_or(false) {
            count += 1;
        }
    }
    count
}
