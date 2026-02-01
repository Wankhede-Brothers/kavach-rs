//! Session state management: load/save from TOON file.
//! Path: ~/.local/shared/shared-ai/stm/session-state.toon

use std::collections::HashMap;
use std::fs;
use std::io::Write;
use std::path::PathBuf;

pub struct SessionState {
    pub id: String,
    pub today: String,
    pub project: String,
    pub work_dir: String,
    pub research_done: bool,
    pub memory_queried: bool,
    pub ceo_invoked: bool,
    pub nlu_parsed: bool,
    pub training_cutoff: String,
    pub post_compact: bool,
    pub compact_count: i32,
    pub turn_count: i32,
    pub last_reinforce_turn: i32,
    pub reinforce_every_n: i32,
    pub current_task: String,
    pub session_id: String,
    pub intent_type: String,
    pub intent_domain: String,
    pub intent_sub_agents: Vec<String>,
    pub intent_skills: Vec<String>,
    pub task_status: String,
    pub research_topic: String,
    pub files_modified: Vec<String>,
    pub aegis_verified: bool,
    pub tasks_created: i32,
    pub tasks_completed: i32,
}

impl Default for SessionState {
    fn default() -> Self {
        let today = chrono::Local::now().format("%Y-%m-%d").to_string();
        let wd = std::env::current_dir()
            .map(|p| p.to_string_lossy().to_string())
            .unwrap_or_default();
        let id = generate_session_id(&wd);
        Self {
            id: id.clone(),
            today,
            project: detect_project(),
            work_dir: wd,
            research_done: false,
            memory_queried: false,
            ceo_invoked: false,
            nlu_parsed: false,
            training_cutoff: "2025-01".into(),
            post_compact: false,
            compact_count: 0,
            turn_count: 0,
            last_reinforce_turn: 0,
            reinforce_every_n: 15,
            current_task: String::new(),
            session_id: id,
            intent_type: String::new(),
            intent_domain: String::new(),
            intent_sub_agents: Vec::new(),
            intent_skills: Vec::new(),
            task_status: String::new(),
            research_topic: String::new(),
            files_modified: Vec::new(),
            aegis_verified: false,
            tasks_created: 0,
            tasks_completed: 0,
        }
    }
}

fn generate_session_id(work_dir: &str) -> String {
    use std::hash::{Hash, Hasher};
    let mut hasher = std::collections::hash_map::DefaultHasher::new();
    work_dir.hash(&mut hasher);
    chrono::Local::now()
        .format("%Y%m%d")
        .to_string()
        .hash(&mut hasher);
    format!("sess_{:016x}", hasher.finish())
}

fn detect_project() -> String {
    if let Ok(wd) = std::env::current_dir() {
        if wd.join(".git").exists() {
            return wd
                .file_name()
                .map(|n| n.to_string_lossy().to_string())
                .unwrap_or_default();
        }
    }
    String::new()
}

pub fn state_path() -> PathBuf {
    let home = dirs::home_dir().unwrap_or_default();
    home.join(".local")
        .join("shared")
        .join("shared-ai")
        .join("stm")
        .join("session-state.toon")
}

