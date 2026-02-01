//! Memory RPC: JSON-RPC interface for memory operations.
//! Reads JSON from stdin, dispatches to memory bank operations.

use std::io::Write;

pub fn run() -> anyhow::Result<()> {
    let stdin = std::io::stdin();
    let input: Result<serde_json::Value, _> = serde_json::from_reader(stdin.lock());

    let req = match input {
        Ok(v) => v,
        Err(e) => {
            write_error(-32700, &format!("Parse error: {e}"), None)?;
            return Ok(());
        }
    };

    let method = req
        .get("method")
        .and_then(|v| v.as_str())
        .unwrap_or("");
    let id = req.get("id").cloned();

    match method {
        "memory.list" => {
            let mem_dir = memory_dir();
            let mut files = Vec::new();
            collect_files(&mem_dir, &mut files);
            let paths: Vec<String> = files
                .iter()
                .filter_map(|f| f.strip_prefix(&mem_dir).ok())
                .map(|f| f.to_string_lossy().to_string())
                .collect();
            write_result(&serde_json::json!({"files": paths}), id.as_ref())?;
        }
        "memory.read" => {
            let path = req
                .get("params")
                .and_then(|p| p.get("path"))
                .and_then(|v| v.as_str())
                .unwrap_or("");
            if path.is_empty() {
                write_error(-32602, "Missing params.path", id.as_ref())?;
            } else {
                let full = memory_dir().join(path);
                match std::fs::read_to_string(&full) {
                    Ok(content) => {
                        write_result(
                            &serde_json::json!({"path": path, "content": content}),
                            id.as_ref(),
                        )?;
                    }
                    Err(e) => {
                        write_error(-32000, &format!("Read failed: {e}"), id.as_ref())?;
                    }
                }
            }
        }
        "memory.status" => {
            let mem_dir = memory_dir();
            let mut count = 0;
            collect_files(&mem_dir, &mut Vec::new());
            count_files(&mem_dir, &mut count);
            write_result(
                &serde_json::json!({"path": mem_dir.to_string_lossy(), "files": count}),
                id.as_ref(),
            )?;
        }
        _ => {
            write_error(-32601, &format!("Method not found: {method}"), id.as_ref())?;
        }
    }

    Ok(())
}

fn memory_dir() -> std::path::PathBuf {
    let home = dirs::home_dir().unwrap_or_default();
    home.join(".local")
        .join("shared")
        .join("shared-ai")
        .join("memory")
}

fn collect_files(dir: &std::path::Path, files: &mut Vec<std::path::PathBuf>) {
    let Ok(entries) = std::fs::read_dir(dir) else {
        return;
    };
    for entry in entries.flatten() {
        let path = entry.path();
        if path.is_dir() {
            collect_files(&path, files);
        } else {
            files.push(path);
        }
    }
}

fn count_files(dir: &std::path::Path, count: &mut usize) {
    let Ok(entries) = std::fs::read_dir(dir) else {
        return;
    };
    for entry in entries.flatten() {
        let path = entry.path();
        if path.is_dir() {
            count_files(&path, count);
        } else {
            *count += 1;
        }
    }
}

fn write_result(result: &serde_json::Value, id: Option<&serde_json::Value>) -> anyhow::Result<()> {
    let resp = serde_json::json!({
        "jsonrpc": "2.0",
        "result": result,
        "id": id,
    });
    let stdout = std::io::stdout();
    let mut handle = stdout.lock();
    writeln!(handle, "{}", serde_json::to_string(&resp)?)?;
    Ok(())
}

fn write_error(code: i32, message: &str, id: Option<&serde_json::Value>) -> anyhow::Result<()> {
    let resp = serde_json::json!({
        "jsonrpc": "2.0",
        "error": {"code": code, "message": message},
        "id": id,
    });
    let stdout = std::io::stdout();
    let mut handle = stdout.lock();
    writeln!(handle, "{}", serde_json::to_string(&resp)?)?;
    Ok(())
}
