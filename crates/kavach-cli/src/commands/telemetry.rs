use std::io::Write;

use clap::Subcommand;

use crate::session;

#[derive(Subcommand)]
pub enum TelemetryCommand {
    Report,
}

pub fn dispatch(cmd: TelemetryCommand) -> anyhow::Result<()> {
    match cmd {
        TelemetryCommand::Report => run_report(),
    }
}

fn run_report() -> anyhow::Result<()> {
    let sess = session::get_or_create_session();
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[TELEMETRY:REPORT]")?;
    writeln!(w, "session: {}", sess.id)?;
    writeln!(w, "project: {}", sess.project)?;
    writeln!(w, "date: {}", sess.today)?;
    writeln!(w)?;

    writeln!(w, "[SESSION_METRICS]")?;
    writeln!(w, "turns: {}", sess.turn_count)?;
    writeln!(w, "compact_count: {}", sess.compact_count)?;
    writeln!(w, "tasks_created: {}", sess.tasks_created)?;
    writeln!(w, "tasks_completed: {}", sess.tasks_completed)?;
    writeln!(w, "files_modified: {}", sess.files_modified.len())?;
    writeln!(w)?;

    writeln!(w, "[STATE_FLAGS]")?;
    writeln!(w, "research_done: {}", sess.research_done)?;
    writeln!(w, "memory_queried: {}", sess.memory_queried)?;
    writeln!(w, "ceo_invoked: {}", sess.ceo_invoked)?;
    writeln!(w, "nlu_parsed: {}", sess.nlu_parsed)?;
    writeln!(w, "aegis_verified: {}", sess.aegis_verified)?;
    writeln!(w)?;

    writeln!(w, "[INTENT]")?;
    writeln!(
        w,
        "type: {}",
        if sess.intent_type.is_empty() {
            "none"
        } else {
            &sess.intent_type
        }
    )?;
    writeln!(
        w,
        "domain: {}",
        if sess.intent_domain.is_empty() {
            "none"
        } else {
            &sess.intent_domain
        }
    )?;

    Ok(())
}
