pub mod aegis;
pub mod dag;
pub mod verify;
pub mod task_health;

use clap::Subcommand;

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
        OrchCommand::Verify => verify::run(),
        OrchCommand::TaskHealth => task_health::run(),
        OrchCommand::Dag => dag::run(),
    }
}
