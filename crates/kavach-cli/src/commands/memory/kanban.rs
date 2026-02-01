//! Kanban dashboard: Sprint/Kanban visual dashboard with Aegis-Guard verification.
//! Path: ~/.local/shared/shared-ai/memory/kanban/<project>/kanban.toon
//! Modes: default (visual), --status (TOON), --sutra (SP/1.0), --all via --project

use std::io::Write;
use std::path::{Path, PathBuf};

use clap::Args;

#[derive(Args)]
pub struct KanbanArgs {
    #[arg(long)]
    pub status: bool,
    #[arg(long)]
    pub visual: bool,
    #[arg(long)]
    pub sutra: bool,
    #[arg(short = 'p', long)]
    pub project: Option<String>,
}

struct KanbanCard {
    id: String,
    column: String,
    title: String,
    priority: String,
    aegis_status: String,
    lint_issues: i32,
    warnings: i32,
    core_bugs: i32,
}

struct KanbanBoard {
    project: String,
    updated: String,
    phases: Vec<Vec<KanbanCard>>,
    loop_count: i32,
}

const COL_BACKLOG: &str = "backlog";
const COL_IN_PROGRESS: &str = "in_progress";
const COL_TESTING: &str = "testing";
const COL_VERIFIED: &str = "verified";
const COL_DONE: &str = "done";

pub fn run(args: KanbanArgs) -> anyhow::Result<()> {
    let kanban_dir = memory_dir().join("kanban");
    let project = args
        .project
        .unwrap_or_else(|| detect_kanban_project(&kanban_dir));
    let board = load_kanban_toon(&kanban_dir, &project);

    if args.sutra {
        output_sutra(&board)
    } else if args.status {
        output_toon_status(&board)
    } else {
        output_visual(&board)
    }
}

fn output_sutra(board: &KanbanBoard) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();
    let today = chrono::Local::now().format("%Y-%m-%d").to_string();
    let counts = count_by_column(board);

    writeln!(w, "[META]")?;
    writeln!(w, "protocol: SP/1.0")?;
    writeln!(w, "from: kavach/kanban")?;
    writeln!(w, "to: CEO")?;
    writeln!(w, "date: {today}")?;
    writeln!(w, "project: {}", board.project)?;
    writeln!(w)?;

    writeln!(w, "[KANBAN_STATE]")?;
    writeln!(w, "backlog: {}", counts.0)?;
    writeln!(w, "in_progress: {}", counts.1)?;
    writeln!(w, "testing: {}", counts.2)?;
    writeln!(w, "verified: {}", counts.3)?;
    writeln!(w, "done: {}", counts.4)?;
    writeln!(w, "loop_count: {}", board.loop_count)?;
    writeln!(w)?;

    writeln!(w, "[AEGIS_QUEUE]")?;
    let testing = tasks_by_column(board, COL_TESTING);
    if !testing.is_empty() {
        writeln!(w, "stage: TESTING")?;
        writeln!(w, "checks: lint,warnings,core_bugs,unit_tests")?;
        for t in &testing {
            writeln!(w, "task: {},{},{}", t.id, t.title, t.aegis_status)?;
        }
    }
    let verified = tasks_by_column(board, COL_VERIFIED);
    if !verified.is_empty() {
        writeln!(w, "stage: VERIFIED")?;
        writeln!(w, "checks: algorithm,dead_code,suppressed,hidden_bugs")?;
        for t in &verified {
            writeln!(w, "task: {},{},{}", t.id, t.title, t.aegis_status)?;
        }
    }
    writeln!(w)?;

    let failed = failed_tasks(board);
    if !failed.is_empty() {
        writeln!(w, "[AEGIS_FAILURES]")?;
        writeln!(w, "action: REPORT_TO_CEO")?;
        writeln!(w, "result: LOOP_CONTINUES")?;
        for t in &failed {
            writeln!(
                w,
                "failed: {},{},lint:{},warn:{},bugs:{}",
                t.id, t.title, t.lint_issues, t.warnings, t.core_bugs
            )?;
        }
        writeln!(w)?;
    }

    let total = counts.0 + counts.1 + counts.2 + counts.3 + counts.4;
    let progress = if total > 0 {
        (counts.4 * 100) / total
    } else {
        0
    };

    writeln!(w, "[PROMISE]")?;
    if progress == 100 && failed.is_empty() {
        writeln!(w, "status: PRODUCTION_READY")?;
        writeln!(w, "signal: <promise>PRODUCTION_READY</promise>")?;
    } else {
        writeln!(w, "status: IN_PROGRESS")?;
        writeln!(w, "progress: {progress}%")?;
        writeln!(w, "signal: LOOP_CONTINUES")?;
    }

    Ok(())
}

