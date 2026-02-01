pub mod agents;
pub mod gates;
pub mod lint;
pub mod memory;
pub mod orch;
pub mod quality;
pub mod session;
pub mod skills;
pub mod status;
pub mod telemetry;

use clap::{Parser, Subcommand};

#[derive(Parser)]
#[command(name = "kavach", about = "Brahmastra Stack - Universal AI CLI Enforcement")]
pub struct Cli {
    #[command(subcommand)]
    pub command: Option<Command>,
}

#[derive(Subcommand)]
pub enum Command {
    /// Print TOON status
    Status,
    /// Hook alias for gates intent
    #[command(name = "intent", hide = true)]
    Intent {
        #[arg(long)]
        hook: bool,
    },
    /// Gate subcommands
    Gates {
        #[command(subcommand)]
        command: gates::GatesCommand,
    },
    /// Memory subcommands
    Memory {
        #[command(subcommand)]
        command: memory::MemoryCommand,
    },
    /// Session subcommands
    Session {
        #[command(subcommand)]
        command: session::SessionCommand,
    },
    /// Orchestration subcommands
    Orch {
        #[command(subcommand)]
        command: orch::OrchCommand,
    },
    /// Agent management
    Agents(agents::AgentsArgs),
    /// Skill management
    Skills(skills::SkillsArgs),
    /// Lint files
    Lint(lint::LintArgs),
    /// Quality check
    Quality(quality::QualityArgs),
    /// Telemetry
    Telemetry {
        #[command(subcommand)]
        command: telemetry::TelemetryCommand,
    },
}

pub fn dispatch(cli: Cli) -> anyhow::Result<()> {
    match cli.command {
        None => {
            cli_print("kavach: Brahmastra Stack enforcement CLI. Use --help for commands.");
            Ok(())
        }
        Some(cmd) => match cmd {
            Command::Status => status::run(),
            Command::Intent { hook } => gates::intent::run(hook),
            Command::Gates { command } => gates::dispatch(command),
            Command::Memory { command } => memory::dispatch(command),
            Command::Session { command } => session::dispatch(command),
            Command::Orch { command } => orch::dispatch(command),
            Command::Agents(args) => agents::run(args),
            Command::Skills(args) => skills::run(args),
            Command::Lint(args) => lint::run(args),
            Command::Quality(args) => quality::run(args),
            Command::Telemetry { command } => telemetry::dispatch(command),
        },
    }
}

/// CLI stdout writer. This is a CLI binary — stdout is the correct output channel.
pub fn cli_print(msg: &str) {
    use std::io::Write;
    let stdout = std::io::stdout();
    let mut handle = stdout.lock();
    let _ = writeln!(handle, "{}", msg);
}

pub fn cli_print_fmt(msg: String) {
    cli_print(&msg);
}
