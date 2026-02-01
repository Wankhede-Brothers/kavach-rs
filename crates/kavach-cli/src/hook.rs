//! Hook I/O infrastructure for Claude Code hooks.
//! Reads JSON from stdin, outputs JSON responses to stdout.
//! This is a CLI binary — stdout is the correct output channel for hook JSON.
//!
//! Design: exit_* functions return Err(HookExit) as a sentinel.
//! Callers propagate with `?`. The top-level `run()` catches HookExit and returns Ok(()).

use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::collections::HashMap;
use std::io::{BufReader, Write};

/// Hook input from Claude Code (read from stdin).
#[derive(Debug, Deserialize, Default)]
#[serde(default)]
pub struct HookInput {
    #[serde(alias = "sessionId")]
    pub session_id: String,
    #[serde(alias = "toolName")]
    pub tool_name: String,
    #[serde(alias = "toolInput")]
    pub tool_input: Option<Value>,
    pub prompt: String,
    #[serde(alias = "hookEventName")]
    pub hook_event_name: String,
}

impl HookInput {
    pub fn get_tool_name(&self) -> &str {
        &self.tool_name
    }

    pub fn get_string(&self, key: &str) -> String {
        self.tool_input
            .as_ref()
            .and_then(|v| v.get(key))
            .and_then(|v| v.as_str())
            .unwrap_or("")
            .to_string()
    }

    pub fn get_prompt(&self) -> String {
        if !self.prompt.is_empty() {
            return self.prompt.clone();
        }
        self.get_string("prompt")
    }
}

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct HookSpecificOutput {
    pub hook_event_name: String,
    pub permission_decision: String,
    pub permission_decision_reason: String,
    #[serde(skip_serializing_if = "String::is_empty")]
    pub additional_context: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub updated_input: Option<Value>,
}

#[derive(Debug, Serialize)]
pub struct HookResponse {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub decision: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub reason: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    #[serde(rename = "additionalContext")]
    pub additional_context: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    #[serde(rename = "hookSpecificOutput")]
    pub hook_specific_output: Option<HookSpecificOutput>,
}

#[derive(Debug, Serialize)]
#[serde(rename_all = "camelCase")]
pub struct UserPromptSubmitResponse {
    pub hook_event_name: String,
    pub additional_context: String,
}

/// Sentinel error: hook already wrote JSON to stdout, caller should return Ok(()).
#[derive(Debug)]
pub struct HookExit;
impl std::fmt::Display for HookExit {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.write_str("hook exit")
    }
}
impl std::error::Error for HookExit {}

/// Check if an anyhow::Error is a HookExit sentinel.
pub fn is_hook_exit(err: &anyhow::Error) -> bool {
    err.downcast_ref::<HookExit>().is_some()
}

// --- Internal: write JSON to stdout (CLI binary output channel) ---
fn cli_write_json<T: Serialize>(v: &T) {
    if let Ok(data) = serde_json::to_string(v) {
        let stdout = std::io::stdout();
        let mut handle = stdout.lock();
        let _ = writeln!(handle, "{}", data);
    }
}

// --- Read ---

pub fn read_hook_input() -> anyhow::Result<HookInput> {
    let reader = BufReader::new(std::io::stdin().lock());
    let input: HookInput = serde_json::from_reader(reader)?;
    Ok(input)
}

pub fn must_read_hook_input() -> HookInput {
    read_hook_input().unwrap_or_default()
}

// --- Helpers ---

pub fn today() -> String {
    chrono::Local::now().format("%Y-%m-%d").to_string()
}

pub fn toon_block(name: &str, kvs: &HashMap<String, String>) -> String {
    let mut result = format!("[{name}]\n");
    for (k, v) in kvs {
        result += &format!("{k}: {v}\n");
    }
    result
}

// --- Exit functions (write JSON + return Err(HookExit)) ---

pub fn exit_silent() -> anyhow::Result<()> {
    cli_write_json(&HookResponse {
        decision: Some("approve".into()),
        reason: Some("ok".into()),
        additional_context: None,
        hook_specific_output: None,
    });
    Err(HookExit.into())
}

pub fn exit_block_toon(gate: &str, reason: &str) -> anyhow::Result<()> {
    let mut kvs = HashMap::new();
    kvs.insert("gate".into(), gate.into());
    kvs.insert("reason".into(), reason.into());
    kvs.insert("date".into(), today());
    let ctx = toon_block("BLOCK", &kvs);
    cli_write_json(&HookResponse {
        decision: None,
        reason: None,
        additional_context: None,
        hook_specific_output: Some(HookSpecificOutput {
            hook_event_name: "PreToolUse".into(),
            permission_decision: "deny".into(),
            permission_decision_reason: reason.into(),
            additional_context: ctx,
            updated_input: None,
        }),
    });
    Err(HookExit.into())
}

pub fn exit_modify_toon(gate: &str, kvs: &mut HashMap<String, String>) -> anyhow::Result<()> {
    kvs.insert("date".into(), today());
    let ctx = toon_block(gate, kvs);
    cli_write_json(&HookResponse {
        decision: Some("approve".into()),
        reason: Some(gate.into()),
        additional_context: Some(ctx),
        hook_specific_output: None,
    });
    Err(HookExit.into())
}

pub fn exit_modify_toon_with_module(
    gate: &str,
    kvs: &mut HashMap<String, String>,
    module_content: &str,
) -> anyhow::Result<()> {
    kvs.insert("date".into(), today());
    let mut ctx = toon_block(gate, kvs);
    if !module_content.is_empty() {
        ctx += &format!("\n[MODULE:LAZY_LOADED]\n{module_content}");
    }
    cli_write_json(&HookResponse {
        decision: Some("approve".into()),
        reason: Some(gate.into()),
        additional_context: Some(ctx),
        hook_specific_output: None,
    });
    Err(HookExit.into())
}

pub fn exit_user_prompt_submit_silent() -> anyhow::Result<()> {
    cli_write_json(&UserPromptSubmitResponse {
        hook_event_name: "UserPromptSubmit".into(),
        additional_context: String::new(),
    });
    Err(HookExit.into())
}

pub fn exit_user_prompt_submit_with_context(context: &str) -> anyhow::Result<()> {
    cli_write_json(&UserPromptSubmitResponse {
        hook_event_name: "UserPromptSubmit".into(),
        additional_context: context.into(),
    });
    Err(HookExit.into())
}
