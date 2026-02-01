use crate::commands::cli_print_fmt;
use super::read_hook_stdin;

pub fn run(hook: bool) -> anyhow::Result<()> {
    if hook {
        let _input = read_hook_stdin().ok();
    }
    cli_print_fmt(format!("[STUB] gates subagent: hook={hook}"));
    Ok(())
}
