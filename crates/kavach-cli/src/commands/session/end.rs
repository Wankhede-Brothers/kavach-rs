//! Session end: persist state before termination.
//! Hook: Stop

use std::io::Write;

use crate::session;

pub fn run() -> anyhow::Result<()> {
    let sess = session::get_or_create_session();
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[END]")?;
    writeln!(w, "date: {}", sess.today)?;
    writeln!(w, "session: {}", sess.id)?;
    writeln!(w, "project: {}", sess.project)?;
    writeln!(w)?;

    writeln!(w, "[STATE]")?;
    writeln!(w, "research_done: {}", sess.research_done)?;
    writeln!(w, "memory: {}", sess.memory_queried)?;
    writeln!(w, "ceo: {}", sess.ceo_invoked)?;
    writeln!(w, "aegis: {}", sess.aegis_verified)?;

    let _ = sess.save();

    Ok(())
}
