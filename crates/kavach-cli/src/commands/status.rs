//! Status command: show system health in TOON format.

use std::io::Write;

use crate::session;

pub fn run() -> anyhow::Result<()> {
    let sess = session::get_or_create_session();
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[STATUS]")?;
    writeln!(w, "today: {}", sess.today)?;
    writeln!(w, "cutoff: {}", sess.training_cutoff)?;
    writeln!(w, "session: {}", sess.id)?;
    writeln!(w, "project: {}", sess.project)?;
    writeln!(w)?;

    writeln!(w, "[ENFORCE]")?;
    writeln!(w, "TABULA_RASA: active")?;
    writeln!(w, "DATE_INJECTION: {}", sess.today)?;
    writeln!(w, "NO_AMNESIA: ~/.local/shared/shared-ai/memory/")?;
    writeln!(w, "NO_ASSUMPTION: verify_before_act")?;
    writeln!(w, "DACE: lazy_load,skill_first")?;
    writeln!(w)?;

    writeln!(w, "[STATE]")?;
    writeln!(w, "research_done: {}", bool_state(sess.research_done))?;
    writeln!(w, "memory: {}", bool_state(sess.memory_queried))?;
    writeln!(w, "ceo: {}", bool_state(sess.ceo_invoked))?;
    writeln!(w, "aegis: {}", bool_state(sess.aegis_verified))?;
    writeln!(w, "turn_count: {}", sess.turn_count)?;
    if !sess.current_task.is_empty() {
        writeln!(w, "task: {} ({})", sess.current_task, sess.task_status)?;
    }

    Ok(())
}

fn bool_state(b: bool) -> &'static str {
    if b {
        "done"
    } else {
        "pending"
    }
}
