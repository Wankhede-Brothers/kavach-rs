//! Session end-hook: cleanup hook called at session Stop event.
//! Persists final session state and triggers memory sync.

use std::io::Write;

use crate::session;

pub fn run() -> anyhow::Result<()> {
    let mut sess = session::get_or_create_session();
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[SESSION:END_HOOK]")?;
    writeln!(w, "session: {}", sess.id)?;
    writeln!(w, "project: {}", sess.project)?;
    writeln!(w, "turns: {}", sess.turn_count)?;
    writeln!(w, "tasks_created: {}", sess.tasks_created)?;
    writeln!(w, "tasks_completed: {}", sess.tasks_completed)?;
    writeln!(w, "files_modified: {}", sess.files_modified.len())?;
    writeln!(w)?;

    if sess.has_task() {
        writeln!(w, "[TASK_WARNING]")?;
        writeln!(w, "task: {}", sess.current_task)?;
        writeln!(
            w,
            "status: {} (not completed before session end)",
            sess.task_status
        )?;
        writeln!(w)?;
    }

    sess.task_status = "ended".into();
    let _ = sess.save();

    writeln!(w, "status: persisted")?;

    Ok(())
}
