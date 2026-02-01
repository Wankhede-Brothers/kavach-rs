//! Lint: file linting for code quality issues.
//! Checks: trailing whitespace, line length, tab/space consistency, DACE line count.
//! Output: TOON (default) or JSON format.

use std::io::Write;
use std::path::Path;

use clap::Args;

#[derive(Args)]
pub struct LintArgs {
    #[arg(long)]
    pub fix: bool,
    #[arg(long)]
    pub format: Option<String>,
    pub files: Vec<String>,
}

struct LintResult {
    file: String,
    issues: Vec<LintIssue>,
    fixed: usize,
}

struct LintIssue {
    line: usize,
    code: String,
    message: String,
}

pub fn run(args: LintArgs) -> anyhow::Result<()> {
    if args.files.is_empty() {
        crate::commands::cli_print("Usage: kavach lint [--fix] [--format=toon|json] <file|dir>...");
        return Ok(());
    }

    let mut results = Vec::new();

    for arg in &args.files {
        let path = Path::new(arg);
        if !path.exists() {
            let stderr = std::io::stderr();
            let mut h = stderr.lock();
            let _ = writeln!(h, "error: {arg}: not found");
            continue;
        }

        if path.is_dir() {
            walk_dir(path, &mut results, args.fix);
        } else {
            let result = lint_file(path, args.fix);
            results.push(result);
        }
    }

    let fmt = args.format.as_deref().unwrap_or("toon");
    if fmt == "json" {
        output_json(&results)
    } else {
        output_toon(&results)
    }
}

fn walk_dir(dir: &Path, results: &mut Vec<LintResult>, fix: bool) {
    let entries = match std::fs::read_dir(dir) {
        Ok(e) => e,
        Err(_) => return,
    };

    for entry in entries.flatten() {
        let path = entry.path();
        if path.is_dir() {
            // Skip hidden dirs and target/node_modules
            let name = path.file_name().unwrap_or_default().to_string_lossy();
            if name.starts_with('.') || name == "target" || name == "node_modules" {
                continue;
            }
            walk_dir(&path, results, fix);
        } else if is_lintable(&path) {
            let result = lint_file(&path, fix);
            if !result.issues.is_empty() {
                results.push(result);
            }
        }
    }
}

fn is_lintable(path: &Path) -> bool {
    let ext = path.extension().and_then(|e| e.to_str()).unwrap_or("");
    matches!(ext, "go" | "rs" | "ts" | "tsx" | "js" | "jsx" | "py" | "json" | "yaml" | "yml" | "toon" | "md")
}

fn lint_file(path: &Path, fix: bool) -> LintResult {
    let file_str = path.to_string_lossy().to_string();
    let mut result = LintResult {
        file: file_str,
        issues: Vec::new(),
        fixed: 0,
    };

    let content = match std::fs::read_to_string(path) {
        Ok(c) => c,
        Err(e) => {
            result.issues.push(LintIssue {
                line: 0, code: "E000".into(),
                message: format!("cannot read file: {e}"),
            });
            return result;
        }
    };

    let lines: Vec<&str> = content.lines().collect();
    let ext = path.extension().and_then(|e| e.to_str()).unwrap_or("");

    // W001: trailing whitespace
    for (i, line) in lines.iter().enumerate() {
        if line.ends_with(' ') || line.ends_with('\t') {
            result.issues.push(LintIssue {
                line: i + 1, code: "W001".into(),
                message: "trailing whitespace".into(),
            });
        }
    }

    // W002: line too long (>120)
    for (i, line) in lines.iter().enumerate() {
        if line.len() > 120 {
            result.issues.push(LintIssue {
                line: i + 1, code: "W002".into(),
                message: format!("line too long ({} > 120)", line.len()),
            });
        }
    }

    // W003: Go tabs vs spaces
    if ext == "go" {
        for (i, line) in lines.iter().enumerate() {
            if line.starts_with("    ") && !line.starts_with('\t') {
                result.issues.push(LintIssue {
                    line: i + 1, code: "W003".into(),
                    message: "use tabs instead of spaces for Go indentation".into(),
                });
            }
        }
    }

    // D001: DACE line count
    if lines.len() > 100 {
        result.issues.push(LintIssue {
            line: 1, code: "D001".into(),
            message: format!("DACE: file exceeds 100 lines ({} lines)", lines.len()),
        });
    }

    // Auto-fix trailing whitespace
    if fix && result.issues.iter().any(|i| i.code == "W001") {
        let new_content: String = lines.iter()
            .map(|l| l.trim_end())
            .collect::<Vec<&str>>()
            .join("\n");
        if new_content != content {
            let fixed_count = result.issues.iter().filter(|i| i.code == "W001").count();
            let _ = std::fs::write(path, &new_content);
            result.fixed = fixed_count;
        }
    }

    result
}

fn output_toon(results: &[LintResult]) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();

    let total_issues: usize = results.iter().map(|r| r.issues.len()).sum();

    writeln!(w, "[LINT_RESULTS]")?;
    writeln!(w, "files: {}", results.len())?;
    writeln!(w, "issues: {total_issues}")?;
    writeln!(w)?;

    for r in results {
        if r.issues.is_empty() { continue; }
        writeln!(w, "[FILE:{}]", r.file)?;
        for issue in &r.issues {
            writeln!(w, "  {}: [{}] {}", issue.line, issue.code, issue.message)?;
        }
        if r.fixed > 0 {
            writeln!(w, "  fixed: {}", r.fixed)?;
        }
        writeln!(w)?;
    }

    Ok(())
}

fn output_json(results: &[LintResult]) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[")?;
    for (i, r) in results.iter().enumerate() {
        let comma = if i < results.len() - 1 { "," } else { "" };
        writeln!(w, "  {{\"file\": {:?}, \"issues\": {}}}{comma}", r.file, r.issues.len())?;
    }
    writeln!(w, "]")?;

    Ok(())
}