pub fn load_session_state() -> Option<SessionState> {
    let path = state_path();
    let content = fs::read_to_string(&path).ok()?;
    let today = chrono::Local::now().format("%Y-%m-%d").to_string();

    let mut state = SessionState::default();
    let mut fields: HashMap<String, String> = HashMap::new();

    for line in content.lines() {
        let line = line.trim();
        if line.is_empty() || line.starts_with('#') || line.starts_with('[') {
            continue;
        }
        if let Some(idx) = line.find(':') {
            let key = line[..idx].trim().to_string();
            let value = line[idx + 1..].trim().to_string();
            fields.insert(key, value);
        }
    }

    if let Some(v) = fields.get("id") { state.id = v.clone(); }
    if let Some(v) = fields.get("today") { state.today = v.clone(); }
    if let Some(v) = fields.get("project") { state.project = v.clone(); }
    if let Some(v) = fields.get("workdir") { state.work_dir = v.clone(); }
    if let Some(v) = fields.get("research_done").or(fields.get("research")) {
        state.research_done = v == "true";
    }
    if let Some(v) = fields.get("memory") { state.memory_queried = v == "true"; }
    if let Some(v) = fields.get("ceo") { state.ceo_invoked = v == "true"; }
    if let Some(v) = fields.get("nlu") { state.nlu_parsed = v == "true"; }
    if let Some(v) = fields.get("cutoff") { state.training_cutoff = v.clone(); }
    if let Some(v) = fields.get("post_compact") { state.post_compact = v == "true"; }
    if let Some(v) = fields.get("compact_count") { state.compact_count = v.parse().unwrap_or(0); }
    if let Some(v) = fields.get("turn_count") { state.turn_count = v.parse().unwrap_or(0); }
    if let Some(v) = fields.get("last_reinforce_turn") { state.last_reinforce_turn = v.parse().unwrap_or(0); }
    if let Some(v) = fields.get("reinforce_every_n") { state.reinforce_every_n = v.parse().unwrap_or(15); }
    if let Some(v) = fields.get("session_id") { state.session_id = v.clone(); }
    if let Some(v) = fields.get("task") { state.current_task = v.clone(); }
    if let Some(v) = fields.get("type") { state.intent_type = v.clone(); }
    if let Some(v) = fields.get("domain") { state.intent_domain = v.clone(); }
    if let Some(v) = fields.get("subagents") {
        state.intent_sub_agents = split_csv(v);
    }
    if let Some(v) = fields.get("skills") {
        state.intent_skills = split_csv(v);
    }
    if let Some(v) = fields.get("task_status") { state.task_status = v.clone(); }
    if let Some(v) = fields.get("research_topic") { state.research_topic = v.clone(); }
    if let Some(v) = fields.get("aegis") { state.aegis_verified = v == "true"; }
    if let Some(v) = fields.get("tasks_created") { state.tasks_created = v.parse().unwrap_or(0); }
    if let Some(v) = fields.get("tasks_completed") { state.tasks_completed = v.parse().unwrap_or(0); }

    if state.today != today {
        return None;
    }

    Some(state)
}

pub fn get_or_create_session() -> SessionState {
    load_session_state().unwrap_or_default()
}

impl SessionState {
    pub fn save(&self) -> anyhow::Result<()> {
        let path = state_path();
        if let Some(parent) = path.parent() {
            fs::create_dir_all(parent)?;
        }

        let tmp_path = path.with_extension("toon.tmp");
        let mut f = fs::File::create(&tmp_path)?;

        writeln!(f, "# Session State - SP/1.0")?;
        writeln!(f, "# Auto-generated, do not edit")?;
        writeln!(f)?;
        writeln!(f, "[SESSION]")?;
        writeln!(f, "id: {}", self.id)?;
        writeln!(f, "today: {}", self.today)?;
        writeln!(f, "project: {}", self.project)?;
        writeln!(f, "workdir: {}", self.work_dir)?;
        writeln!(f, "cutoff: {}", self.training_cutoff)?;
        writeln!(f)?;
        writeln!(f, "[STATE]")?;
        writeln!(f, "research_done: {}", self.research_done)?;
        writeln!(f, "memory: {}", self.memory_queried)?;
        writeln!(f, "ceo: {}", self.ceo_invoked)?;
        writeln!(f, "nlu: {}", self.nlu_parsed)?;
        writeln!(f, "turn_count: {}", self.turn_count)?;
        writeln!(f, "last_reinforce_turn: {}", self.last_reinforce_turn)?;
        writeln!(f, "reinforce_every_n: {}", self.reinforce_every_n)?;
        writeln!(f, "session_id: {}", self.session_id)?;
        writeln!(f)?;
        writeln!(f, "[COMPACT]")?;
        writeln!(f, "post_compact: {}", self.post_compact)?;
        writeln!(f, "compact_count: {}", self.compact_count)?;
        writeln!(f)?;
        writeln!(f, "aegis: {}", self.aegis_verified)?;
        writeln!(f, "tasks_created: {}", self.tasks_created)?;
        writeln!(f, "tasks_completed: {}", self.tasks_completed)?;
        writeln!(f, "research_topic: {}", self.research_topic)?;
        writeln!(f)?;
        writeln!(f, "[TASK]")?;
        writeln!(f, "task: {}", self.current_task)?;
        writeln!(f, "task_status: {}", self.task_status)?;
        writeln!(f)?;
        if !self.intent_type.is_empty() {
            writeln!(f, "[INTENT_BRIDGE]")?;
            writeln!(f, "type: {}", self.intent_type)?;
            writeln!(f, "domain: {}", self.intent_domain)?;
            if !self.intent_sub_agents.is_empty() {
                writeln!(f, "subagents: {}", self.intent_sub_agents.join(","))?;
            }
            if !self.intent_skills.is_empty() {
                writeln!(f, "skills: {}", self.intent_skills.join(","))?;
            }
        }

        drop(f);
        fs::rename(&tmp_path, &path)?;
        Ok(())
    }

