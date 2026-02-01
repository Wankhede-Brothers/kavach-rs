use super::cli_print_fmt;

pub fn run() -> anyhow::Result<()> {
    let now = chrono::Local::now().format("%Y-%m-%d %H:%M:%S %Z");
    cli_print_fmt(format!("META:STATUS\n  protocol: SP/1.0\n  binary: kavach-cli (rust)\n  today: {now}\n  status: STUB"));
    Ok(())
}
