//! Skills management: discover, list, get, inject research context.
//! Sources: built-in + ~/.claude/skills/ + .claude/skills/
//! Protocol: SP/3.0 (Sutra Protocol)

use std::io::Write;
use std::path::{Path, PathBuf};

use clap::Args;

#[derive(Args)]
pub struct SkillsArgs {
    #[arg(long)]
    pub get: Option<String>,
    #[arg(long)]
    pub sutra: bool,
    #[arg(long)]
    pub inject: bool,
}

struct Skill {
    name: String,
    category: String,
    description: String,
    triggers: Vec<String>,
    commands: Vec<String>,
    path: String,
    research: Vec<String>,
    patterns: Vec<String>,
}

pub fn run(args: SkillsArgs) -> anyhow::Result<()> {
    let mut skills = discover();

    if let Some(name) = &args.get {
        let skill = match find_skill(&mut skills, name) {
            Some(s) => s,
            None => {
                let stderr = std::io::stderr();
                let mut h = stderr.lock();
                let _ = writeln!(h, "[ERROR] Skill not found: {name}");
                return Ok(());
            }
        };

        if args.inject {
            inject_context(skill);
        }

        if args.sutra {
            output_sutra_single(skill)
        } else {
            output_toon_single(skill)
        }
    } else if args.sutra {
        output_sutra_list(&skills)
    } else {
        output_toon_list(&skills)
    }
}

fn output_toon_list(skills: &[Skill]) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();
    let today = chrono::Local::now().format("%Y-%m-%d").to_string();

    writeln!(w, "[SKILLS]")?;
    writeln!(w, "count: {}", skills.len())?;
    writeln!(w, "date: {today}")?;
    writeln!(w)?;

    let categories = ["git", "session", "memory", "research", "build", "test"];
    for cat in &categories {
        let list: Vec<&Skill> = skills.iter().filter(|s| s.category == *cat).collect();
        if list.is_empty() { continue; }
        writeln!(w, "[{}]", cat.to_uppercase())?;
        for s in &list {
            writeln!(w, "{}: {}", s.name, s.description)?;
        }
        writeln!(w)?;
    }

    let general: Vec<&Skill> = skills.iter()
        .filter(|s| !categories.contains(&s.category.as_str()))
        .collect();
    if !general.is_empty() {
        writeln!(w, "[GENERAL]")?;
        for s in &general {
            writeln!(w, "{}: {}", s.name, s.description)?;
        }
    }

    Ok(())
}

fn output_toon_single(s: &Skill) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[SKILL:{}]", s.name.to_uppercase())?;
    writeln!(w, "name: {}", s.name)?;
    writeln!(w, "category: {}", s.category)?;
    writeln!(w, "description: {}", s.description)?;
    if !s.path.is_empty() {
        writeln!(w, "source: {}", s.path)?;
    }
    writeln!(w)?;

    if !s.triggers.is_empty() {
        writeln!(w, "[TRIGGERS]")?;
        for t in &s.triggers { writeln!(w, "- {t}")?; }
        writeln!(w)?;
    }

    if !s.commands.is_empty() {
        writeln!(w, "[COMMANDS]")?;
        for c in &s.commands { writeln!(w, "- {c}")?; }
        writeln!(w)?;
    }

    if !s.research.is_empty() {
        writeln!(w, "[RESEARCH_CONTEXT]")?;
        for r in &s.research { writeln!(w, "- {r}")?; }
        writeln!(w)?;
    }

    if !s.patterns.is_empty() {
        writeln!(w, "[PATTERNS]")?;
        for p in &s.patterns { writeln!(w, "- {p}")?; }
    }

    Ok(())
}

fn output_sutra_list(skills: &[Skill]) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();
    let today = chrono::Local::now().format("%Y-%m-%d").to_string();

    writeln!(w, "[META]")?;
    writeln!(w, "protocol: SP/3.0")?;
    writeln!(w, "from: kavach/skills")?;
    writeln!(w, "to: Claude")?;
    writeln!(w, "date: {today}")?;
    writeln!(w, "count: {}", skills.len())?;
    writeln!(w)?;

    writeln!(w, "[AVAILABLE_SKILLS]")?;
    for s in skills {
        let triggers = if s.triggers.is_empty() {
            format!("/{}", s.name)
        } else {
            s.triggers.join(",")
        };
        writeln!(w, "{}: {triggers}", s.name)?;
    }
    writeln!(w)?;

    writeln!(w, "[CATEGORIES]")?;
    let categories = ["git", "session", "memory", "research", "build", "test"];
    for cat in &categories {
        let names: Vec<&str> = skills.iter()
            .filter(|s| s.category == *cat)
            .map(|s| s.name.as_str())
            .collect();
        if !names.is_empty() {
            writeln!(w, "{cat}: {}", names.join(","))?;
        }
    }

    Ok(())
}

