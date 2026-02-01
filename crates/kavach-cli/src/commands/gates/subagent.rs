//! Subagent gate: validates Task tool subagent_type and enforces delegation.
//! Routes to pre_tool_ceo logic for validation, adds orchestration context.

use std::collections::HashMap;

use crate::hook::{self, HookInput};
use crate::patterns;
use crate::session;

pub fn run(hook: bool) -> anyhow::Result<()> {
    if !hook {
        crate::commands::cli_print("gates subagent: use --hook flag");
        return Ok(());
    }

    let input = hook::must_read_hook_input();
    let result = dispatch(&input);

    match result {
        Ok(()) => Ok(()),
        Err(e) if hook::is_hook_exit(&e) => Ok(()),
        Err(e) => Err(e),
    }
}

fn dispatch(input: &HookInput) -> anyhow::Result<()> {
    let subagent_type = input.get_string("subagent_type");

    if subagent_type.is_empty() {
        hook::exit_block_toon("SUBAGENT", "Task_requires_subagent_type")?;
    }

    if !patterns::is_valid_agent(&subagent_type) {
        hook::exit_block_toon("SUBAGENT", &format!("unknown_agent:{subagent_type}"))?;
    }

    let prompt = input.get_string("prompt");
    if prompt.len() > 8000 {
        let mut kvs = HashMap::new();
        kvs.insert("warn".into(), format!("prompt_length:{}", prompt.len()));
        hook::exit_modify_toon("SUBAGENT", &mut kvs)?;
    }

    let mut sess = session::get_or_create_session();
    sess.mark_ceo_invoked();

    let engineers = [
        "backend-engineer",
        "frontend-engineer",
        "aegis-guardian",
    ];
    if engineers.contains(&subagent_type.as_str()) {
        let today = hook::today();
        let mut kvs = HashMap::new();
        kvs.insert("agent".into(), subagent_type);
        kvs.insert("date".into(), today);
        kvs.insert("cutoff".into(), "2025-01".into());
        kvs.insert("CEO_FLOW".into(), "DELEGATE->VERIFY->AEGIS".into());
        hook::exit_modify_toon("SUBAGENT_ORCHESTRATION", &mut kvs)?;
    }

    hook::exit_silent()
}
