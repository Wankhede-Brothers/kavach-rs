pub mod init;
pub mod end;
pub mod compact;
pub mod validate;
pub mod resume;
pub mod land;
pub mod end_hook;

use clap::Subcommand;

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
        SessionCommand::Land => land::run(),
        SessionCommand::EndHook => end_hook::run(),
    }
}
