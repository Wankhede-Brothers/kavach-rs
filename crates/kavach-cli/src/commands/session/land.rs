//! Session land: finalize task, run Aegis verification, update memory.

use std::io::Write;

use crate::session;

pub fn run() -> anyhow::Result<()> {
    let mut sess = session::get_or_create_session();
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[SESSION:LAND]")?;
    writeln!(w, "session: {}", sess.id)?;
    writeln!(w, "project: {}", sess.project)?;
    writeln!(w, "date: {}", sess.today)?;
    writeln!(w)?;

    writeln!(w, "[TASK]")?;
    if sess.has_task() {
        writeln!(w, "task: {}", sess.current_task)?;
        writeln!(w, "status: {}", sess.task_status)?;
    } else {
        writeln!(w, "task: none")?;
    }
    writeln!(w)?;

    writeln!(w, "[METRICS]")?;
    writeln!(w, "turns: {}", sess.turn_count)?;
    writeln!(w, "tasks_created: {}", sess.tasks_created)?;
    writeln!(w, "tasks_completed: {}", sess.tasks_completed)?;
    writeln!(w, "files_modified: {}", sess.files_modified.len())?;
    writeln!(w, "research_done: {}", sess.research_done)?;
    writeln!(w, "aegis_verified: {}", sess.aegis_verified)?;
    writeln!(w)?;

    writeln!(w, "[CHECKLIST]")?;
    if !sess.aegis_verified {
        writeln!(w, "  - Run: kavach orch aegis")?;
    }
    writeln!(w, "  - Run: kavach memory sync")?;
    writeln!(w, "  - Run: kavach session end")?;

    sess.task_status = "landing".into();
    let _ = sess.save();

    Ok(())
}
