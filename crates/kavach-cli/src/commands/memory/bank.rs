use clap::Args;
use crate::commands::cli_print_fmt;

#[derive(Args)]
pub struct BankArgs {
    #[arg(long)]
    pub status: bool,
    #[arg(long)]
    pub scan: bool,
    #[arg(long)]
    pub all: bool,
}

pub fn run(args: BankArgs) -> anyhow::Result<()> {
    cli_print_fmt(format!("[STUB] memory bank: status={} scan={} all={}", args.status, args.scan, args.all));
    Ok(())
}
