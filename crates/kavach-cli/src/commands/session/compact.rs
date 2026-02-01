//! Session compact: save state before context compaction.
//! Hook: PreCompact
//! DACE: ~50 tokens output

use std::io::Write;

use crate::session;

pub fn run() -> anyhow::Result<()> {
    let mut sess = session::get_or_create_session();
    let out = std::io::stdout();
    let mut w = out.lock();

    sess.mark_post_compact();

    writeln!(w, "[COMPACT:DACE]")?;
    writeln!(w, "date: {}", sess.today)?;
    writeln!(w, "session: {}", sess.id)?;
    writeln!(w, "project: {}", sess.project)?;
    writeln!(w, "compact_count: {}", sess.compact_count)?;

    if !sess.current_task.is_empty() {
        writeln!(w, "task_saved: {}", sess.current_task)?;
    }

    writeln!(w)?;
    writeln!(w, "[POST_COMPACT]")?;
    writeln!(w, "run: kavach session init")?;
    writeln!(w, "This auto-restores context from Memory Bank")?;

    Ok(())
}
