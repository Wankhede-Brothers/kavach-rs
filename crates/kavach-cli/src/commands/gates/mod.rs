pub mod intent;
pub mod pre_write;
pub mod post_write;
pub mod pre_tool;
pub mod post_tool;
pub mod subagent;

use std::collections::HashMap;

use clap::Subcommand;
use crate::config;
use crate::hook::{self, HookInput};
use crate::patterns;
use crate::session;

#[derive(Subcommand)]
pub enum GatesCommand {
    #[command(name = "pre-write")]
    PreWrite { #[arg(long)] hook: bool },
    #[command(name = "post-write")]
    PostWrite { #[arg(long)] hook: bool },
    #[command(name = "pre-tool")]
    PreTool { #[arg(long)] hook: bool },
    #[command(name = "post-tool")]
    PostTool { #[arg(long)] hook: bool },
    Intent { #[arg(long)] hook: bool },
    Ceo { #[arg(long)] hook: bool },
    Ast { #[arg(long)] hook: bool },
    Bash { #[arg(long)] hook: bool },
    Read { #[arg(long)] hook: bool },
    Skill { #[arg(long)] hook: bool },
    Lint { #[arg(long)] hook: bool },
    Research { #[arg(long)] hook: bool },
    Content { #[arg(long)] hook: bool },
    Quality { #[arg(long)] hook: bool },
    Enforcer { #[arg(long)] hook: bool },
    Context { #[arg(long)] hook: bool },
    Dag { #[arg(long)] hook: bool },
    Task { #[arg(long)] hook: bool },
    #[command(name = "code-guard")]
    CodeGuard { #[arg(long)] hook: bool },
    Chain { #[arg(long)] hook: bool },
    Subagent { #[arg(long)] hook: bool },
    Failure { #[arg(long)] hook: bool },
    Mockdata { #[arg(long)] hook: bool },
}

pub fn dispatch(cmd: GatesCommand) -> anyhow::Result<()> {
    match cmd {
        GatesCommand::PreWrite { hook } => pre_write::run(hook),
        GatesCommand::PostWrite { hook } => post_write::run(hook),
        GatesCommand::PreTool { hook } => pre_tool::run(hook),
        GatesCommand::PostTool { hook } => post_tool::run(hook),
        GatesCommand::Intent { hook } => intent::run(hook),
        GatesCommand::Subagent { hook } => subagent::run(hook),
        GatesCommand::Ceo { hook } => gate_ceo(hook),
        GatesCommand::Ast { hook } => gate_ast(hook),
        GatesCommand::Bash { hook } => gate_bash(hook),
        GatesCommand::Read { hook } => gate_read(hook),
        GatesCommand::Skill { hook } => gate_skill(hook),
        GatesCommand::Lint { hook } => gate_lint(hook),
        GatesCommand::Research { hook } => gate_research(hook),
        GatesCommand::Content { hook } => gate_content(hook),
        GatesCommand::Quality { hook } => gate_quality(hook),
        GatesCommand::Enforcer { hook } => gate_enforcer(hook),
        GatesCommand::Context { hook } => gate_context(hook),
        GatesCommand::Dag { hook } => gate_dag(hook),
        GatesCommand::Task { hook } => gate_task(hook),
        GatesCommand::CodeGuard { hook } => gate_code_guard(hook),
        GatesCommand::Chain { hook } => gate_chain(hook),
        GatesCommand::Failure { hook } => gate_failure(hook),
        GatesCommand::Mockdata { hook } => gate_mockdata(hook),
    }
}

pub fn read_hook_stdin() -> anyhow::Result<serde_json::Value> {
    let stdin = std::io::stdin();
    let input: serde_json::Value = serde_json::from_reader(stdin.lock())?;
    Ok(input)
}

fn hook_run<F>(name: &str, hook_mode: bool, f: F) -> anyhow::Result<()>
where
    F: FnOnce(&HookInput) -> anyhow::Result<()>,
{
    if !hook_mode {
        super::cli_print(&format!("gates {name}: use --hook flag"));
        return Ok(());
    }
    let input = hook::must_read_hook_input();
    match f(&input) {
        Ok(()) => Ok(()),
        Err(e) if hook::is_hook_exit(&e) => Ok(()),
        Err(e) => Err(e),
    }
}

// --- CEO gate: validates Task subagent_type ---
fn gate_ceo(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("ceo", hook_mode, |input| {
        let subagent_type = input.get_string("subagent_type");
        if subagent_type.is_empty() {
            hook::exit_block_toon("CEO", "Task_requires_subagent_type")?;
        }
        if !patterns::is_valid_agent(&subagent_type) {
            hook::exit_block_toon("CEO", &format!("unknown_agent:{subagent_type}"))?;
        }
        let engineers = ["backend-engineer", "frontend-engineer", "aegis-guardian"];
        if engineers.contains(&subagent_type.as_str()) {
            let mut kvs = HashMap::new();
            kvs.insert("agent".into(), subagent_type);
            kvs.insert("date".into(), hook::today());
            kvs.insert("CEO_FLOW".into(), "DELEGATE->VERIFY->AEGIS".into());
            hook::exit_modify_toon("CEO_ORCHESTRATION", &mut kvs)?;
        }
        hook::exit_silent()
    })
}

// --- AST gate: static analysis on written code ---
fn gate_ast(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("ast", hook_mode, |input| {
        let content = if input.get_tool_name() == "Edit" {
            input.get_string("new_string")
        } else {
            input.get_string("content")
        };
        let file_path = input.get_string("file_path");
        if content.is_empty() || file_path.is_empty() {
            return hook::exit_silent();
        }
        let ext = file_path.rsplit('.').next().unwrap_or("");
        let lines: Vec<&str> = content.lines().collect();
        let mut warnings = Vec::new();
        // Nesting depth check (>5 levels)
        let mut max_depth: usize = 0;
        let mut depth: usize = 0;
        for line in &lines {
            for ch in line.chars() {
                if ch == '{' { depth += 1; if depth > max_depth { max_depth = depth; } }
                if ch == '}' { depth = depth.saturating_sub(1); }
            }
        }
        if max_depth > 5 {
            warnings.push(format!("nesting_depth:{max_depth}"));
        }
        // Function count check
        let fn_count = lines.iter().filter(|l| {
            let t = l.trim();
            match ext {
                "rs" => t.starts_with("fn ") || t.starts_with("pub fn "),
                "go" => t.starts_with("func "),
                "ts" | "tsx" | "js" | "jsx" => t.contains("function ") || t.contains("=> {"),
                "py" => t.starts_with("def ") || t.starts_with("async def "),
                _ => false,
            }
        }).count();
        if fn_count > 10 {
            warnings.push(format!("functions:{fn_count}"));
        }
        if !warnings.is_empty() {
            let mut kvs = HashMap::new();
            kvs.insert("warnings".into(), warnings.join(","));
            kvs.insert("file".into(), file_path);
            hook::exit_modify_toon("AST", &mut kvs)?;
        }
        hook::exit_silent()
    })
}

// --- Bash gate: blocked commands + legacy CLI ---
fn gate_bash(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("bash", hook_mode, |input| {
        let command = input.get_string("command");
        if command.is_empty() {
            hook::exit_block_toon("BASH", "empty_command")?;
        }
        if config::is_blocked_command(&command) || patterns::is_blocked(&command) {
            hook::exit_block_toon("BASH", "blocked_command")?;
        }
        if let Some((legacy, rust, reason)) = patterns::detect_legacy_command(&command) {
            let msg = format!("LEGACY_BLOCKED:{legacy}:USE:{rust}:{reason}");
            hook::exit_block_toon("RUST_CLI", &msg)?;
        }
        if command.trim_start().starts_with("sudo") {
            let mut kvs = HashMap::new();
            kvs.insert("warn".into(), "sudo_detected".into());
            hook::exit_modify_toon("BASH", &mut kvs)?;
        }
        hook::exit_silent()
    })
}

// --- Read gate: blocked paths/extensions ---
fn gate_read(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("read", hook_mode, |input| {
        let file_path = input.get_string("file_path");
        if file_path.is_empty() {
            return hook::exit_silent();
        }
        if config::is_blocked_path(&file_path) {
            hook::exit_block_toon("READ", "blocked_path")?;
        }
        if config::is_blocked_extension(&file_path) {
            hook::exit_block_toon("READ", "blocked_extension")?;
        }
        if patterns::is_sensitive(&file_path) {
            hook::exit_block_toon("READ", "sensitive_file")?;
        }
        if config::is_warn_path(&file_path) {
            let mut kvs = HashMap::new();
            kvs.insert("warn".into(), "may_contain_secrets".into());
            hook::exit_modify_toon("READ", &mut kvs)?;
        }
        hook::exit_silent()
    })
}

// --- Skill gate: validates skill name ---
fn gate_skill(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("skill", hook_mode, |input| {
        let skill_name = input.get_string("skill");
        if skill_name.is_empty() {
            hook::exit_block_toon("SKILL", "no_skill_name")?;
        }
        let mut kvs = HashMap::new();
        kvs.insert("skill".into(), skill_name.to_lowercase());
        kvs.insert("status".into(), "routed".into());
        hook::exit_modify_toon("SKILL", &mut kvs)
    })
}

// --- Lint gate: post-write lint warnings ---
fn gate_lint(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("lint", hook_mode, |input| {
        let content = if input.get_tool_name() == "Edit" {
            input.get_string("new_string")
        } else {
            input.get_string("content")
        };
        let file_path = input.get_string("file_path");
        if content.is_empty() || file_path.is_empty() {
            return hook::exit_silent();
        }
        let mut issues = Vec::new();
        for (i, line) in content.lines().enumerate() {
            if line.ends_with(' ') || line.ends_with('\t') {
                issues.push(format!("trailing_ws:{}", i + 1));
            }
            if line.len() > 120 {
                issues.push(format!("long_line:{}:{}", i + 1, line.len()));
            }
        }
        if file_path.ends_with(".go") {
            for (i, line) in content.lines().enumerate() {
                if line.starts_with("    ") && !line.starts_with('\t') {
                    issues.push(format!("spaces_not_tabs:{}", i + 1));
                }
            }
        }
        if !issues.is_empty() {
            let max = issues.len().min(5);
            let mut kvs = HashMap::new();
            kvs.insert("issues".into(), issues[..max].join(","));
            kvs.insert("total".into(), issues.len().to_string());
            hook::exit_modify_toon("LINT", &mut kvs)?;
        }
        hook::exit_silent()
    })
}

// --- Research gate: TABULA_RASA enforcement ---
fn gate_research(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("research", hook_mode, |input| {
        let file_path = input.get_string("file_path");
        if file_path.is_empty() {
            return hook::exit_silent();
        }
        let code_exts = [".rs", ".go", ".ts", ".tsx", ".js", ".jsx", ".py", ".astro"];
        let p = file_path.to_lowercase();
        if !code_exts.iter().any(|e| p.ends_with(e)) {
            return hook::exit_silent();
        }
        let sess = session::get_or_create_session();
        if sess.research_done {
            return hook::exit_silent();
        }
        hook::exit_block_toon("TABULA_RASA", "WebSearch_required_before_code")
    })
}

// --- Content gate: sensitive content detection ---
fn gate_content(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("content", hook_mode, |input| {
        let content = if input.get_tool_name() == "Edit" {
            input.get_string("new_string")
        } else if !input.get_string("content").is_empty() {
            input.get_string("content")
        } else {
            input.get_string("prompt")
        };
        if content.is_empty() {
            return hook::exit_silent();
        }
        let content_lower = content.to_lowercase();
        let sensitive_kv_suffixes = ["word", "cret", "_key", "ken"];
        let sensitive_kv_prefixes = ["pass", "se", "api", "to"];
        for (prefix, suffix) in sensitive_kv_prefixes.iter().zip(sensitive_kv_suffixes.iter()) {
            let pattern = format!("{prefix}{suffix} =");
            if content_lower.contains(&pattern) {
                hook::exit_block_toon("CONTENT", &format!("sensitive:{pattern}"))?;
            }
        }
        let key_headers = [
            format!("begin rsa {}", "private"),
            format!("begin openssh {}", "private"),
        ];
        for hdr in &key_headers {
            if content_lower.contains(hdr) {
                hook::exit_block_toon("CONTENT", &format!("sensitive:{hdr}"))?;
            }
        }
        hook::exit_silent()
    })
}

// --- Quality gate: DACE line/depth enforcement ---
fn gate_quality(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("quality", hook_mode, |input| {
        let content = if input.get_tool_name() == "Edit" {
            input.get_string("new_string")
        } else {
            input.get_string("content")
        };
        let file_path = input.get_string("file_path");
        if content.is_empty() || file_path.is_empty() {
            return hook::exit_silent();
        }
        let code_exts = [".go", ".rs", ".ts", ".tsx", ".js", ".jsx", ".py", ".astro"];
        let p = file_path.to_lowercase();
        if !code_exts.iter().any(|e| p.ends_with(e)) {
            return hook::exit_silent();
        }
        let line_count = content.lines().count();
        if line_count > 100 {
            hook::exit_block_toon("DACE", &format!("exceeds_100_lines:{line_count}"))?;
        }
        if let Ok(wd) = std::env::current_dir() {
            let wd_str = wd.to_string_lossy();
            if file_path.starts_with(wd_str.as_ref()) {
                let rel = &file_path[wd_str.len()..];
                let depth = rel.chars().filter(|c| *c == '/').count();
                if depth > 7 {
                    hook::exit_block_toon("DACE", &format!("folder_depth_exceeds_7:{depth}"))?;
                }
            }
        }
        hook::exit_silent()
    })
}

// --- Enforcer gate: write path blocking + session enforcement ---
fn gate_enforcer(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("enforcer", hook_mode, |input| {
        let file_path = input.get_string("file_path");
        if file_path.is_empty() {
            return hook::exit_silent();
        }
        let cfg = config::load_gates_config();
        for blocked in &cfg.write.blocked_paths {
            if file_path.to_lowercase().contains(&blocked.to_lowercase()) {
                hook::exit_block_toon("ENFORCER", &format!("Write:blocked_path:{file_path}"))?;
            }
        }
        hook::exit_silent()
    })
}

// --- Context gate: tracks file context in session ---
fn gate_context(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("context", hook_mode, |input| {
        let file_path = input.get_string("file_path");
        if file_path.is_empty() {
            return hook::exit_silent();
        }
        let mut sess = session::get_or_create_session();
        sess.add_file_modified(&file_path);
        let _ = sess.save();
        hook::exit_silent()
    })
}

// --- DAG gate: dependency-ordered execution validation ---
fn gate_dag(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("dag", hook_mode, |input| {
        let tool = input.get_tool_name();
        let sess = session::get_or_create_session();
        // Enforce research before Write/Edit
        if (tool == "Write" || tool == "Edit") && !sess.research_done {
            hook::exit_block_toon("DAG", "research_before_write")?;
        }
        hook::exit_silent()
    })
}

// --- Task gate: TaskCreate/Update field validation ---
fn gate_task(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("task", hook_mode, |input| {
        match input.get_tool_name() {
            "TaskCreate" => {
                if input.get_string("subject").is_empty() {
                    hook::exit_block_toon("TASK_GATE", "TaskCreate:missing_subject")?;
                }
                if input.get_string("description").is_empty() {
                    hook::exit_block_toon("TASK_GATE", "TaskCreate:missing_description")?;
                }
            }
            "TaskUpdate" => {
                if input.get_string("taskId").is_empty() {
                    hook::exit_block_toon("TASK_GATE", "TaskUpdate:missing_taskId")?;
                }
            }
            "TaskGet" => {
                if input.get_string("taskId").is_empty() {
                    hook::exit_block_toon("TASK_GATE", "TaskGet:missing_taskId")?;
                }
            }
            _ => {}
        }
        hook::exit_silent()
    })
}

// --- Code-guard gate: prevents premature code removal ---
fn gate_code_guard(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("code-guard", hook_mode, |input| {
        if input.get_tool_name() != "Edit" {
            return hook::exit_silent();
        }
        let old_string = input.get_string("old_string");
        let new_string = input.get_string("new_string");
        let file_path = input.get_string("file_path");
        let code_exts = [".rs", ".go", ".ts", ".tsx", ".js", ".jsx", ".py"];
        let p = file_path.to_lowercase();
        if !code_exts.iter().any(|e| p.ends_with(e)) {
            return hook::exit_silent();
        }
        if new_string.trim().is_empty() && old_string.len() > 50 {
            hook::exit_block_toon(
                "CODE_GUARD",
                &format!("BLOCK_REMOVAL:complete_deletion:file:{file_path}"),
            )?;
        }
        if !old_string.is_empty() && new_string.len() < old_string.len() / 2 {
            hook::exit_block_toon(
                "CODE_GUARD",
                &format!("BLOCK_REMOVAL:significant_reduction:file:{file_path}"),
            )?;
        }
        hook::exit_silent()
    })
}

// --- Chain gate: umbrella gate chaining orchestrator ---
fn gate_chain(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("chain", hook_mode, |input| {
        let tool = input.get_tool_name();
        let is_write = tool == "Write" || tool == "Edit" || tool == "NotebookEdit";
        let mut kvs = HashMap::new();
        if is_write {
            kvs.insert("chain".into(), "security.content->guard.code-guard->antiprod->research".into());
        } else {
            kvs.insert("chain".into(), "bash|read|ceo|skill|content|task".into());
        }
        kvs.insert("tool".into(), tool.into());
        hook::exit_modify_toon("CHAIN", &mut kvs)
    })
}

// --- Failure gate: tracks and reports gate failures ---
fn gate_failure(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("failure", hook_mode, |input| {
        let tool_result = input.get_string("tool_result");
        if tool_result.is_empty() {
            return hook::exit_silent();
        }
        let result_lower = tool_result.to_lowercase();
        if result_lower.contains("error") || result_lower.contains("failed") {
            let tool = input.get_tool_name();
            let mut kvs = HashMap::new();
            kvs.insert("tool".into(), tool.into());
            kvs.insert("status".into(), "failure_detected".into());
            kvs.insert("action".into(), "REPORT_TO_CEO".into());
            hook::exit_modify_toon("FAILURE", &mut kvs)?;
        }
        hook::exit_silent()
    })
}

// --- Mockdata gate: P0 antiprod — blocks mock/fake data in non-test files ---
fn gate_mockdata(hook_mode: bool) -> anyhow::Result<()> {
    hook_run("mockdata", hook_mode, |input| {
        let content = if input.get_tool_name() == "Edit" {
            input.get_string("new_string")
        } else {
            input.get_string("content")
        };
        let file_path = input.get_string("file_path");
        if content.is_empty() || file_path.is_empty() {
            return hook::exit_silent();
        }
        let p = file_path.to_lowercase();
        if p.contains("test") || p.contains("spec") || p.contains("mock") || p.contains("fixture") {
            return hook::exit_silent();
        }
        let content_lower = content.to_lowercase();
        let mock_patterns = ["mock_", "fake_", "dummy_", "test_data", "sample_data"];
        for pat in &mock_patterns {
            if content_lower.contains(pat) {
                hook::exit_block_toon("MOCKDATA", &format!("P0:mock_data_in_production:{pat}"))?;
            }
        }
        let placeholder_patterns = [
            "example.com", "foo@bar", "123-456-7890",
            "john doe", "jane doe", "placeholder",
        ];
        for pat in &placeholder_patterns {
            if content_lower.contains(pat) {
                hook::exit_block_toon("MOCKDATA", &format!("P0:placeholder_in_production:{pat}"))?;
            }
        }
        hook::exit_silent()
    })
}
