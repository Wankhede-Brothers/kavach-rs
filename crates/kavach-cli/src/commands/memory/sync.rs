//! Memory sync: sync tool events to Memory Bank (STM session log + scratchpad).
//! Hook: PostToolUse — tracks task create/update, file changes, bash, agents.
//! Manual: kavach memory sync --task "name" --status completed

use std::io::Write;
use std::path::PathBuf;

use clap::Args;

use crate::hook;
use crate::session;

#[derive(Args)]
pub struct SyncArgs {
    #[arg(long)]
    pub hook: bool,
    #[arg(long)]
    pub task: Option<String>,
    #[arg(long)]
    pub status: Option<String>,
}

pub fn run(args: SyncArgs) -> anyhow::Result<()> {
    let project = detect_project();
    let today = chrono::Local::now().format("%Y-%m-%d").to_string();

    if args.hook {
        return run_sync_hook(&project, &today);
    }

    if let (Some(task), Some(status)) = (&args.task, &args.status) {
        update_scratchpad(&project, &today, task, status);
        let out = std::io::stdout();
        let mut w = out.lock();
        writeln!(w, "[SYNC] {task} -> {status}")?;
        return Ok(());
    }

    crate::commands::cli_print("memory sync: use --hook or --task + --status");
    Ok(())
}

fn run_sync_hook(project: &str, today: &str) -> anyhow::Result<()> {
    let input = hook::must_read_hook_input();
    let tool = input.get_tool_name();

    match tool {
        "TaskCreate" => {
            let subject = input.get_string("subject");
            if !subject.is_empty() {
                append_stm_event(project, "task_created", &subject);
            }
        }
        "TaskUpdate" => {
            let status = input.get_string("status");
            let subject = input.get_string("subject");
            let task_id = input.get_string("taskId");
            if !task_id.is_empty() && !status.is_empty() {
                append_stm_event(project, &format!("task_{status}"), &subject);
                if status == "completed" && !subject.is_empty() {
                    update_scratchpad(project, today, &subject, "completed");
                }
                if status == "in_progress" && !subject.is_empty() {
                    update_scratchpad(project, today, &subject, "in_progress");
                }
            }
        }
        "Write" | "Edit" => {
            let file_path = input.get_string("file_path");
            if !file_path.is_empty() {
                let mut sess = session::get_or_create_session();
                sess.add_file_modified(&file_path);
                let _ = sess.save();
                append_stm_event(project, &format!("file_{tool}"), &file_path);
            }
        }
        "Bash" => {
            let command = input.get_string("command");
            if !command.is_empty() && is_significant_bash(&command) {
                append_stm_event(project, "bash", &command);
            }
        }
        "Task" => {
            let desc = input.get_string("description");
            let agent = input.get_string("subagent_type");
            if !desc.is_empty() {
                append_stm_event(project, &format!("agent_{agent}"), &desc);
            }
        }
        _ => {}
    }

    hook::exit_silent()?;
    Ok(())
}

fn is_significant_bash(cmd: &str) -> bool {
    let significant = [
        "build", "test", "deploy", "cargo", "go ", "bun ",
        "npm ", "git commit", "git push", "git merge",
    ];
    let lower = cmd.to_lowercase();
    significant.iter().any(|s| lower.contains(s))
}

fn append_stm_event(project: &str, event_type: &str, detail: &str) {
    let stm_dir = stm_dir();
    let _ = std::fs::create_dir_all(&stm_dir);

    let log_path = stm_dir.join("session-log.toon");
    let timestamp = chrono::Local::now().format("%H:%M:%S").to_string();

    let entry = format!("{timestamp} [{event_type}] {project}: {detail}\n");

    let mut file = std::fs::OpenOptions::new()
        .create(true)
        .append(true)
        .open(&log_path)
        .unwrap_or_else(|_| std::fs::File::create(&log_path).expect("stm log"));
    let _ = file.write_all(entry.as_bytes());
}

fn update_scratchpad(project: &str, today: &str, task: &str, status: &str) {
    if project.is_empty() {
        return;
    }

    let stm_dir = stm_dir();
    let project_dir = stm_dir.join("projects").join(project);
    let _ = std::fs::create_dir_all(&project_dir);

    let path = project_dir.join("scratchpad.toon");
    let content = format!(
        "# Project Scratchpad - SP/1.0\n# Auto-updated by kavach memory sync\n\n\
         [SCRATCHPAD:{project}]\nworkdir: {wd}\nupdated: {today}\n\n\
         [TASK]\nintent: {task}\nstatus: {status}\n",
        wd = std::env::current_dir()
            .map(|p| p.to_string_lossy().to_string())
            .unwrap_or_default(),
    );

    let _ = std::fs::write(&path, content);
}

fn stm_dir() -> PathBuf {
    let home = dirs::home_dir().unwrap_or_default();
    home.join(".local")
        .join("shared")
        .join("shared-ai")
        .join("stm")
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
