use crate::commands::cli_print_fmt;

pub fn run() -> anyhow::Result<()> {
    let stdin = std::io::stdin();
    let input: Result<serde_json::Value, _> = serde_json::from_reader(stdin.lock());
    match input {
        Ok(v) => cli_print_fmt(format!("[STUB] memory rpc: received {:?}", v.get("method"))),
        Err(_) => cli_print_fmt("[STUB] memory rpc: no valid JSON on stdin".into()),
    }
    Ok(())
}
