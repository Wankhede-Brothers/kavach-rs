use clap::Args;
use super::cli_print_fmt;

#[derive(Args)]
pub struct QualityArgs {
    #[arg(long)]
    pub format: Option<String>,
    #[arg(long)]
    pub verbose: bool,
    pub paths: Vec<String>,
}

pub fn run(args: QualityArgs) -> anyhow::Result<()> {
    cli_print_fmt(format!("[STUB] quality: format={:?} verbose={} paths={:?}", args.format, args.verbose, args.paths));
    Ok(())
}
