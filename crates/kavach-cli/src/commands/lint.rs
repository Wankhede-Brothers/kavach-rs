use clap::Args;
use super::cli_print_fmt;

#[derive(Args)]
pub struct LintArgs {
    #[arg(long)]
    pub fix: bool,
    #[arg(long)]
    pub format: Option<String>,
    pub files: Vec<String>,
}

pub fn run(args: LintArgs) -> anyhow::Result<()> {
    cli_print_fmt(format!("[STUB] lint: fix={} format={:?} files={:?}", args.fix, args.format, args.files));
    Ok(())
}