fn output_sutra_single(s: &Skill) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();
    let today = chrono::Local::now().format("%Y-%m-%d").to_string();

    writeln!(w, "[META]")?;
    writeln!(w, "protocol: SP/3.0")?;
    writeln!(w, "from: kavach/skills")?;
    writeln!(w, "skill: {}", s.name)?;
    writeln!(w, "date: {today}")?;
    writeln!(w)?;

    writeln!(w, "[SKILL:{}]", s.name.to_uppercase())?;
    writeln!(w, "name: {}", s.name)?;
    writeln!(w, "category: {}", s.category)?;
    writeln!(w, "desc: {}", s.description)?;
    writeln!(w)?;

    writeln!(w, "[CONTEXT]")?;
    writeln!(w, "triggers: {}", s.triggers.join(","))?;
    writeln!(w, "commands: {}", s.commands.join(","))?;
    writeln!(w)?;

    if !s.research.is_empty() || !s.patterns.is_empty() {
        writeln!(w, "[INJECTED_CONTEXT]")?;
        writeln!(w, "research_entries: {}", s.research.len())?;
        writeln!(w, "pattern_entries: {}", s.patterns.len())?;
    }

    Ok(())
}

// --- Discovery ---

fn discover() -> Vec<Skill> {
    let mut skills = builtin_skills();

    let home = dirs::home_dir().unwrap_or_default();
    let global_dir = home.join(".claude").join("skills");
    let discovered_global = from_dir(&global_dir);
    skills = merge_skills(skills, discovered_global);

    if let Ok(wd) = std::env::current_dir() {
        let project_dir = wd.join(".claude").join("skills");
        let discovered_project = from_dir(&project_dir);
        skills = merge_skills(skills, discovered_project);
    }

    skills
}

fn from_dir(dir: &Path) -> Vec<Skill> {
    let mut skills = Vec::new();
    let entries = match std::fs::read_dir(dir) {
        Ok(e) => e,
        Err(_) => return skills,
    };

    for entry in entries.flatten() {
        if !entry.path().is_dir() { continue; }
        let skill_path = entry.path().join("SKILL.md");
        if !skill_path.exists() { continue; }
        if let Some(skill) = parse_skill_file(&skill_path) {
            skills.push(skill);
        }
    }

    skills
}

fn parse_skill_file(path: &Path) -> Option<Skill> {
    let content = std::fs::read_to_string(path).ok()?;
    let dir_name = path.parent()?.file_name()?.to_string_lossy().to_string();

    let mut skill = Skill {
        name: dir_name,
        category: String::new(),
        description: String::new(),
        triggers: Vec::new(),
        commands: Vec::new(),
        path: path.to_string_lossy().to_string(),
        research: Vec::new(),
        patterns: Vec::new(),
    };

    let mut in_frontmatter = false;
    let mut frontmatter_done = false;

    for line in content.lines() {
        if line == "---" {
            if !in_frontmatter {
                in_frontmatter = true;
            } else {
                frontmatter_done = true;
                in_frontmatter = false;
            }
            continue;
        }

        if in_frontmatter {
            if let Some(rest) = line.strip_prefix("name:") {
                skill.name = rest.trim().to_string();
            } else if let Some(rest) = line.strip_prefix("category:") {
                skill.category = rest.trim().to_string();
            } else if let Some(rest) = line.strip_prefix("description:") {
                skill.description = rest.trim().to_string();
            }
        }

        if frontmatter_done && skill.description.is_empty() {
            let trimmed = line.trim();
            if !trimmed.is_empty() && !trimmed.starts_with('#') {
                skill.description = trimmed.to_string();
            }
        }
    }

    Some(skill)
}

fn merge_skills(builtin: Vec<Skill>, discovered: Vec<Skill>) -> Vec<Skill> {
    let mut result = Vec::with_capacity(builtin.len() + discovered.len());
    let mut seen = std::collections::HashSet::new();

    for s in discovered {
        seen.insert(s.name.clone());
        result.push(s);
    }
    for s in builtin {
        if !seen.contains(&s.name) {
            result.push(s);
        }
    }

    result
}

