use clap::Args;
use crate::commands::cli_print_fmt;

#[derive(Args)]
pub struct WriteArgs {
    #[arg(short = 'c', long)]
    pub category: Option<String>,
    #[arg(short = 'k', long)]
    pub key: Option<String>,
    #[arg(short = 'p', long)]
    pub project: Option<String>,
}

pub fn run(args: WriteArgs) -> anyhow::Result<()> {
    cli_print_fmt(format!("[STUB] memory write: category={:?} key={:?} project={:?}", args.category, args.key, args.project));
    Ok(())
}
