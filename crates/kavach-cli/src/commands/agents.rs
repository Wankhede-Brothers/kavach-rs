//! Agent management: discover, list, get, inject research context.
//! Sources: built-in hierarchy + ~/.claude/agents/ + .claude/agents/
//! Protocol: SP/3.0 (Sutra Protocol)

use std::io::Write;
use std::path::{Path, PathBuf};

use clap::Args;

#[derive(Args)]
pub struct AgentsArgs {
    #[arg(long)]
    pub get: Option<String>,
    #[arg(long)]
    pub sutra: bool,
    #[arg(long)]
    pub inject: bool,
}

struct Agent {
    name: String,
    level: i32,
    model: String,
    description: String,
    triggers: Vec<String>,
    tools: Vec<String>,
    path: String,
    research: Vec<String>,
    patterns: Vec<String>,
}

const LEVEL_NLU: i32 = -1;
const LEVEL_CEO: i32 = 0;
const LEVEL_ENGINEER: i32 = 1;
const LEVEL_REVIEW: i32 = 2;
const LEVEL_AEGIS: i32 = 3;

fn level_name(level: i32) -> &'static str {
    match level {
        LEVEL_NLU => "L-1 (NLU)",
        LEVEL_CEO => "L0 (CEO)",
        LEVEL_ENGINEER => "L1 (Engineers)",
        LEVEL_REVIEW => "L1.5 (Review)",
        LEVEL_AEGIS => "L2 (Aegis)",
        _ => "Unknown",
    }
}

pub fn run(args: AgentsArgs) -> anyhow::Result<()> {
    let mut agents = discover();

    if let Some(name) = &args.get {
        let agent = match find_agent(&mut agents, name) {
            Some(a) => a,
            None => {
                let stderr = std::io::stderr();
                let mut h = stderr.lock();
                let _ = writeln!(h, "[ERROR] Agent not found: {name}");
                return Ok(());
            }
        };

        if args.inject {
            inject_context(agent);
        }

        if args.sutra {
            output_sutra_single(agent)
        } else {
            output_toon_single(agent)
        }
    } else if args.sutra {
        output_sutra_list(&agents)
    } else {
        output_toon_list(&agents)
    }
}

fn output_toon_list(agents: &[Agent]) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();
    let today = chrono::Local::now().format("%Y-%m-%d").to_string();

    writeln!(w, "[AGENTS]")?;
    writeln!(w, "count: {}", agents.len())?;
    writeln!(w, "date: {today}")?;
    writeln!(w)?;

    let levels = [
        LEVEL_NLU,
        LEVEL_CEO,
        LEVEL_ENGINEER,
        LEVEL_REVIEW,
        LEVEL_AEGIS,
    ];
    for lvl in &levels {
        let list: Vec<&Agent> = agents.iter().filter(|a| a.level == *lvl).collect();
        if list.is_empty() {
            continue;
        }
        writeln!(w, "[{}]", level_name(*lvl))?;
        for a in &list {
            writeln!(w, "{}: {} ({})", a.name, a.description, a.model)?;
        }
        writeln!(w)?;
    }

    writeln!(w, "[MODELS]")?;
    writeln!(w, "opus: ceo,research-director,aegis-guardian")?;
    writeln!(
        w,
        "sonnet: backend,frontend,devops,security,qa,code-reviewer"
    )?;
    writeln!(w, "haiku: nlu-intent-analyzer")?;

    Ok(())
}

fn output_toon_single(a: &Agent) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();

    writeln!(w, "[AGENT:{}]", a.name.to_uppercase())?;
    writeln!(w, "name: {}", a.name)?;
    writeln!(w, "level: {}", a.level)?;
    writeln!(w, "model: {}", a.model)?;
    writeln!(w, "description: {}", a.description)?;
    if !a.path.is_empty() {
        writeln!(w, "source: {}", a.path)?;
    }
    writeln!(w)?;

    if !a.triggers.is_empty() {
        writeln!(w, "[TRIGGERS]")?;
        for t in &a.triggers {
            writeln!(w, "- {t}")?;
        }
        writeln!(w)?;
    }

    if !a.tools.is_empty() {
        writeln!(w, "[TOOLS]")?;
        writeln!(w, "{}", a.tools.join(", "))?;
        writeln!(w)?;
    }

    if !a.research.is_empty() {
        writeln!(w, "[RESEARCH_CONTEXT]")?;
        for r in &a.research {
            writeln!(w, "- {r}")?;
        }
        writeln!(w)?;
    }

    if !a.patterns.is_empty() {
        writeln!(w, "[PATTERNS]")?;
        for p in &a.patterns {
            writeln!(w, "- {p}")?;
        }
    }

    Ok(())
}

