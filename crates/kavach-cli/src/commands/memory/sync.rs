use clap::Args;
use crate::commands::cli_print_fmt;

#[derive(Args)]
pub struct SyncArgs {
    #[arg(long)]
    pub hook: bool,
    #[arg(long)]
    pub task: bool,
    #[arg(long)]
    pub status: bool,
}

pub fn run(args: SyncArgs) -> anyhow::Result<()> {
    cli_print_fmt(format!("[STUB] memory sync: hook={} task={} status={}", args.hook, args.task, args.status));
    Ok(())
}
