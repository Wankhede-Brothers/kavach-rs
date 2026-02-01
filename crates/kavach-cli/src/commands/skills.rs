use clap::Args;
use super::cli_print_fmt;

#[derive(Args)]
pub struct SkillsArgs {
    #[arg(long)]
    pub get: Option<String>,
    #[arg(long)]
    pub sutra: bool,
    #[arg(long)]
    pub inject: bool,
}

pub fn run(args: SkillsArgs) -> anyhow::Result<()> {
    if let Some(name) = &args.get {
        cli_print_fmt(format!("[STUB] skills --get {name}"));
    } else if args.sutra {
        cli_print_fmt("[STUB] skills --sutra".into());
    } else if args.inject {
        cli_print_fmt("[STUB] skills --inject".into());
    } else {
        cli_print_fmt("[STUB] skills: list all skills".into());
    }
    Ok(())
}
