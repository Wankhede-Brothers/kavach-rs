use clap::Args;
use crate::commands::cli_print_fmt;

#[derive(Args)]
pub struct AegisArgs {
    #[arg(long)]
    pub hook: bool,
    #[arg(long)]
    pub task: bool,
}

pub fn run(args: AegisArgs) -> anyhow::Result<()> {
    cli_print_fmt(format!("[STUB] orch aegis: hook={} task={}", args.hook, args.task));
    Ok(())
}
