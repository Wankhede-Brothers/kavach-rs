use std::env;
use std::sync::OnceLock;

mod commands;
mod config;
mod hook;
mod patterns;
mod session;

use clap::Parser;
use commands::Cli;

/// Symlink dispatch table: binary name -> equivalent subcommand args
static SYMLINK_MAP: OnceLock<Vec<(&'static str, Vec<&'static str>)>> = OnceLock::new();

fn symlink_table() -> &'static Vec<(&'static str, Vec<&'static str>)> {
    SYMLINK_MAP.get_or_init(|| {
        vec![
            ("intent-gate", vec!["gates", "intent", "--hook"]),
            ("kavach-pre-write", vec!["gates", "pre-write", "--hook"]),
            ("kavach-post-write", vec!["gates", "post-write", "--hook"]),
            ("kavach-pre-tool", vec!["gates", "pre-tool", "--hook"]),
            ("kavach-post-tool", vec!["gates", "post-tool", "--hook"]),
            ("kavach-subagent", vec!["gates", "subagent", "--hook"]),
            ("ast-gate", vec!["gates", "ast", "--hook"]),
            ("bash-sanitizer", vec!["gates", "bash", "--hook"]),
            ("read-blocker", vec!["gates", "read", "--hook"]),
            ("skill-gate", vec!["gates", "skill", "--hook"]),
            ("lint-gate", vec!["gates", "lint", "--hook"]),
            ("research-gate", vec!["gates", "research", "--hook"]),
            ("content-gate", vec!["gates", "content", "--hook"]),
            ("quality-gate", vec!["gates", "quality", "--hook"]),
            ("enforcer", vec!["gates", "enforcer", "--hook"]),
            ("kavach-context", vec!["gates", "context", "--hook"]),
            ("kavach-dag", vec!["gates", "dag", "--hook"]),
            ("kavach-task", vec!["gates", "task", "--hook"]),
            ("kavach-code-guard", vec!["gates", "code-guard", "--hook"]),
            ("kavach-chain", vec!["gates", "chain", "--hook"]),
            ("ceo-gate", vec!["gates", "ceo", "--hook"]),
            ("kavach-failure", vec!["gates", "failure", "--hook"]),
            ("kavach-mockdata", vec!["gates", "mockdata", "--hook"]),
            ("memory-bank", vec!["memory", "bank"]),
            ("session-init", vec!["session", "init"]),
            ("memory-write", vec!["memory", "write"]),
            ("memory-rpc", vec!["memory", "rpc"]),
            ("stm-updater", vec!["memory", "stm"]),
            ("rpc-inject", vec!["memory", "inject"]),
            ("spec-inject", vec!["memory", "spec"]),
            ("session-validate", vec!["session", "validate"]),
            ("session-end", vec!["session", "end"]),
            ("session-resume", vec!["session", "resume"]),
            ("pre-compact", vec!["session", "compact"]),
            ("aegis-auto", vec!["orch", "aegis", "--hook"]),
            ("autonomous-orch", vec!["orch", "auto"]),
            ("post-verify", vec!["orch", "post"]),
        ]
    })
}

fn main() -> anyhow::Result<()> {
    let binary_name = env::args()
        .next()
        .and_then(|p| {
            std::path::Path::new(&p)
                .file_name()
                .map(|f| f.to_string_lossy().into_owned())
        })
        .unwrap_or_default();

    // Check symlink dispatch
    if let Some((_, ref mapped_args)) = symlink_table()
        .iter()
        .find(|(name, _)| *name == binary_name.as_str())
    {
        let mut full_args = vec!["kavach".to_string()];
        full_args.extend(mapped_args.iter().map(|s| s.to_string()));
        // Append any extra args from actual invocation
        full_args.extend(env::args().skip(1));
        let cli = Cli::parse_from(full_args);
        return commands::dispatch(cli);
    }

    let cli = Cli::parse();
    commands::dispatch(cli)
}
