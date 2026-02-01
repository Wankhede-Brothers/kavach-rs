//! Post-write umbrella gate (PostToolUse:Write|Edit|NotebookEdit).
//! Hierarchy: ANTIPROD(P0->P3) -> QUALITY -> LINT -> CONTEXT -> MEMORY

use crate::hook::{self, HookInput};
use crate::session;

pub fn run(hook_mode: bool) -> anyhow::Result<()> {
    if !hook_mode {
        crate::commands::cli_print("gates post-write: use --hook flag");
        return Ok(());
    }

    let input = hook::must_read_hook_input();
    let result = dispatch(&input);

    match result {
        Ok(()) => Ok(()),
        Err(e) if hook::is_hook_exit(&e) => Ok(()),
        Err(e) => Err(e),
    }
}

fn dispatch(input: &HookInput) -> anyhow::Result<()> {
    let mut session = session::get_or_create_session();

    let file_path = input.get_string("file_path");
    let content = if input.get_tool_name() == "Edit" {
        input.get_string("new_string")
    } else {
        input.get_string("content")
    };

    // L2: ANTIPROD - P0->P3 hierarchy
    if !content.is_empty() && !file_path.is_empty() {
        run_antiprod_check(&file_path, &content)?;
    }

    // L2: QUALITY - folder depth, line count
    if !content.is_empty() && !file_path.is_empty() {
        run_quality_check(&file_path, &content)?;
    }

    // L2: LINT - whitespace checks
    if !content.is_empty() && !file_path.is_empty() {
        run_lint_check(&file_path, &content);
    }

    // L2: CONTEXT + MEMORY - track file modification
    if !file_path.is_empty() {
        session.add_file_modified(&file_path);
        let _ = session.save();
    }

    hook::exit_silent()
}

fn run_antiprod_check(file_path: &str, content: &str) -> anyhow::Result<()> {
    if is_allowlisted(file_path) {
        return Ok(());
    }

    let ext = file_ext(file_path);
    let base = file_path.rsplit('/').next().unwrap_or("").to_lowercase();

    // P1: Rust-specific
    if ext == ".rs" {
        // Rust macros that indicate non-production code - build at runtime
        let dbg_macro = ["db", "g!", "("].concat();
        let todo_macro = ["to", "do!", "("].concat();
        let unimpl_macro = ["unimp", "lemented!", "("].concat();
        let panic_macro = ["pan", "ic!", "("].concat();
        if content.contains(&dbg_macro) {
            let msg = format!("PROD_LEAK:{}:Remove -- runs in release builds. Use tracing.", dbg_macro);
            hook::exit_block_toon("ANTIPROD", &msg)?;
        }
        if content.contains(&todo_macro) || content.contains(&unimpl_macro) {
            hook::exit_block_toon("ANTIPROD", "PROD_LEAK:macro:Implement before shipping.")?;
        }
        if base != "main.rs" && content.contains(&panic_macro) {
            hook::exit_block_toon("ANTIPROD", "PROD_LEAK:macro:Return Result/Option instead.")?;
        }
        // unsafe without SAFETY comment
        if content.contains("unsafe {") && !content.contains("// SAFETY:") {
            hook::exit_block_toon("ANTIPROD", "PROD_LEAK:unsafe block:Justify with // SAFETY: comment or remove.")?;
        }
        // #[allow(dead_code)] suppression
        if content.contains("#[allow(dead_code)]") {
            hook::exit_block_toon("ANTIPROD", "TYPE_LOOSE:#[allow(dead_code)]:Remove dead code instead of suppressing.")?;
        }
        // #[allow(unused suppression
        if content.contains("#[allow(unused") {
            hook::exit_block_toon("ANTIPROD", "TYPE_LOOSE:#[allow(unused)]:Remove unused code instead of suppressing.")?;
        }
    }

    // P1: JS/TS console.log
    if is_frontend_file(file_path) {
        if content.contains("console.log(") || content.contains("console.debug(") {
            hook::exit_block_toon("ANTIPROD", "PROD_LEAK:console.log:Remove debug output or use structured logger.")?;
        }
    }

    // P1: TODO/FIXME (universal, build at runtime)
    let todo_upper = ["TO", "DO"].concat();
    let fixme_upper = ["FIX", "ME"].concat();
    let content_upper = content.to_uppercase();
    if content_upper.contains(&todo_upper) || content_upper.contains(&fixme_upper) {
        // Check it's actually a comment pattern, not a variable name
        for line in content.lines() {
            let trimmed = line.trim().to_uppercase();
            if trimmed.contains(&format!("// {todo_upper}")) || trimmed.contains(&format!("# {todo_upper}"))
                || trimmed.contains(&format!("// {fixme_upper}")) || trimmed.contains(&format!("# {fixme_upper}"))
            {
                hook::exit_block_toon("ANTIPROD", &format!("PROD_LEAK:{todo_upper}/{fixme_upper}:Implement or create ticket."))?;
            }
        }
    }

    // P1: localhost in non-config files
    if is_non_config_file(file_path) && content.contains("://localhost") {
        hook::exit_block_toon("ANTIPROD", "PROD_LEAK:localhost:Use config/environment variable for URLs.")?;
    }

    // P2: .unwrap() in Rust handler files
    if ext == ".rs" && is_handler_file(file_path) && content.contains(".unwrap()") {
        hook::exit_block_toon("ANTIPROD", "ERROR_BLIND:.unwrap():Use ? operator instead of .unwrap() in handlers.")?;
    }

    // P2: empty catch in JS/TS
    if is_frontend_file(file_path) && content.contains(".catch(() => {})") {
        hook::exit_block_toon("ANTIPROD", "ERROR_BLIND:.catch(() => {}):Handle errors.")?;
    }

    // P3: as any in TS
    if is_frontend_file(file_path) && content.contains("as any") {
        hook::exit_block_toon("ANTIPROD", "TYPE_LOOSE:as any:Use proper type narrowing.")?;
    }

    // P1: Python print()
    if ext == ".py" && content.contains("print(") {
        hook::exit_block_toon("ANTIPROD", "PROD_LEAK:print():Use logging module instead.")?;
    }

    // P1: Go fmt.Print
    if ext == ".go" && content.contains("fmt.Print") {
        hook::exit_block_toon("ANTIPROD", "PROD_LEAK:fmt.Print:Use structured logger.")?;
    }

    // Docker checks
    if is_dockerfile(file_path) {
        if content.contains("FROM ") && content.contains(":latest") {
            hook::exit_block_toon("ANTIPROD", "PROD_LEAK:FROM :latest:Pin image version.")?;
        }
    }

    // chmod 777 (build at runtime)
    let chmod_pattern = format!("chmod {}", "777");
    if content.contains(&chmod_pattern) {
        hook::exit_block_toon("ANTIPROD", &format!("PROD_LEAK:{chmod_pattern}:Use least-privilege permissions."))?;
    }

    Ok(())
}

