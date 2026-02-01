pub mod init;
pub mod end;
pub mod compact;
pub mod validate;
pub mod resume;

use clap::Subcommand;
use super::cli_print_fmt;

#[derive(Subcommand)]
pub enum SessionCommand {
    Init,
    Validate,
    End,
    Compact,
    Resume,
    Land,
    #[command(name = "end-hook")]
    EndHook,
}

pub fn dispatch(cmd: SessionCommand) -> anyhow::Result<()> {
    match cmd {
        SessionCommand::Init => init::run(),
        SessionCommand::Validate => validate::run(),
        SessionCommand::End => end::run(),
        SessionCommand::Compact => compact::run(),
        SessionCommand::Resume => resume::run(),
        SessionCommand::Land => { cli_print_fmt("[STUB] session land".into()); Ok(()) }
        SessionCommand::EndHook => { cli_print_fmt("[STUB] session end-hook".into()); Ok(()) }
    }
}
