//! Session init: initialize session with date injection.
//! Hook: SessionStart, UserPromptSubmit
//! DACE: ~100 tokens normal, ~50 tokens post-compact

use std::io::Write;

use crate::session;

pub fn run() -> anyhow::Result<()> {
    let mut sess = session::get_or_create_session();
    let out = std::io::stdout();
    let mut w = out.lock();

    if sess.is_post_compact() {
        run_post_compact(&mut w, &mut sess)?;
        return Ok(());
    }

    let session_type = get_session_type(&sess);

    writeln!(w, "[META]")?;
    writeln!(w, "protocol: SP/1.0")?;
    writeln!(w, "date: {}", sess.today)?;
    writeln!(w, "session: {}", sess.id)?;
    writeln!(w, "type: {session_type}")?;
    writeln!(w)?;

    writeln!(w, "[TABULA_RASA]")?;
    writeln!(w, "cutoff: {}", sess.training_cutoff)?;
    writeln!(w, "today: {}", sess.today)?;
    writeln!(w, "rule: WebSearch_BEFORE_code")?;
    writeln!(w, "blocked: I_think,I_believe,I_recall,Based_on_my_knowledge")?;
    writeln!(w)?;

    writeln!(w, "[NO_AMNESIA]")?;
    writeln!(w, "memory: ~/.local/shared/shared-ai/memory/")?;
    writeln!(w, "forbidden: I_have_no_memory,I_dont_have_access")?;
    writeln!(w)?;

    writeln!(w, "[SESSION]")?;
    writeln!(w, "id: {}", sess.id)?;
    writeln!(w, "project: {}", sess.project)?;
    writeln!(w, "research_mode: always")?;
    writeln!(w, "research_done: {}", sess.research_done)?;
    writeln!(w, "memory: {}", sess.memory_queried)?;
    writeln!(w)?;

    writeln!(w, "[MEMORY] total: 0 | query: kavach memory bank")?;
    writeln!(w)?;

    writeln!(w, "[DACE] mode: lazy_load,skill_first,on_demand")?;

    sess.mark_memory_queried();

    Ok(())
}

fn run_post_compact(w: &mut impl Write, sess: &mut session::SessionState) -> anyhow::Result<()> {
    sess.clear_post_compact();

    writeln!(w, "[META]")?;
    writeln!(w, "protocol: SP/1.0")?;
    writeln!(w, "date: {}", sess.today)?;
    writeln!(w, "session: {}", sess.id)?;
    writeln!(w, "type: post_compact_recovery")?;
    writeln!(w, "compact_count: {}", sess.compact_count)?;
    writeln!(w)?;

    writeln!(w, "[SESSION]")?;
    writeln!(w, "id: {}", sess.id)?;
    writeln!(w, "project: {}", sess.project)?;
    writeln!(w, "research_mode: always")?;
    writeln!(w, "research_done: {}", sess.research_done)?;
    writeln!(w, "memory: {}", sess.memory_queried)?;
    writeln!(w)?;

    if !sess.current_task.is_empty() {
        writeln!(w, "[TASK:RESTORED] {} | status: {}", sess.current_task, sess.task_status)?;
        writeln!(w)?;
    }

    writeln!(w, "[MEMORY] total: 0 | query: kavach memory bank")?;
    writeln!(w)?;

    writeln!(w, "[DACE] mode: lazy_load,skill_first,on_demand | CONTEXT_RESTORED")?;

    sess.mark_memory_queried();

    Ok(())
}

fn get_session_type(s: &session::SessionState) -> &'static str {
    if s.post_compact {
        "post_compact_recovery"
    } else if s.research_done || s.memory_queried {
        "resumed_session"
    } else {
        "fresh_session"
    }
}