fn output_toon_status(board: &KanbanBoard) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();
    let counts = count_by_column(board);
    let priorities = count_by_priority(board);
    let total = counts.0 + counts.1 + counts.2 + counts.3 + counts.4;
    let progress = if total > 0 {
        (counts.4 * 100) / total
    } else {
        0
    };

    writeln!(w, "[KANBAN]")?;
    writeln!(w, "project: {}", board.project)?;
    writeln!(w, "updated: {}", board.updated)?;
    writeln!(w, "loop_count: {}", board.loop_count)?;
    writeln!(w)?;

    writeln!(w, "[PIPELINE]")?;
    writeln!(w, "backlog: {}", counts.0)?;
    writeln!(w, "in_progress: {}", counts.1)?;
    writeln!(w, "testing: {}", counts.2)?;
    writeln!(w, "verified: {}", counts.3)?;
    writeln!(w, "done: {}", counts.4)?;
    writeln!(w)?;

    writeln!(w, "[PRIORITY]")?;
    writeln!(w, "critical: {}", priorities.0)?;
    writeln!(w, "high: {}", priorities.1)?;
    writeln!(w, "medium: {}", priorities.2)?;
    writeln!(w, "low: {}", priorities.3)?;
    writeln!(w)?;

    writeln!(w, "[AEGIS]")?;
    let failed_count = failed_tasks(board).len();
    writeln!(w, "failed: {failed_count}")?;
    writeln!(
        w,
        "action: {}",
        if failed_count > 0 {
            "REPORT_TO_CEO"
        } else {
            "CONTINUE"
        }
    )?;
    writeln!(w)?;

    writeln!(w, "[PROGRESS]")?;
    let bar = progress_bar(progress as usize, 20);
    writeln!(w, "{bar} {progress}% ({}/{})", counts.4, total)?;

    Ok(())
}