    pub fn increment_turn(&mut self) {
        self.turn_count += 1;
        let _ = self.save();
    }

    pub fn reset_research_for_new_prompt(&mut self) {
        if !self.has_task() {
            self.research_done = false;
            self.ceo_invoked = false;
            let _ = self.save();
        }
    }

    pub fn mark_nlu_parsed(&mut self) {
        self.nlu_parsed = true;
        let _ = self.save();
    }

    pub fn mark_ceo_invoked(&mut self) {
        self.ceo_invoked = true;
        let _ = self.save();
    }

    pub fn is_post_compact(&self) -> bool {
        self.post_compact
    }

    pub fn clear_post_compact(&mut self) {
        self.post_compact = false;
        let _ = self.save();
    }

    pub fn needs_reinforcement(&self) -> bool {
        let threshold = if self.reinforce_every_n > 0 {
            self.reinforce_every_n
        } else {
            15
        };
        self.turn_count - self.last_reinforce_turn >= threshold
    }

    pub fn mark_reinforcement_done(&mut self) {
        self.last_reinforce_turn = self.turn_count;
        let _ = self.save();
    }

    pub fn has_task(&self) -> bool {
        !self.current_task.is_empty()
    }

    pub fn store_intent(
        &mut self,
        intent_type: &str,
        domain: &str,
        sub_agents: &[String],
        skills: &[String],
    ) {
        self.intent_type = intent_type.into();
        self.intent_domain = domain.into();
        self.intent_sub_agents = sub_agents.to_vec();
        self.intent_skills = skills.to_vec();
        let _ = self.save();
    }

    pub fn mark_post_compact(&mut self) {
        self.post_compact = true;
        self.compact_count += 1;
        let _ = self.save();
    }

    pub fn mark_memory_queried(&mut self) {
        self.memory_queried = true;
        let _ = self.save();
    }

    pub fn mark_research_done(&mut self) {
        self.research_done = true;
        let _ = self.save();
    }

    pub fn mark_research_done_with_topic(&mut self, topic: &str) {
        self.research_done = true;
        if !topic.is_empty() {
            self.research_topic = topic.into();
        }
        let _ = self.save();
    }

    pub fn set_current_task(&mut self, task: &str) {
        self.current_task = task.into();
        self.task_status = "in_progress".into();
        let _ = self.save();
    }

    pub fn clear_task(&mut self) {
        self.current_task.clear();
        self.task_status.clear();
        let _ = self.save();
    }

    pub fn add_file_modified(&mut self, path: &str) {
        if !self.files_modified.contains(&path.to_string()) {
            self.files_modified.push(path.into());
        }
    }
}

fn split_csv(s: &str) -> Vec<String> {
    s.split(',')
        .map(|p| p.trim().to_string())
        .filter(|p| !p.is_empty())
        .collect()
}
