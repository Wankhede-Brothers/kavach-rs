//! Task health: monitors session task metrics and warns on anomalies.

use std::io::Write;

use crate::session;

pub fn run() -> anyhow::Result<()> {
    let sess = session::get_or_create_session();
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[TASK:HEALTH]")?;
    writeln!(w, "session: {}", sess.id)?;
    writeln!(w, "project: {}", sess.project)?;
    writeln!(w)?;

    writeln!(w, "[METRICS]")?;
    writeln!(w, "turns: {}", sess.turn_count)?;
    writeln!(w, "tasks_created: {}", sess.tasks_created)?;
    writeln!(w, "tasks_completed: {}", sess.tasks_completed)?;
    writeln!(w, "current_task: {}", if sess.current_task.is_empty() { "none" } else { &sess.current_task })?;
    writeln!(w, "task_status: {}", if sess.task_status.is_empty() { "idle" } else { &sess.task_status })?;
    writeln!(w)?;

    let mut warnings = Vec::new();

    // Warn if many turns without completing a task
    if sess.turn_count > 20 && sess.tasks_completed == 0 {
        warnings.push("no_tasks_completed_after_20_turns");
    }

    // Warn if task created but not started
    if sess.tasks_created > 0 && sess.task_status.is_empty() {
        warnings.push("tasks_created_but_none_active");
    }

    // Warn if many tasks created but few completed
    if sess.tasks_created > 5 && sess.tasks_completed < sess.tasks_created / 2 {
        warnings.push("low_completion_rate");
    }

    // Warn if no research done
    if !sess.research_done && sess.turn_count > 5 {
        warnings.push("no_research_done");
    }

    if warnings.is_empty() {
        writeln!(w, "[STATUS] HEALTHY")?;
    } else {
        writeln!(w, "[WARNINGS]")?;
        for warning in &warnings {
            writeln!(w, "  - {warning}")?;
        }
    }

    Ok(())
}