fn output_sutra_list(agents: &[Agent]) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();
    let today = chrono::Local::now().format("%Y-%m-%d").to_string();

    writeln!(w, "[META]")?;
    writeln!(w, "protocol: SP/3.0")?;
    writeln!(w, "from: kavach/agents")?;
    writeln!(w, "to: CEO")?;
    writeln!(w, "date: {today}")?;
    writeln!(w, "count: {}", agents.len())?;
    writeln!(w)?;

    writeln!(w, "[AGENT_HIERARCHY]")?;
    writeln!(w, "L-1: nlu-intent-analyzer (haiku)")?;
    writeln!(w, "L0: ceo, research-director (opus)")?;
    writeln!(w, "L1: backend, frontend, devops, security, qa (sonnet)")?;
    writeln!(w, "L1.5: code-reviewer (sonnet)")?;
    writeln!(w, "L2: aegis-guardian (opus)")?;
    writeln!(w)?;

    writeln!(w, "[DELEGATION_RULES]")?;
    writeln!(w, "CEO_ONLY: orchestration,delegation,decisions")?;
    writeln!(w, "ENGINEER: implementation,code_changes")?;
    writeln!(w, "AEGIS: verification,quality_gate")?;
    writeln!(w)?;

    writeln!(w, "[AVAILABLE]")?;
    for a in agents {
        writeln!(w, "{}: {}", a.name, a.model)?;
    }

    Ok(())
}

fn output_sutra_single(a: &Agent) -> anyhow::Result<()> {
    let out = std::io::stdout();
    let mut w = out.lock();
    let today = chrono::Local::now().format("%Y-%m-%d").to_string();

    writeln!(w, "[META]")?;
    writeln!(w, "protocol: SP/3.0")?;
    writeln!(w, "from: kavach/agents")?;
    writeln!(w, "to: {}", a.name)?;
    writeln!(w, "date: {today}")?;
    writeln!(w)?;

    writeln!(w, "[AGENT:{}]", a.name.to_uppercase())?;
    writeln!(w, "name: {}", a.name)?;
    writeln!(w, "level: {}", a.level)?;
    writeln!(w, "model: {}", a.model)?;
    writeln!(w, "desc: {}", a.description)?;
    writeln!(w)?;

    writeln!(w, "[CONTEXT]")?;
    writeln!(w, "tools: {}", a.tools.join(","))?;
    writeln!(w, "triggers: {}", a.triggers.join(","))?;
    writeln!(w)?;

    if !a.research.is_empty() || !a.patterns.is_empty() {
        writeln!(w, "[INJECTED_CONTEXT]")?;
        writeln!(w, "research_entries: {}", a.research.len())?;
        writeln!(w, "pattern_entries: {}", a.patterns.len())?;
    }

    Ok(())
}

// --- Discovery ---

fn discover() -> Vec<Agent> {
    let mut agents = builtin_agents();

    let home = dirs::home_dir().unwrap_or_default();
    let global_dir = home.join(".claude").join("agents");
    let discovered_global = from_dir(&global_dir);
    agents = merge_agents(agents, discovered_global);

    if let Ok(wd) = std::env::current_dir() {
        let project_dir = wd.join(".claude").join("agents");
        let discovered_project = from_dir(&project_dir);
        agents = merge_agents(agents, discovered_project);
    }

    agents.sort_by_key(|a| a.level);
    agents
}

fn from_dir(dir: &Path) -> Vec<Agent> {
    let mut agents = Vec::new();
    let entries = match std::fs::read_dir(dir) {
        Ok(e) => e,
        Err(_) => return agents,
    };

    for entry in entries.flatten() {
        let path = entry.path();
        if path.is_dir() || !path.extension().map(|e| e == "md").unwrap_or(false) {
            continue;
        }
        if let Some(agent) = parse_agent_file(&path) {
            agents.push(agent);
        }
    }

    agents
}

