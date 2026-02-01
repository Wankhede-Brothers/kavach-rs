//! Post-tool umbrella gate (PostToolUse for non-write tools).
//! Routes by tool name: memory | context | research | task

use std::collections::HashMap;

use crate::hook::{self, HookInput};
use crate::session;

pub fn run(hook_mode: bool) -> anyhow::Result<()> {
    if !hook_mode {
        crate::commands::cli_print("gates post-tool: use --hook flag");
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
    match input.get_tool_name() {
        "Bash" => post_tool_bash(input),
        "Read" | "Glob" | "Grep" => post_tool_read(input),
        "Task" => post_tool_task(input),
        "WebSearch" | "WebFetch" => post_tool_research(input),
        "TaskCreate" => post_tool_task_create(input),
        "TaskUpdate" => post_tool_task_update(input),
        "TaskOutput" => post_tool_task_output(input),
        _ => hook::exit_silent(),
    }
}

fn post_tool_bash(_input: &HookInput) -> anyhow::Result<()> {
    hook::exit_silent()
}

fn post_tool_read(_input: &HookInput) -> anyhow::Result<()> {
    hook::exit_silent()
}

fn post_tool_task(input: &HookInput) -> anyhow::Result<()> {
    let agent_type = input.get_string("subagent_type");
    if !agent_type.is_empty() {
        let mut kvs = HashMap::new();
        kvs.insert("agent".into(), agent_type);
        kvs.insert("status".into(), "completed".into());
        hook::exit_modify_toon("AGENT_COMPLETE", &mut kvs)?;
    }
    hook::exit_silent()
}

fn post_tool_research(input: &HookInput) -> anyhow::Result<()> {
    let mut session = session::get_or_create_session();

    let topic = if input.get_tool_name() == "WebSearch" {
        input.get_string("query")
    } else {
        input.get_string("url")
    };

    session.mark_research_done_with_topic(&topic);

    hook::exit_silent()
}

fn post_tool_task_create(input: &HookInput) -> anyhow::Result<()> {
    let mut session = session::get_or_create_session();
    let subject = input.get_string("subject");

    session.tasks_created += 1;
    session.set_current_task(&subject);

    hook::exit_silent()
}

fn post_tool_task_update(input: &HookInput) -> anyhow::Result<()> {
    let mut session = session::get_or_create_session();
    let status = input.get_string("status");

    if status == "completed" || status == "deleted" {
        session.tasks_completed += 1;
        session.clear_task();
    } else if status == "in_progress" {
        let subject = input.get_string("subject");
        if !subject.is_empty() {
            session.set_current_task(&subject);
        }
    }

    hook::exit_silent()
}

fn post_tool_task_output(_input: &HookInput) -> anyhow::Result<()> {
    hook::exit_silent()
}
