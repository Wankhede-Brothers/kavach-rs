//! Memory bank: query and manage TOON-based memory.
//! Path: ~/.local/shared/shared-ai/memory/
//! DACE: project-scoped by default (active + global only)

use std::io::Write;
use std::path::{Path, PathBuf};

use clap::Args;

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
    if args.status {
        show_memory_status()
    } else if args.scan {
        scan_memory_bank()
    } else if args.all {
        show_all_projects()
    } else {
        show_project_scoped()
    }
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

const CATEGORIES: &[&str] = &[
    "decisions", "graph", "kanban", "patterns",
    "proposals", "research", "roadmaps", "STM",
];

const PROJECT_CATEGORIES: &[&str] = &[
    "decisions", "kanban", "patterns", "proposals", "research", "roadmaps",
];

fn show_memory_status() -> anyhow::Result<()> {
    let mem = memory_dir();
    let project = detect_project();
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[MEMORY_BANK]")?;
    writeln!(w, "path: ~/.local/shared/shared-ai/memory/")?;
    writeln!(w, "project: {project}")?;
    writeln!(w)?;

    writeln!(w, "[CATEGORIES]")?;
    for cat in CATEGORIES {
        let path = mem.join(cat);
        let count = count_toon_recursive(&path);
        writeln!(w, "{cat}: {count}")?;
    }
    writeln!(w)?;

    writeln!(w, "[PROJECT_FILES]")?;
    for cat in PROJECT_CATEGORIES {
        let pf = mem.join(cat).join(&project).join(format!("{cat}.toon"));
        let status = if pf.exists() { "OK" } else { "-" };
        writeln!(w, "{cat}: {status}")?;
    }
    writeln!(w)?;

    writeln!(w, "[ROOT_FILES]")?;
    for f in &["GOVERNANCE.toon", "index.toon", "volatile.toon"] {
        let status = if mem.join(f).exists() { "OK" } else { "-" };
        let name = f.trim_end_matches(".toon");
        writeln!(w, "{name}: {status}")?;
    }

    Ok(())
}

fn show_project_scoped() -> anyhow::Result<()> {
    let mem = memory_dir();
    let project = detect_project();
    let out = std::io::stdout();
    let mut w = out.lock();

    let wd = std::env::current_dir()
        .map(|p| p.to_string_lossy().to_string())
        .unwrap_or_default();

    writeln!(w, "[MEMORY]")?;
    writeln!(w, "path: ~/.local/shared/shared-ai/memory/")?;
    writeln!(w, "project: {project}")?;
    writeln!(w, "workdir: {wd}")?;
    writeln!(w, "scope: PROJECT_ISOLATED (active + global only)")?;
    writeln!(w)?;

    writeln!(w, "[DETECTION]")?;
    if std::env::var("KAVACH_PROJECT").is_ok() {
        writeln!(w, "method: KAVACH_PROJECT env var")?;
    } else if std::env::current_dir().map(|d| d.join(".git").exists()).unwrap_or(false) {
        writeln!(w, "method: .git root detection")?;
    } else {
        writeln!(w, "method: fallback (global)")?;
    }
    writeln!(w)?;

    writeln!(w, "[PROJECT_DOCS]")?;
    let mut project_total = 0;
    for cat in PROJECT_CATEGORIES {
        let path = mem.join(cat).join(&project);
        let count = count_toon_in_dir(&path);
        project_total += count;
        if count > 0 {
            writeln!(w, "{cat}: {count}")?;
        }
    }
    writeln!(w, "project_total: {project_total}")?;
    writeln!(w)?;

    writeln!(w, "[GLOBAL_DOCS]")?;
    let mut global_total = 0;
    for cat in PROJECT_CATEGORIES {
        let path = mem.join(cat).join("global");
        let count = count_toon_in_dir(&path);
        global_total += count;
        if count > 0 {
            writeln!(w, "{cat}: {count}")?;
        }
    }
    writeln!(w, "global_total: {global_total}")?;
    writeln!(w)?;

    let stm_path = mem.join("STM");
    let stm_count = count_toon_in_dir(&stm_path);
    writeln!(w, "[STM]")?;
    writeln!(w, "files: {stm_count}")?;
    writeln!(w)?;

    writeln!(w, "[ROOT]")?;
    for f in &["GOVERNANCE.toon", "index.toon", "volatile.toon"] {
        if mem.join(f).exists() {
            let name = f.trim_end_matches(".toon");
            writeln!(w, "{name}: OK")?;
        }
    }
    writeln!(w)?;

    let total = project_total + global_total + stm_count;
    writeln!(w, "[TOTAL] {total} (project: {project_total}, global: {global_total}, stm: {stm_count})")?;

    Ok(())
}

