//! Session validate: check session health + memory bank accessibility.
//! Returns PASS/FAIL with exit code for CI/hook use.

use std::io::Write;
use std::path::Path;

use crate::session;

pub fn run() -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();
    let today = chrono::Local::now().format("%Y-%m-%d").to_string();
    let mut pass = true;

    writeln!(w, "[VALIDATE]")?;
    writeln!(w, "date: {today}")?;
    writeln!(w)?;

    // 1. Session state exists and date matches
    let session = session::load_session_state();
    match &session {
        Some(s) if s.today == today => {
            writeln!(w, "session: PASS (id: {})", s.session_id)?;
        }
        Some(s) => {
            writeln!(w, "session: FAIL (stale date: {} != {today})", s.today)?;
            pass = false;
        }
        None => {
            writeln!(w, "session: FAIL (no session state)")?;
            pass = false;
        }
    }

    // 2. Memory bank directory accessible
    let mem_dir = memory_dir();
    if mem_dir.exists() && mem_dir.is_dir() {
        let count = count_toon_recursive(&mem_dir);
        writeln!(w, "memory_bank: PASS ({count} files)")?;
    } else {
        writeln!(w, "memory_bank: FAIL (directory missing)")?;
        pass = false;
    }

    // 3. Governance file
    let gov = mem_dir.join("GOVERNANCE.toon");
    if gov.exists() {
        writeln!(w, "governance: PASS")?;
    } else {
        writeln!(w, "governance: WARN (not found)")?;
    }

    // 4. STM directory
    let stm = stm_dir();
    if stm.exists() {
        writeln!(w, "stm: PASS")?;
    } else {
        writeln!(w, "stm: WARN (not found)")?;
    }

    // 5. Research state
    if let Some(s) = &session {
        let research_status = if s.research_done { "done" } else { "pending" };
        writeln!(w, "research: {research_status}")?;
    }

    writeln!(w)?;
    if pass {
        writeln!(w, "[RESULT] PASS")?;
    } else {
        writeln!(w, "[RESULT] FAIL")?;
        std::process::exit(1);
    }

    Ok(())
}

fn memory_dir() -> std::path::PathBuf {
    let home = dirs::home_dir().unwrap_or_default();
    home.join(".local")
        .join("shared")
        .join("shared-ai")
        .join("memory")
}

fn stm_dir() -> std::path::PathBuf {
    let home = dirs::home_dir().unwrap_or_default();
    home.join(".local")
        .join("shared")
        .join("shared-ai")
        .join("stm")
}

fn count_toon_recursive(dir: &Path) -> usize {
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