fn run_quality_check(file_path: &str, content: &str) -> anyhow::Result<()> {
    let ext = file_ext(file_path);
    let code_exts = [".go", ".rs", ".ts", ".tsx", ".js", ".jsx", ".py", ".astro"];
    if !code_exts.iter().any(|e| ext == *e) {
        return Ok(());
    }

    // Folder depth check (max 7)
    if let Ok(wd) = std::env::current_dir() {
        let wd_str = wd.to_string_lossy();
        if file_path.starts_with(wd_str.as_ref()) {
            let rel = &file_path[wd_str.len()..];
            let depth = rel.chars().filter(|c| *c == '/').count();
            if depth > 7 {
                hook::exit_block_toon("DACE", &format!("folder_depth_exceeds_7:{depth}"))?;
            }
        }
    }

    // Line count check (max 100 for new content being written)
    let line_count = content.lines().count();
    if line_count > 100 {
        hook::exit_block_toon("DACE", &format!("exceeds_100_lines:{line_count}"))?;
    }

    Ok(())
}

fn run_lint_check(file_path: &str, content: &str) {
    let mut issues = Vec::new();
    for (i, line) in content.lines().enumerate() {
        if line.ends_with(' ') || line.ends_with('\t') {
            issues.push(format!("trailing_ws:{}", i + 1));
        }
    }

    // Go: spaces instead of tabs
    if file_ext(file_path) == ".go" {
        for (i, line) in content.lines().enumerate() {
            if line.starts_with("    ") && !line.starts_with("\t") {
                issues.push(format!("spaces:{}", i + 1));
            }
        }
    }

    if !issues.is_empty() {
        let max = issues.len().min(3);
        let stderr = std::io::stderr();
        let mut handle = stderr.lock();
        let _ = std::io::Write::write_all(
            &mut handle,
            format!("[LINT] {}\n", issues[..max].join(",")).as_bytes(),
        );
    }
}

fn file_ext(path: &str) -> &str {
    path.rfind('.').map(|i| &path[i..]).unwrap_or("")
}

fn is_frontend_file(path: &str) -> bool {
    let exts = [".ts", ".tsx", ".js", ".jsx", ".astro", ".vue", ".svelte"];
    let p = path.to_lowercase();
    exts.iter().any(|e| p.ends_with(e))
}

fn is_handler_file(path: &str) -> bool {
    let p = path.to_lowercase();
    p.contains("handler") || p.contains("routes") || p.contains("lib.rs")
}

fn is_non_config_file(path: &str) -> bool {
    let p = path.to_lowercase();
    let config_patterns = [
        "config", ".env", "astro.config", "vite.config", "next.config",
        "wrangler.toml", "docker-compose", ".toml", "constants",
    ];
    let non_code_exts = [".md", ".txt", ".json", ".yaml", ".yml", ".csv", ".xml", ".html"];
    if config_patterns.iter().any(|pat| p.contains(pat)) {
        return false;
    }
    if non_code_exts.iter().any(|ext| p.ends_with(ext)) {
        return false;
    }
    true
}

fn is_dockerfile(path: &str) -> bool {
    let base = path.rsplit('/').next().unwrap_or("").to_lowercase();
    base.starts_with("dockerfile") || base == "containerfile"
}

fn is_allowlisted(path: &str) -> bool {
    let p = path.to_lowercase();
    p.contains("test") || p.contains("spec") || p.contains("mock")
        || p.contains("fixture") || p.contains("__test")
}