fn find_skill<'a>(skills: &'a mut [Skill], name: &str) -> Option<&'a mut Skill> {
    let lower = name.to_lowercase();
    skills.iter_mut().find(|s| s.name.to_lowercase() == lower)
}

fn inject_context(skill: &mut Skill) {
    let mem_dir = memory_dir();
    let project = detect_project();

    let research_path = mem_dir.join("research").join(&project).join("research.toon");
    if research_path.exists() {
        skill.research = load_entries(&research_path, &skill.name);
    }

    let patterns_path = mem_dir.join("patterns").join(&project).join("patterns.toon");
    if patterns_path.exists() {
        skill.patterns = load_entries(&patterns_path, &skill.name);
    }
}

fn load_entries(path: &Path, skill_name: &str) -> Vec<String> {
    let content = match std::fs::read_to_string(path) {
        Ok(c) => c,
        Err(_) => return Vec::new(),
    };

    let lower_skill = skill_name.to_lowercase();
    content.lines()
        .map(|l| l.trim())
        .filter(|l| !l.is_empty() && !l.starts_with('#'))
        .filter(|l| l.to_lowercase().contains(&lower_skill))
        .map(|l| l.to_string())
        .collect()
}

// --- Built-in Skills ---

fn builtin_skills() -> Vec<Skill> {
    vec![
        Skill { name: "commit".into(), category: "git".into(), description: "Create git commit with conventional format".into(), triggers: vec!["/commit".into(), "commit changes".into()], commands: vec!["git add".into(), "git commit".into()], path: String::new(), research: Vec::new(), patterns: Vec::new() },
        Skill { name: "review-pr".into(), category: "git".into(), description: "Review pull request changes".into(), triggers: vec!["/review-pr".into(), "review pr".into()], commands: vec!["gh pr view".into(), "gh pr diff".into()], path: String::new(), research: Vec::new(), patterns: Vec::new() },
        Skill { name: "create-pr".into(), category: "git".into(), description: "Create pull request".into(), triggers: vec!["/create-pr".into(), "create pr".into()], commands: vec!["gh pr create".into()], path: String::new(), research: Vec::new(), patterns: Vec::new() },
        Skill { name: "init".into(), category: "session".into(), description: "Initialize session with context".into(), triggers: vec!["/init".into(), "start session".into()], commands: vec!["kavach session init".into()], path: String::new(), research: Vec::new(), patterns: Vec::new() },
        Skill { name: "status".into(), category: "session".into(), description: "Show system status".into(), triggers: vec!["/status".into(), "show status".into()], commands: vec!["kavach status".into()], path: String::new(), research: Vec::new(), patterns: Vec::new() },
        Skill { name: "memory".into(), category: "memory".into(), description: "Query memory bank".into(), triggers: vec!["/memory".into(), "memory bank".into()], commands: vec!["kavach memory bank".into()], path: String::new(), research: Vec::new(), patterns: Vec::new() },
        Skill { name: "kanban".into(), category: "memory".into(), description: "Show kanban dashboard".into(), triggers: vec!["/kanban".into(), "show kanban".into()], commands: vec!["kavach memory kanban".into()], path: String::new(), research: Vec::new(), patterns: Vec::new() },
        Skill { name: "plan".into(), category: "research".into(), description: "Plan implementation approach".into(), triggers: vec!["/plan".into(), "create plan".into()], commands: vec!["EnterPlanMode".into()], path: String::new(), research: Vec::new(), patterns: Vec::new() },
        Skill { name: "explore".into(), category: "research".into(), description: "Explore codebase".into(), triggers: vec!["/explore".into(), "explore code".into()], commands: vec!["Task(Explore)".into()], path: String::new(), research: Vec::new(), patterns: Vec::new() },
        Skill { name: "build".into(), category: "build".into(), description: "Build project".into(), triggers: vec!["/build".into(), "build project".into()], commands: vec!["make".into(), "go build".into(), "npm run build".into()], path: String::new(), research: Vec::new(), patterns: Vec::new() },
        Skill { name: "test".into(), category: "test".into(), description: "Run tests".into(), triggers: vec!["/test".into(), "run tests".into()], commands: vec!["go test".into(), "npm test".into(), "pytest".into()], path: String::new(), research: Vec::new(), patterns: Vec::new() },
    ]
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

fn memory_dir() -> PathBuf {
    let home = dirs::home_dir().unwrap_or_default();
    home.join(".local").join("shared").join("shared-ai").join("memory")
}
