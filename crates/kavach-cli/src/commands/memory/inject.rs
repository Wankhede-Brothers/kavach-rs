//! Memory inject: injects context from memory bank into hook responses.
//! Used by RPC and hook pipelines to provide DACE-optimized context.

use std::io::Write;
use std::path::PathBuf;

pub fn run() -> anyhow::Result<()> {
    let mem_dir = memory_dir();
    let out = std::io::stdout();
    let mut w = out.lock();

    if !mem_dir.exists() {
        writeln!(w, "[INJECT:EMPTY]")?;
        writeln!(w, "status: no memory bank found")?;
        return Ok(());
    }

    let project = detect_project();
    let mut injected = Vec::new();

    // Inject project-specific decisions
    let decisions_dir = mem_dir.join("decisions").join(&project);
    if decisions_dir.exists() {
        if let Ok(entries) = std::fs::read_dir(&decisions_dir) {
            for entry in entries.flatten() {
                let path = entry.path();
                if path.extension().map(|e| e == "toon").unwrap_or(false) {
                    if let Ok(content) = std::fs::read_to_string(&path) {
                        let name = path.file_stem().unwrap_or_default().to_string_lossy();
                        injected.push(format!("[DECISION:{name}]\n{content}"));
                    }
                }
            }
        }
    }

    // Inject project patterns
    let patterns_path = mem_dir.join("patterns").join(&project).join("patterns.toon");
    if patterns_path.exists() {
        if let Ok(content) = std::fs::read_to_string(&patterns_path) {
            injected.push(format!("[PATTERNS]\n{content}"));
        }
    }

    writeln!(w, "[INJECT:CONTEXT]")?;
    writeln!(w, "project: {project}")?;
    writeln!(w, "sections: {}", injected.len())?;
    writeln!(w)?;

    for section in &injected {
        writeln!(w, "{section}")?;
        writeln!(w)?;
    }

    Ok(())
}

fn memory_dir() -> PathBuf {
    let home = dirs::home_dir().unwrap_or_default();
    home.join(".local")
        .join("shared")
        .join("shared-ai")
        .join("memory")
}

fn detect_project() -> String {
    if let Ok(wd) = std::env::current_dir() {
        if wd.join(".git").exists() {
            return wd
                .file_name()
                .map(|n| n.to_string_lossy().to_string())
                .unwrap_or_default();
        }
    }
    String::new()
}
