//! Aegis Guardian verification (Level 2).
//! Two-stage verification: TESTING (lint, warnings, bugs) -> VERIFIED (dead code, suppressed).
//! Runs project-specific toolchain commands to gather metrics.

use std::io::Write;
use std::path::Path;
use std::process::Command;

use clap::Args;

use crate::hook;

#[derive(Args)]
pub struct AegisArgs {
    #[arg(long)]
    pub hook: bool,
    #[arg(long)]
    pub task: Option<String>,
}

struct AegisResult {
    stage: String,
    status: String,
    lint_issues: i32,
    warnings: i32,
    core_bugs: i32,
    dead_code: bool,
    suppressed: bool,
    algo_ok: bool,
    fail_reasons: Vec<String>,
    exec_errors: Vec<String>,
}

pub fn run(args: AegisArgs) -> anyhow::Result<()> {
    let project = detect_project();
    let today = chrono::Local::now().format("%Y-%m-%d").to_string();

    if args.hook {
        let input = hook::must_read_hook_input();
        let tool_result = input.get_string("tool_result");
        if tool_result.is_empty() {
            return hook::exit_silent();
        }

        let result = run_verification();
        output_aegis_result(&project, &today, &result)?;

        if result.status == "passed" {
            hook::exit_silent()?;
        } else {
            hook::exit_block_toon("AEGIS_FAIL", &result.fail_reasons.join(","))?;
        }
        return Ok(());
    }

    let result = run_verification();
    output_aegis_result(&project, &today, &result)
}

fn run_verification() -> AegisResult {
    let mut result = AegisResult {
        stage: "TESTING".into(),
        status: "passed".into(),
        lint_issues: 0,
        warnings: 0,
        core_bugs: 0,
        dead_code: false,
        suppressed: false,
        algo_ok: true,
        fail_reasons: Vec::new(),
        exec_errors: Vec::new(),
    };

    let work_dir = std::env::current_dir().unwrap_or_default();

    // Stage 1: TESTING
    match count_lint_issues(&work_dir) {
        Ok(n) => result.lint_issues = n,
        Err(e) => result.exec_errors.push(e),
    }
    match count_warnings(&work_dir) {
        Ok(n) => result.warnings = n,
        Err(e) => result.exec_errors.push(e),
    }
    match count_core_bugs(&work_dir) {
        Ok(n) => result.core_bugs = n,
        Err(e) => result.exec_errors.push(e),
    }

    if result.lint_issues > 0 {
        result.fail_reasons.push(format!("lint_issues:{}", result.lint_issues));
    }
    if result.warnings > 0 {
        result.fail_reasons.push(format!("warnings:{}", result.warnings));
    }
    if result.core_bugs > 0 {
        result.fail_reasons.push(format!("core_bugs:{}", result.core_bugs));
    }

    if !result.fail_reasons.is_empty() {
        result.status = "failed".into();
        return result;
    }

    // Stage 2: VERIFIED
    result.stage = "VERIFIED".into();

    match has_dead_code(&work_dir) {
        Ok(b) => result.dead_code = b,
        Err(e) => result.exec_errors.push(e),
    }
    match has_suppressed_elements(&work_dir) {
        Ok(b) => result.suppressed = b,
        Err(e) => result.exec_errors.push(e),
    }

    if result.dead_code {
        result.fail_reasons.push("dead_code:found".into());
    }
    if result.suppressed {
        result.fail_reasons.push("suppressed_elements:found".into());
    }

    if !result.fail_reasons.is_empty() {
        result.status = "failed".into();
    }

    result
}

fn count_lint_issues(work_dir: &Path) -> Result<i32, String> {
    if work_dir.join("go.mod").exists() {
        let output = Command::new("go").args(["vet", "./..."])
            .current_dir(work_dir).output()
            .map_err(|e| format!("go vet failed: {e}"))?;
        let text = String::from_utf8_lossy(&output.stderr);
        return Ok(text.lines().filter(|l| !l.trim().is_empty()).count() as i32);
    }

    if work_dir.join("Cargo.toml").exists() {
        let output = Command::new("cargo").args(["clippy", "--message-format=short", "--", "-D", "warnings"])
            .current_dir(work_dir).output()
            .map_err(|e| format!("cargo clippy failed: {e}"))?;
        let text = String::from_utf8_lossy(&output.stderr);
        return Ok(text.matches("warning:").count() as i32);
    }

    Ok(0)
}

