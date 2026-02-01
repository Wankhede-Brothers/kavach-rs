pub mod intent;
pub mod pre_write;
pub mod post_write;
pub mod pre_tool;
pub mod post_tool;
pub mod subagent;

use clap::Subcommand;
use super::cli_print_fmt;

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
        // All other gates use the generic stub
        GatesCommand::Ceo { hook } => gate_stub("ceo", hook),
        GatesCommand::Ast { hook } => gate_stub("ast", hook),
        GatesCommand::Bash { hook } => gate_stub("bash", hook),
        GatesCommand::Read { hook } => gate_stub("read", hook),
        GatesCommand::Skill { hook } => gate_stub("skill", hook),
        GatesCommand::Lint { hook } => gate_stub("lint", hook),
        GatesCommand::Research { hook } => gate_stub("research", hook),
        GatesCommand::Content { hook } => gate_stub("content", hook),
        GatesCommand::Quality { hook } => gate_stub("quality", hook),
        GatesCommand::Enforcer { hook } => gate_stub("enforcer", hook),
        GatesCommand::Context { hook } => gate_stub("context", hook),
        GatesCommand::Dag { hook } => gate_stub("dag", hook),
        GatesCommand::Task { hook } => gate_stub("task", hook),
        GatesCommand::CodeGuard { hook } => gate_stub("code-guard", hook),
        GatesCommand::Chain { hook } => gate_stub("chain", hook),
        GatesCommand::Failure { hook } => gate_stub("failure", hook),
        GatesCommand::Mockdata { hook } => gate_stub("mockdata", hook),
    }
}

fn gate_stub(name: &str, hook: bool) -> anyhow::Result<()> {
    if hook {
        // Read stdin, approve
        let _input = read_hook_stdin().ok();
    }
    cli_print_fmt(format!("[STUB] gates {name}: hook={hook}"));
    Ok(())
}

pub fn read_hook_stdin() -> anyhow::Result<serde_json::Value> {
    let stdin = std::io::stdin();
    let input: serde_json::Value = serde_json::from_reader(stdin.lock())?;
    Ok(input)
}
