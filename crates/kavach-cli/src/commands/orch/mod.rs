pub mod aegis;
pub mod dag;

use clap::Subcommand;
use super::cli_print_fmt;

#[derive(Subcommand)]
pub enum OrchCommand {
    Aegis(aegis::AegisArgs),
    Verify,
    #[command(name = "task-health")]
    TaskHealth,
    Dag,
}

pub fn dispatch(cmd: OrchCommand) -> anyhow::Result<()> {
    match cmd {
        OrchCommand::Aegis(args) => aegis::run(args),
        OrchCommand::Verify => { cli_print_fmt("[STUB] orch verify".into()); Ok(()) }
        OrchCommand::TaskHealth => { cli_print_fmt("[STUB] orch task-health".into()); Ok(()) }
        OrchCommand::Dag => dag::run(),
    }
}
