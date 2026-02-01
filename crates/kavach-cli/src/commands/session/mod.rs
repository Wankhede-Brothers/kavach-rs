pub mod compact;
pub mod end;
pub mod end_hook;
pub mod init;
pub mod land;
pub mod resume;
pub mod validate;

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
