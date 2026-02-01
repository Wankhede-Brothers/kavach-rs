//! Memory view: browse memory bank contents with filtering.
//! Lists TOON files from ~/.local/shared/shared-ai/memory/

use std::io::Write;
use std::path::PathBuf;

pub fn run() -> anyhow::Result<()> {
    let mem_dir = memory_dir();
    let out = std::io::stdout();
    let mut w = out.lock();

    if !mem_dir.exists() {
        writeln!(w, "[MEMORY:VIEW]")?;
        writeln!(w, "status: empty")?;
        writeln!(w, "path: {}", mem_dir.display())?;
        return Ok(());
    }

    let mut files = Vec::new();
    collect_toon_files(&mem_dir, &mut files);
    files.sort();

    writeln!(w, "[MEMORY:VIEW]")?;
    writeln!(w, "path: {}", mem_dir.display())?;
    writeln!(w, "total: {}", files.len())?;
    writeln!(w)?;

    for file in &files {
        let rel = file
            .strip_prefix(&mem_dir)
            .unwrap_or(file)
            .to_string_lossy();
        let size = std::fs::metadata(file)
            .map(|m| m.len())
            .unwrap_or(0);
        writeln!(w, "  {rel} ({size}b)")?;
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

fn collect_toon_files(dir: &std::path::Path, files: &mut Vec<PathBuf>) {
    let Ok(entries) = std::fs::read_dir(dir) else {
        return;
    };
    for entry in entries.flatten() {
        let path = entry.path();
        if path.is_dir() {
            collect_toon_files(&path, files);
        } else if path.extension().map(|e| e == "toon").unwrap_or(false) {
            files.push(path);
        }
    }
}
