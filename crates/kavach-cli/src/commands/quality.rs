//! Quality: code quality analysis with DACE scoring.
//! Metrics: lines, functions, imports, DACE score, complexity.
//! Output: TOON (default) or JSON format.

use std::io::Write;
use std::path::Path;

use clap::Args;

#[derive(Args)]
pub struct QualityArgs {
    #[arg(long)]
    pub format: Option<String>,
    #[arg(long)]
    pub verbose: bool,
    pub paths: Vec<String>,
}

struct QualityResult {
    file: String,
    lines: usize,
    functions: usize,
    imports: usize,
    dace_score: i32,
    complexity: String,
}

pub fn run(args: QualityArgs) -> anyhow::Result<()> {
    if args.paths.is_empty() {
        crate::commands::cli_print("Usage: kavach quality [--verbose] [--format=toon|json] <file|dir>...");
        return Ok(());
    }

    let mut results = Vec::new();

    for arg in &args.paths {
        let path = Path::new(arg);
        if !path.exists() {
            let stderr = std::io::stderr();
            let mut h = stderr.lock();
            let _ = writeln!(h, "error: {arg}: not found");
            continue;
        }

        if path.is_dir() {
            walk_dir(path, &mut results);
        } else {
            results.push(analyze_file(path));
        }
    }

    let fmt = args.format.as_deref().unwrap_or("toon");
    if fmt == "json" {
        output_json(&results)
    } else {
        output_toon(&results, args.verbose)
    }
}

fn walk_dir(dir: &Path, results: &mut Vec<QualityResult>) {
    let entries = match std::fs::read_dir(dir) {
        Ok(e) => e,
        Err(_) => return,
    };

    for entry in entries.flatten() {
        let path = entry.path();
        if path.is_dir() {
            let name = path.file_name().unwrap_or_default().to_string_lossy();
            if name.starts_with('.') || name == "target" || name == "node_modules" {
                continue;
            }
            walk_dir(&path, results);
        } else if is_analyzable(&path) {
            results.push(analyze_file(&path));
        }
    }
}

fn is_analyzable(path: &Path) -> bool {
    let ext = path.extension().and_then(|e| e.to_str()).unwrap_or("");
    matches!(ext, "go" | "rs" | "ts" | "tsx" | "js" | "jsx" | "py")
}

fn analyze_file(path: &Path) -> QualityResult {
    let file_str = path.to_string_lossy().to_string();
    let content = match std::fs::read_to_string(path) {
        Ok(c) => c,
        Err(_) => {
            return QualityResult {
                file: file_str, lines: 0, functions: 0,
                imports: 0, dace_score: 0, complexity: "unknown".into(),
            };
        }
    };

    let lines: Vec<&str> = content.lines().collect();
    let ext = path.extension().and_then(|e| e.to_str()).unwrap_or("");

    let functions = count_functions(&lines, ext);
    let imports = count_imports(&lines, ext);

    let mut result = QualityResult {
        file: file_str,
        lines: lines.len(),
        functions,
        imports,
        dace_score: 100,
        complexity: String::new(),
    };

    // DACE score
    if result.lines > 100 {
        let deduct = ((result.lines - 100) / 10).min(50) as i32;
        result.dace_score -= deduct;
    }
    if result.functions > 10 {
        let deduct = ((result.functions - 10) * 2).min(20) as i32;
        result.dace_score -= deduct;
    }
    if result.dace_score < 0 { result.dace_score = 0; }

    // Complexity
    result.complexity = if result.lines <= 50 && result.functions <= 5 {
        "low".into()
    } else if result.lines <= 100 && result.functions <= 10 {
        "medium".into()
    } else {
        "high".into()
    };

    result
}

fn count_functions(lines: &[&str], ext: &str) -> usize {
    let mut count = 0;
    for line in lines {
        let trimmed = line.trim();
        match ext {
            "go" => {
                if trimmed.starts_with("func ") { count += 1; }
            }
            "rs" => {
                if trimmed.starts_with("fn ") || trimmed.starts_with("pub fn ") { count += 1; }
            }
            "ts" | "tsx" | "js" | "jsx" => {
                if trimmed.contains("function ") || trimmed.contains("=> {") || trimmed.contains("async ") {
                    count += 1;
                }
            }
            "py" => {
                if trimmed.starts_with("def ") || trimmed.starts_with("async def ") { count += 1; }
            }
            _ => {}
        }
    }
    count
}

fn count_imports(lines: &[&str], ext: &str) -> usize {
    let mut count = 0;
    for line in lines {
        let trimmed = line.trim();
        match ext {
            "go" => {
                if trimmed.starts_with("import ") || trimmed == "import (" { count += 1; }
                if trimmed.starts_with('"') && trimmed.ends_with('"') { count += 1; }
            }
            "rs" => {
                if trimmed.starts_with("use ") { count += 1; }
            }
            "ts" | "tsx" | "js" | "jsx" => {
                if trimmed.starts_with("import ") || trimmed.starts_with("require(") { count += 1; }
            }
            "py" => {
                if trimmed.starts_with("import ") || trimmed.starts_with("from ") { count += 1; }
            }
            _ => {}
        }
    }
    count
}

fn output_toon(results: &[QualityResult], verbose: bool) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();

    let total_lines: usize = results.iter().map(|r| r.lines).sum();
    let total_functions: usize = results.iter().map(|r| r.functions).sum();
    let avg_dace = if results.is_empty() { 0 } else {
        results.iter().map(|r| r.dace_score).sum::<i32>() / results.len() as i32
    };

    writeln!(w, "[QUALITY_SUMMARY]")?;
    writeln!(w, "files: {}", results.len())?;
    writeln!(w, "total_lines: {total_lines}")?;
    writeln!(w, "total_functions: {total_functions}")?;
    writeln!(w, "avg_dace_score: {avg_dace}")?;
    writeln!(w)?;

    if verbose {
        for r in results {
            writeln!(w, "[FILE:{}]", r.file)?;
            writeln!(w, "  lines: {}", r.lines)?;
            writeln!(w, "  functions: {}", r.functions)?;
            writeln!(w, "  imports: {}", r.imports)?;
            writeln!(w, "  dace_score: {}", r.dace_score)?;
            writeln!(w, "  complexity: {}", r.complexity)?;
            writeln!(w)?;
        }
    }

    Ok(())
}

fn output_json(results: &[QualityResult]) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[")?;
    for (i, r) in results.iter().enumerate() {
        let comma = if i < results.len() - 1 { "," } else { "" };
        writeln!(w, "  {{\"file\": {:?}, \"lines\": {}, \"dace_score\": {}, \"complexity\": {:?}}}{comma}",
            r.file, r.lines, r.dace_score, r.complexity)?;
    }
    writeln!(w, "]")?;

    Ok(())
}
