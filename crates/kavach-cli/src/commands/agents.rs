use clap::Args;
use super::cli_print_fmt;

#[derive(Args)]
pub struct AgentsArgs {
    #[arg(long)]
    pub get: Option<String>,
    #[arg(long)]
    pub sutra: bool,
    #[arg(long)]
    pub inject: bool,
}

pub fn run(args: AgentsArgs) -> anyhow::Result<()> {
    if let Some(name) = &args.get {
        cli_print_fmt(format!("[STUB] agents --get {name}"));
    } else if args.sutra {
        cli_print_fmt("[STUB] agents --sutra".into());
    } else if args.inject {
        cli_print_fmt("[STUB] agents --inject".into());
    } else {
        cli_print_fmt("[STUB] agents: list all agents".into());
    }
    Ok(())
}