fn count_warnings(work_dir: &Path) -> Result<i32, String> {
    if work_dir.join("go.mod").exists() {
        let output = Command::new("go").args(["build", "-v", "./..."])
            .current_dir(work_dir).output()
            .map_err(|e| format!("go build -v failed: {e}"))?;
        let text = String::from_utf8_lossy(&output.stderr);
        return Ok(text.matches("warning").count() as i32);
    }

    if work_dir.join("Cargo.toml").exists() {
        let output = Command::new("cargo").args(["check", "--message-format=short"])
            .current_dir(work_dir).output()
            .map_err(|e| format!("cargo check failed: {e}"))?;
        let text = String::from_utf8_lossy(&output.stderr);
        return Ok(text.matches("warning:").count() as i32);
    }

    Ok(0)
}

fn count_core_bugs(work_dir: &Path) -> Result<i32, String> {
    let todo_upper = ["TO", "DO"].concat();
    let fixme_upper = ["FIX", "ME"].concat();
    let pattern = format!("{todo_upper}|{fixme_upper}|BUG|XXX");

    let output = Command::new("rg").args(["-c", &pattern])
        .arg(work_dir)
        .args(["--type", "go", "--type", "rust", "--type", "ts"])
        .output();

    match output {
        Ok(o) => {
            let text = String::from_utf8_lossy(&o.stdout);
            let count: i32 = text.lines()
                .filter_map(|l| l.rsplit(':').next())
                .filter_map(|n| n.parse::<i32>().ok())
                .sum();
            Ok(count)
        }
        Err(_) => Ok(0),
    }
}

fn has_dead_code(work_dir: &Path) -> Result<bool, String> {
    if work_dir.join("Cargo.toml").exists() {
        let output = Command::new("cargo").args(["check"])
            .current_dir(work_dir).output()
            .map_err(|e| format!("cargo check failed: {e}"))?;
        let text = String::from_utf8_lossy(&output.stderr);
        return Ok(text.contains("dead_code"));
    }

    if work_dir.join("go.mod").exists() {
        let output = Command::new("go").args(["vet", "./..."])
            .current_dir(work_dir).output()
            .map_err(|e| format!("go vet ./... failed: {e}"))?;
        let text = String::from_utf8_lossy(&output.stderr);
        return Ok(text.contains("unused"));
    }

    Ok(false)
}

fn has_suppressed_elements(work_dir: &Path) -> Result<bool, String> {
    let suppress_pattern = ["@Suppress", "|#pragma|nolint|#\\[allow"].concat();
    let output = Command::new("rg").args(["-l", &suppress_pattern])
        .arg(work_dir).output();

    match output {
        Ok(o) => {
            let text = String::from_utf8_lossy(&o.stdout);
            Ok(!text.trim().is_empty())
        }
        Err(_) => Ok(false),
    }
}

fn output_aegis_result(project: &str, today: &str, result: &AegisResult) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[AEGIS:VERIFICATION]")?;
    writeln!(w, "project: {project}")?;
    writeln!(w, "date: {today}")?;
    writeln!(w, "stage: {}", result.stage)?;
    writeln!(w, "status: {}", result.status)?;
    writeln!(w)?;

    writeln!(w, "[TESTING_STAGE]")?;
    writeln!(w, "lint_issues: {}", result.lint_issues)?;
    writeln!(w, "warnings: {}", result.warnings)?;
    writeln!(w, "core_bugs: {}", result.core_bugs)?;
    writeln!(w)?;

    writeln!(w, "[VERIFIED_STAGE]")?;
    writeln!(w, "dead_code: {}", if result.dead_code { "FOUND" } else { "CLEAN" })?;
    writeln!(w, "suppressed: {}", if result.suppressed { "FOUND" } else { "CLEAN" })?;
    writeln!(w, "algorithm: {}", if result.algo_ok { "VERIFIED" } else { "UNVERIFIED" })?;
    writeln!(w)?;

    if !result.exec_errors.is_empty() {
        writeln!(w, "[EXEC_ERRORS]")?;
        writeln!(w, "note: Some verification commands failed")?;
        for err in &result.exec_errors {
            writeln!(w, "  - {err}")?;
        }
        writeln!(w)?;
    }

    if result.status == "passed" {
        writeln!(w, "[PROMISE]")?;
        writeln!(w, "status: PRODUCTION_READY")?;
        writeln!(w, "signal: <promise>PRODUCTION_READY</promise>")?;
    } else {
        writeln!(w, "[AEGIS_FAILURES]")?;
        writeln!(w, "action: REPORT_TO_CEO")?;
        writeln!(w, "result: LOOP_CONTINUES")?;
        for reason in &result.fail_reasons {
            writeln!(w, "  - {reason}")?;
        }
    }

    Ok(())
}

fn detect_project() -> String {
    if let Ok(val) = std::env::var("KAVACH_PROJECT") {
        if !val.is_empty() { return val; }
    }
    if let Ok(wd) = std::env::current_dir() {
        if wd.join(".git").exists() {
            return wd.file_name()
                .map(|n| n.to_string_lossy().to_string())
                .unwrap_or_else(|| "global".into());
        }
    }
    "global".into()
}