fn show_all_projects() -> anyhow::Result<()> {
    let mem = memory_dir();
    let project = detect_project();
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[MEMORY]")?;
    writeln!(w, "path: ~/.local/shared/shared-ai/memory/")?;
    writeln!(w, "active_project: {project}")?;
    writeln!(w, "scope: ALL_PROJECTS (--all flag)")?;
    writeln!(w)?;

    writeln!(w, "[DOCS]")?;
    let mut total = 0;
    for cat in CATEGORIES {
        let count = count_toon_recursive(&mem.join(cat));
        total += count;
        writeln!(w, "{cat}: {count}")?;
    }
    writeln!(w, "total: {total}")?;
    writeln!(w)?;

    writeln!(w, "[PROJECTS]")?;
    let kanban_dir = mem.join("kanban");
    if let Ok(entries) = std::fs::read_dir(&kanban_dir) {
        for entry in entries.flatten() {
            if entry.path().is_dir() {
                let name = entry.file_name().to_string_lossy().to_string();
                let indicator = if name == project {
                    " <- ACTIVE"
                } else if name == "global" {
                    " (shared)"
                } else {
                    ""
                };
                writeln!(w, "- {name}{indicator}")?;
            }
        }
    }
    writeln!(w)?;

    writeln!(w, "[ROOT]")?;
    for f in &["GOVERNANCE.toon", "index.toon", "volatile.toon"] {
        let name = f.trim_end_matches(".toon");
        let status = if mem.join(f).exists() { "OK" } else { "-" };
        writeln!(w, "{name}: {status}")?;
    }

    Ok(())
}

fn scan_memory_bank() -> anyhow::Result<()> {
    let mem = memory_dir();
    let project = detect_project();
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[SCAN]")?;
    writeln!(w, "path: ~/.local/shared/shared-ai/memory/")?;
    writeln!(w, "project: {project}")?;
    writeln!(w)?;

    writeln!(w, "[INDEX]")?;
    let mut total = 0;
    for cat in CATEGORIES {
        let count = count_toon_recursive(&mem.join(cat));
        total += count;
        writeln!(w, "{cat}: {count}")?;
    }
    writeln!(w, "total: {total}")?;
    writeln!(w)?;

    writeln!(w, "[PROJECTS]")?;
    let kanban_dir = mem.join("kanban");
    if let Ok(entries) = std::fs::read_dir(&kanban_dir) {
        for entry in entries.flatten() {
            if entry.path().is_dir() {
                let name = entry.file_name().to_string_lossy().to_string();
                if name != "global" {
                    let indicator = if name == project { " (current)" } else { "" };
                    writeln!(w, "- {name}{indicator}")?;
                }
            }
        }
    }

    Ok(())
}

fn count_toon_in_dir(dir: &Path) -> usize {
    let Ok(entries) = std::fs::read_dir(dir) else {
        return 0;
    };
    entries
        .flatten()
        .filter(|e| {
            !e.path().is_dir()
                && e.path()
                    .extension()
                    .map(|ext| ext == "toon")
                    .unwrap_or(false)
        })
        .count()
}

fn count_toon_recursive(dir: &Path) -> usize {
    if !dir.exists() {
        return 0;
    }
    let mut count = 0;
    let Ok(entries) = std::fs::read_dir(dir) else {
        return 0;
    };
    for entry in entries.flatten() {
        let path = entry.path();
        if path.is_dir() {
            count += count_toon_recursive(&path);
        } else if path.extension().map(|e| e == "toon").unwrap_or(false) {
            count += 1;
        }
    }
    count
}
