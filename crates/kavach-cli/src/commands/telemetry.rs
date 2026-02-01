use clap::Subcommand;
use super::cli_print_fmt;

#[derive(Subcommand)]
pub enum TelemetryCommand {
    Report,
}

pub fn dispatch(cmd: TelemetryCommand) -> anyhow::Result<()> {
    match cmd {
        TelemetryCommand::Report => {
            cli_print_fmt("[STUB] telemetry report".into());
            Ok(())
        }
    }
}