fn parse_agent_file(path: &Path) -> Option<Agent> {
    let content = std::fs::read_to_string(path).ok()?;
    let file_stem = path.file_stem()?.to_string_lossy().to_string();

    let mut agent = Agent {
        name: file_stem,
        level: LEVEL_ENGINEER,
        model: "sonnet".into(),
        description: String::new(),
        triggers: Vec::new(),
        tools: Vec::new(),
        path: path.to_string_lossy().to_string(),
        research: Vec::new(),
        patterns: Vec::new(),
    };

    let mut in_frontmatter = false;

    for line in content.lines() {
        if line == "---" {
            in_frontmatter = !in_frontmatter;
            continue;
        }

        if in_frontmatter {
            if let Some(rest) = line.strip_prefix("name:") {
                agent.name = rest.trim().to_string();
            } else if let Some(rest) = line.strip_prefix("level:") {
                agent.level = parse_level(rest.trim());
            } else if let Some(rest) = line.strip_prefix("model:") {
                agent.model = rest.trim().to_string();
            } else if let Some(rest) = line.strip_prefix("description:") {
                agent.description = rest.trim().to_string();
            }
        }
    }

    Some(agent)
}

fn parse_level(s: &str) -> i32 {
    match s {
        "-1" | "nlu" => LEVEL_NLU,
        "0" | "ceo" | "l0" => LEVEL_CEO,
        "1" | "engineer" | "l1" => LEVEL_ENGINEER,
        "1.5" | "2" | "review" => LEVEL_REVIEW,
        "3" | "aegis" | "l2" => LEVEL_AEGIS,
        _ => LEVEL_ENGINEER,
    }
}

fn merge_agents(builtin: Vec<Agent>, discovered: Vec<Agent>) -> Vec<Agent> {
    let mut result = Vec::with_capacity(builtin.len() + discovered.len());
    let mut seen = std::collections::HashSet::new();

    for a in discovered {
        seen.insert(a.name.clone());
        result.push(a);
    }
    for a in builtin {
        if !seen.contains(&a.name) {
            result.push(a);
        }
    }

    result
}

fn find_agent<'a>(agents: &'a mut [Agent], name: &str) -> Option<&'a mut Agent> {
    agents.iter_mut().find(|a| a.name == name)
}

fn inject_context(agent: &mut Agent) {
    let mem_dir = memory_dir();
    let project = detect_project();

    let research_path = mem_dir
        .join("research")
        .join(&project)
        .join("research.toon");
    if research_path.exists() {
        agent.research = load_entries(
            &research_path,
            &agent.name,
            &["verified:", "finding:", "fact:"],
        );
    }

    let patterns_path = mem_dir
        .join("patterns")
        .join(&project)
        .join("patterns.toon");
    if patterns_path.exists() {
        agent.patterns = load_entries(
            &patterns_path,
            &agent.name,
            &["pattern:", "solution:", "template:"],
        );
    }
}

fn load_entries(path: &Path, agent_name: &str, markers: &[&str]) -> Vec<String> {
    let content = match std::fs::read_to_string(path) {
        Ok(c) => c,
        Err(_) => return Vec::new(),
    };

    let lower_agent = agent_name.to_lowercase();
    content
        .lines()
        .map(|l| l.trim())
        .filter(|l| !l.is_empty() && !l.starts_with('#'))
        .filter(|l| {
            let lower = l.to_lowercase();
            lower.contains(&lower_agent) || markers.iter().any(|m| lower.contains(m))
        })
        .map(|l| l.to_string())
        .collect()
}

// --- Built-in Agents ---