fn output_visual(board: &KanbanBoard) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();
    let counts = count_by_column(board);
    let total = counts.0 + counts.1 + counts.2 + counts.3 + counts.4;
    let progress = if total > 0 {
        (counts.4 * 100) / total
    } else {
        0
    };
    let sep = "=".repeat(79);
    let dash = "-".repeat(79);

    writeln!(w, "+{sep}+")?;
    writeln!(
        w,
        "|                    KANBAN DASHBOARD: {:<38} |",
        board.project
    )?;
    writeln!(w, "+{sep}+")?;
    writeln!(
        w,
        "|  Updated: {:<20}  Loop Count: {:<5}  Progress: {:>3}%        |",
        board.updated, board.loop_count, progress
    )?;
    writeln!(w, "+{sep}+")?;
    writeln!(w)?;

    writeln!(w, "+{dash}+")?;
    writeln!(
        w,
        "|                           PRODUCTION PIPELINE                                 |"
    )?;
    writeln!(
        w,
        "+----------+-----------+-----------+-----------+-----------+-------------------+"
    )?;
    writeln!(
        w,
        "| BACKLOG  |IN_PROGRESS|  TESTING  | VERIFIED  |   DONE    |      STATUS       |"
    )?;
    writeln!(
        w,
        "+----------+-----------+-----------+-----------+-----------+-------------------+"
    )?;

    let icon = status_icon(progress as usize);
    writeln!(
        w,
        "|   {:>3}    |    {:>3}    |    {:>3}    |    {:>3}    |    {:>3}    | {} |",
        counts.0, counts.1, counts.2, counts.3, counts.4, icon
    )?;
    writeln!(
        w,
        "+----------+-----------+-----------+-----------+-----------+-------------------+"
    )?;
    writeln!(w)?;

    writeln!(w, "+{dash}+")?;
    let bar = progress_bar(progress as usize, 40);
    writeln!(
        w,
        "|  Progress: {bar} {:>3}%                            |",
        progress
    )?;
    writeln!(w, "+{dash}+")?;
    writeln!(w)?;

    let failed = failed_tasks(board);
    let testing = tasks_by_column(board, COL_TESTING);
    let verified = tasks_by_column(board, COL_VERIFIED);

    writeln!(w, "+{dash}+")?;
    writeln!(
        w,
        "|                           AEGIS-GUARD STATUS                                  |"
    )?;
    writeln!(w, "+{dash}+")?;

    if !testing.is_empty() {
        writeln!(
            w,
            "|  TESTING STAGE (Lint, Warnings, Core Bugs):                                  |"
        )?;
        for t in &testing {
            let ic = aegis_icon(&t.aegis_status);
            writeln!(
                w,
                "|    {ic} {:<50} [{:<7}]        |",
                truncate(&t.title, 50),
                t.aegis_status
            )?;
        }
    }
    if !verified.is_empty() {
        writeln!(
            w,
            "|  VERIFIED STAGE (Algorithm, Dead Code, Suppressed):                          |"
        )?;
        for t in &verified {
            let ic = aegis_icon(&t.aegis_status);
            writeln!(
                w,
                "|    {ic} {:<50} [{:<7}]        |",
                truncate(&t.title, 50),
                t.aegis_status
            )?;
        }
    }
    if !failed.is_empty() {
        writeln!(w, "+{dash}+")?;
        writeln!(
            w,
            "|  FAILURES REPORTED TO CEO - LOOP CONTINUES                                   |"
        )?;
        for t in &failed {
            writeln!(w, "|    X {:<60}          |", truncate(&t.title, 60))?;
        }
    }
    writeln!(w, "+{dash}+")?;
    writeln!(w)?;

    writeln!(w, "+{dash}+")?;
    if progress == 100 && failed.is_empty() {
        writeln!(
            w,
            "|                    [OK] PROMISE: PRODUCTION_READY                            |"
        )?;
    } else {
        writeln!(
            w,
            "|                    [..] PROMISE: IN_PROGRESS (LOOP CONTINUES)                |"
        )?;
    }
    writeln!(w, "+{dash}+")?;

    Ok(())
}

fn count_by_column(board: &KanbanBoard) -> (i32, i32, i32, i32, i32) {
    let (mut b, mut ip, mut t, mut v, mut d) = (0, 0, 0, 0, 0);
    for phase in &board.phases {
        for c in phase {
            match c.column.as_str() {
                COL_BACKLOG => b += 1,
                COL_IN_PROGRESS => ip += 1,
                COL_TESTING => t += 1,
                COL_VERIFIED => v += 1,
                COL_DONE => d += 1,
                _ => {}
            }
        }
    }
    (b, ip, t, v, d)
}

fn count_by_priority(board: &KanbanBoard) -> (i32, i32, i32, i32) {
    let (mut c, mut h, mut m, mut l) = (0, 0, 0, 0);
    for phase in &board.phases {
        for card in phase {
            match card.priority.as_str() {
                "critical" => c += 1,
                "high" => h += 1,
                "medium" => m += 1,
                "low" => l += 1,
                _ => {}
            }
        }
    }
    (c, h, m, l)
}

fn tasks_by_column<'a>(board: &'a KanbanBoard, col: &str) -> Vec<&'a KanbanCard> {
    board
        .phases
        .iter()
        .flat_map(|p| p.iter())
        .filter(|c| c.column == col)
        .collect()
}

fn failed_tasks(board: &KanbanBoard) -> Vec<&KanbanCard> {
    board
        .phases
        .iter()
        .flat_map(|p| p.iter())
        .filter(|c| c.aegis_status == "failed")
        .collect()
}

fn progress_bar(percent: usize, width: usize) -> String {
    let filled = (percent * width) / 100;
    let empty = width.saturating_sub(filled);
    format!("[{}{}]", "#".repeat(filled), "-".repeat(empty))
}

