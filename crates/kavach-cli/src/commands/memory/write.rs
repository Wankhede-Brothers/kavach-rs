//! Memory write: write TOON content to memory bank (project-isolated).
//! Input: stdin (TOON-formatted content)
//! Path: ~/.local/shared/shared-ai/memory/{category}/{project}/{key}.toon

use std::io::{self, Read, Write};
use std::path::PathBuf;

use clap::Args;

#[derive(Args)]
pub struct WriteArgs {
    #[arg(short = 'c', long)]
    pub category: Option<String>,
    #[arg(short = 'k', long)]
    pub key: Option<String>,
    #[arg(short = 'p', long)]
    pub project: Option<String>,
}

pub fn run(args: WriteArgs) -> anyhow::Result<()> {
    let category = match &args.category {
        Some(c) => c.clone(),
        None => {
            let stderr = io::stderr();
            let mut h = stderr.lock();
            let _ = writeln!(h, "Error: --category required");
            return Ok(());
        }
    };

    let key = match &args.key {
        Some(k) => k.clone(),
        None => {
            let stderr = io::stderr();
            let mut h = stderr.lock();
            let _ = writeln!(h, "Error: --key required");
            return Ok(());
        }
    };

    // Read content from stdin
    let mut content = String::new();
    io::stdin().read_to_string(&mut content)?;

    if content.trim().is_empty() {
        let stderr = io::stderr();
        let mut h = stderr.lock();
        let _ = writeln!(h, "Error: no content provided on stdin");
        return Ok(());
    }

    // Determine project
    let project = args.project.unwrap_or_else(detect_project);

    // Valid categories
    let valid = ["decisions", "graph", "kanban", "patterns", "proposals", "research", "roadmaps"];
    if !valid.contains(&category.as_str()) {
        let stderr = io::stderr();
        let mut h = stderr.lock();
        let _ = writeln!(h, "Error: invalid category '{}'. Valid: {}", category, valid.join(", "));
        return Ok(());
    }

    // Build path and ensure directory exists
    let mem = memory_dir();
    let project_dir = mem.join(&category).join(&project);
    std::fs::create_dir_all(&project_dir)?;

    let path = project_dir.join(format!("{key}.toon"));

    // Write content
    std::fs::write(&path, &content)?;

    let out = io::stdout();
    let mut w = out.lock();
    writeln!(w, "Saved to {}", path.display())?;

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
    if let Ok(val) = std::env::var("KAVACH_PROJECT") {
        if !val.is_empty() {
            return val;
        }
    }
    if let Ok(wd) = std::env::current_dir() {
        if wd.join(".git").exists() {
            return wd
                .file_name()
                .map(|n| n.to_string_lossy().to_string())
                .unwrap_or_else(|| "global".into());
        }
    }
    "global".into()
}