fn builtin_agents() -> Vec<Agent> {
    vec![
        Agent {
            name: "nlu-intent-analyzer".into(),
            level: LEVEL_NLU,
            model: "haiku".into(),
            description: "Fast intent parsing - routes to CEO".into(),
            triggers: vec!["user_prompt".into()],
            tools: vec!["Read".into(), "Grep".into(), "Glob".into()],
            path: String::new(),
            research: Vec::new(),
            patterns: Vec::new(),
        },
        Agent {
            name: "ceo".into(),
            level: LEVEL_CEO,
            model: "opus".into(),
            description: "Orchestrator - delegates, never writes code".into(),
            triggers: vec!["task_delegation".into(), "complex_request".into()],
            tools: vec![
                "Task".into(),
                "Read".into(),
                "Grep".into(),
                "Glob".into(),
                "WebSearch".into(),
            ],
            path: String::new(),
            research: Vec::new(),
            patterns: Vec::new(),
        },
        Agent {
            name: "research-director".into(),
            level: LEVEL_CEO,
            model: "opus".into(),
            description: "Evidence-based research findings".into(),
            triggers: vec!["research_needed".into(), "verify_facts".into()],
            tools: vec![
                "WebSearch".into(),
                "WebFetch".into(),
                "Read".into(),
                "Grep".into(),
            ],
            path: String::new(),
            research: Vec::new(),
            patterns: Vec::new(),
        },
        Agent {
            name: "backend-engineer".into(),
            level: LEVEL_ENGINEER,
            model: "sonnet".into(),
            description: "Rust backend - Axum, Tonic, Zig".into(),
            triggers: vec!["backend_task".into(), "api_implementation".into()],
            tools: vec![
                "Read".into(),
                "Edit".into(),
                "Write".into(),
                "Bash".into(),
                "Grep".into(),
            ],
            path: String::new(),
            research: Vec::new(),
            patterns: Vec::new(),
        },
        Agent {
            name: "frontend-engineer".into(),
            level: LEVEL_ENGINEER,
            model: "sonnet".into(),
            description: "TypeScript + React frontend".into(),
            triggers: vec!["frontend_task".into(), "ui_implementation".into()],
            tools: vec![
                "Read".into(),
                "Edit".into(),
                "Write".into(),
                "Bash".into(),
                "Grep".into(),
            ],
            path: String::new(),
            research: Vec::new(),
            patterns: Vec::new(),
        },
        Agent {
            name: "devops-engineer".into(),
            level: LEVEL_ENGINEER,
            model: "sonnet".into(),
            description: "Docker, K8s, CI/CD pipelines".into(),
            triggers: vec!["devops_task".into(), "deployment".into()],
            tools: vec!["Read".into(), "Edit".into(), "Write".into(), "Bash".into()],
            path: String::new(),
            research: Vec::new(),
            patterns: Vec::new(),
        },
        Agent {
            name: "security-engineer".into(),
            level: LEVEL_ENGINEER,
            model: "sonnet".into(),
            description: "Security analysis, OWASP compliance".into(),
            triggers: vec!["security_review".into(), "vulnerability_check".into()],
            tools: vec!["Read".into(), "Grep".into(), "WebSearch".into()],
            path: String::new(),
            research: Vec::new(),
            patterns: Vec::new(),
        },
        Agent {
            name: "qa-lead".into(),
            level: LEVEL_ENGINEER,
            model: "sonnet".into(),
            description: "Test strategy and coverage".into(),
            triggers: vec!["testing_task".into(), "coverage_analysis".into()],
            tools: vec!["Read".into(), "Edit".into(), "Write".into(), "Bash".into()],
            path: String::new(),
            research: Vec::new(),
            patterns: Vec::new(),
        },
        Agent {
            name: "code-reviewer".into(),
            level: LEVEL_REVIEW,
            model: "sonnet".into(),
            description: "Post-implementation code review".into(),
            triggers: vec!["code_review".into(), "pr_review".into()],
            tools: vec!["Read".into(), "Grep".into(), "Bash".into()],
            path: String::new(),
            research: Vec::new(),
            patterns: Vec::new(),
        },
        Agent {
            name: "aegis-guardian".into(),
            level: LEVEL_AEGIS,
            model: "opus".into(),
            description: "Verification Guardian - Quality, Security, Testing".into(),
            triggers: vec!["final_verification".into(), "quality_gate".into()],
            tools: vec!["Read".into(), "Grep".into(), "Bash".into()],
            path: String::new(),
            research: Vec::new(),
            patterns: Vec::new(),
        },
    ]
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
