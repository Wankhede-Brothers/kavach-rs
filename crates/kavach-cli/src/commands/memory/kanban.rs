use clap::Args;
use crate::commands::cli_print_fmt;

#[derive(Args)]
pub struct KanbanArgs {
    #[arg(long)]
    pub status: bool,
    #[arg(long)]
    pub visual: bool,
    #[arg(long)]
    pub sutra: bool,
    #[arg(short = 'p', long)]
    pub project: Option<String>,
}

pub fn run(args: KanbanArgs) -> anyhow::Result<()> {
    cli_print_fmt(format!("[STUB] memory kanban: status={} visual={} sutra={} project={:?}", args.status, args.visual, args.sutra, args.project));
    Ok(())
}
