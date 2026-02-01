//! Memory spec: loads and injects SDD (Specs Driven Development) specs.
//! Reads .spec.toon files from project directory for pre-tool context.

use std::io::Write;

pub fn run() -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();

    let specs = discover_specs();

    writeln!(w, "[SPEC:INJECT]")?;
    writeln!(w, "total: {}", specs.len())?;
    writeln!(w)?;

    if specs.is_empty() {
        writeln!(w, "status: no specs found")?;
        writeln!(
            w,
            "hint: Create .spec.toon files in project root or specs/ directory"
        )?;
        return Ok(());
    }

    for (name, content) in &specs {
        writeln!(w, "[SPEC:{name}]")?;
        // Output first 20 lines max (DACE budget)
        for line in content.lines().take(20) {
            writeln!(w, "  {line}")?;
        }
        let total_lines = content.lines().count();
        if total_lines > 20 {
            writeln!(w, "  ... ({} more lines)", total_lines - 20)?;
        }
        writeln!(w)?;
    }

    Ok(())
}

fn discover_specs() -> Vec<(String, String)> {
    let mut specs = Vec::new();

    let wd = std::env::current_dir().unwrap_or_default();

    // Check project root for .spec.toon files
    scan_dir_for_specs(&wd, &mut specs);

    // Check specs/ directory
    let specs_dir = wd.join("specs");
    if specs_dir.exists() {
        scan_dir_for_specs(&specs_dir, &mut specs);
    }

    // Check .claude/specs/ directory
    let claude_specs = wd.join(".claude").join("specs");
    if claude_specs.exists() {
        scan_dir_for_specs(&claude_specs, &mut specs);
    }

    specs
}

fn scan_dir_for_specs(dir: &std::path::Path, specs: &mut Vec<(String, String)>) {
    let Ok(entries) = std::fs::read_dir(dir) else {
        return;
    };
    for entry in entries.flatten() {
        let path = entry.path();
        let name = path
            .file_name()
            .unwrap_or_default()
            .to_string_lossy()
            .to_string();
        if name.ends_with(".spec.toon") || name.ends_with(".spec.md") {
            if let Ok(content) = std::fs::read_to_string(&path) {
                let spec_name = name
                    .trim_end_matches(".spec.toon")
                    .trim_end_matches(".spec.md")
                    .to_string();
                specs.push((spec_name, content));
            }
        }
    }
}