fn status_icon(progress: usize) -> &'static str {
    if progress == 100 {
        "[OK] PRODUCTION "
    } else if progress >= 75 {
        "[..] ALMOST READY"
    } else if progress >= 50 {
        "[..] IN PROGRESS "
    } else {
        "[  ] EARLY STAGE "
    }
}

fn aegis_icon(status: &str) -> &'static str {
    match status {
        "passed" => "[OK]",
        "failed" => "[X]",
        "blocked" => "[!]",
        _ => "[ ]",
    }
}

fn truncate(s: &str, max: usize) -> String {
    if s.len() <= max {
        s.to_string()
    } else {
        format!("{}...", &s[..max.saturating_sub(3)])
    }
}

fn load_kanban_toon(kanban_dir: &Path, project: &str) -> KanbanBoard {
    let mut board = KanbanBoard {
        project: project.to_string(),
        updated: chrono::Local::now().format("%Y-%m-%d").to_string(),
        phases: Vec::new(),
        loop_count: 0,
    };

    let toon_path = kanban_dir.join(project).join("kanban.toon");
    let content = match std::fs::read_to_string(&toon_path) {
        Ok(c) => c,
        Err(_) => return board,
    };

    let mut current_phase: Option<usize> = None;

    for line in content.lines() {
        let line = line.trim();
        if line.is_empty() || line.starts_with('#') {
            continue;
        }

        if let Some(rest) = line.strip_prefix("KANBAN:") {
            board.project = rest.to_string();
        } else if let Some(rest) = line.strip_prefix("updated:") {
            board.updated = rest.trim().to_string();
        } else if let Some(rest) = line.strip_prefix("loop_count:") {
            board.loop_count = rest.trim().parse().unwrap_or(0);
        } else if line.starts_with("PHASE_") && line.contains("_CARDS") {
            if let Some(num_str) = line.split('_').nth(1) {
                if let Ok(n) = num_str.parse::<usize>() {
                    current_phase = Some(n);
                    while board.phases.len() <= n {
                        board.phases.push(Vec::new());
                    }
                }
            }
        } else if current_phase.is_some() && (line.starts_with('p') || line.starts_with('P')) {
            if let Some(card) = parse_card_line(line) {
                let idx = current_phase.unwrap();
                while board.phases.len() <= idx {
                    board.phases.push(Vec::new());
                }
                board.phases[idx].push(card);
            }
        }
    }

    board
}

fn parse_card_line(line: &str) -> Option<KanbanCard> {
    let parts: Vec<&str> = line.splitn(10, ',').collect();
    if parts.len() < 4 {
        return None;
    }

    let mut card = KanbanCard {
        id: parts[0].trim().to_string(),
        column: parts[1].trim().to_string(),
        title: parts[2].trim().to_string(),
        priority: parts[3].trim().to_string(),
        aegis_status: "pending".to_string(),
        lint_issues: 0,
        warnings: 0,
        core_bugs: 0,
    };

    if parts.len() >= 6 {
        card.aegis_status = parts[5].trim().to_string();
    }
    if parts.len() >= 7 {
        card.lint_issues = parts[6].trim().parse().unwrap_or(0);
    }
    if parts.len() >= 8 {
        card.warnings = parts[7].trim().parse().unwrap_or(0);
    }
    if parts.len() >= 9 {
        card.core_bugs = parts[8].trim().parse().unwrap_or(0);
    }

    Some(card)
}

fn detect_kanban_project(kanban_dir: &Path) -> String {
    let project = detect_project();
    if kanban_dir.join(&project).join("kanban.toon").exists() {
        return project;
    }
    if let Ok(entries) = std::fs::read_dir(kanban_dir) {
        for entry in entries.flatten() {
            if entry.path().is_dir() {
                let name = entry.file_name().to_string_lossy().to_string();
                if name != "global" {
                    return name;
                }
            }
        }
    }
    "global".into()
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

fn memory_dir() -> PathBuf {
    let home = dirs::home_dir().unwrap_or_default();
    home.join(".local")
        .join("shared")
        .join("shared-ai")
        .join("memory")
}
