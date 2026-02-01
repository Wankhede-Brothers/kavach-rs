//! Intent classification gate with multi-tier cascade.
//! Tier 0: Trivial -> silent exit (0 tokens)
//! Tier 1: Status query -> BINARY_FIRST directive (~50 tokens)
//! Tier 2: Post-compact recovery / periodic reinforcement (~80 tokens)
//! Tier 3: Full NLU classification (~200 tokens)

use crate::hook;
use crate::session;

pub fn run(hook_mode: bool) -> anyhow::Result<()> {
    if !hook_mode {
        crate::commands::cli_print("gates intent: use --hook flag");
        return Ok(());
    }

    let input = hook::must_read_hook_input();
    let result = run_intent_gate(&input);

    match result {
        Ok(()) => Ok(()),
        Err(e) if hook::is_hook_exit(&e) => Ok(()),
        Err(e) => Err(e),
    }
}

fn run_intent_gate(input: &hook::HookInput) -> anyhow::Result<()> {
    let prompt = input.get_prompt().to_lowercase().trim().to_string();

    if prompt.is_empty() {
        return hook::exit_user_prompt_submit_silent();
    }

    // TIER 0: Trivial prompts
    if is_simple_query(&prompt) {
        return hook::exit_user_prompt_submit_silent();
    }

    // TIER 1: Status queries
    if is_status_query(&prompt) {
        return hook::exit_user_prompt_submit_with_context(&status_directive());
    }

    // TIER 2+: Load session
    let mut sess = session::get_or_create_session();
    sess.increment_turn();
    sess.reset_research_for_new_prompt();
    let today = hook::today();

    let mut context_blocks: Vec<String> = Vec::new();

    // TIER 2: Post-compact recovery
    if sess.is_post_compact() {
        context_blocks.push(post_compact_recovery(&sess));
        sess.clear_post_compact();
        sess.mark_reinforcement_done();
    }

    // TIER 2: Periodic reinforcement
    if sess.needs_reinforcement() && !sess.is_post_compact() {
        context_blocks.push(periodic_reinforcement(&sess));
        sess.mark_reinforcement_done();
    }

    // TIER 3: NLU classification
    let intent = classify_intent(&prompt);
    if intent.intent_type != "unclassified" || !intent.domain.is_empty() {
        sess.mark_nlu_parsed();
        sess.store_intent(
            &intent.intent_type,
            &intent.domain,
            &intent.sub_agents,
            &intent.skills,
        );
        context_blocks.push(format_intent_directive(&intent, &today, &sess));
    }

    if !context_blocks.is_empty() {
        return hook::exit_user_prompt_submit_with_context(&context_blocks.join("\n\n"));
    }

    hook::exit_user_prompt_submit_silent()
}

// --- Types ---

struct IntentClassification {
    intent_type: String,
    domain: String,
    skills: Vec<String>,
    agent: String,
    sub_agents: Vec<String>,
    research_req: bool,
    confidence: String,
}

// --- Tier helpers ---

fn is_simple_query(prompt: &str) -> bool {
    let simple = [
        "hello", "hi", "hey", "thanks", "thank you", "bye", "yes", "no", "ok", "okay",
    ];
    simple.contains(&prompt)
}

fn is_status_query(prompt: &str) -> bool {
    let triggers = [
        "status",
        "project status",
        "what is the status",
        "show status",
        "check status",
    ];
    triggers.iter().any(|t| prompt.contains(t))
}

fn is_implementation_intent(t: &str) -> bool {
    ["implement", "debug", "refactor", "optimize", "fix", "audit", "docs", "unclassified"]
        .contains(&t)
}

fn is_delegation_required(t: &str) -> bool {
    ["implement", "debug", "refactor", "optimize", "audit"].contains(&t)
}

// --- NLU classification ---

fn classify_intent(prompt: &str) -> IntentClassification {
    let mut intent = IntentClassification {
        intent_type: String::new(),
        domain: String::new(),
        skills: Vec::new(),
        agent: String::new(),
        sub_agents: Vec::new(),
        research_req: true,
        confidence: "medium".into(),
    };

    // Intent type classification (priority order, first match wins)
    let debug_words = ["fix", "bug", "error", "broken", "crash", "failing", "not working"];
    let perf_words = ["optimize", "faster", "slow", "performance", "speed up"];
    let refactor_words = ["refactor", "restructure", "clean up", "technical debt"];
    let research_words = ["research", "explore", "explain", "how does", "what is"];
    let docs_words = ["document", "documentation", "readme", "api docs"];
    let audit_words = ["audit", "review", "vulnerability", "compliance"];
    let implement_words = ["implement", "create", "build", "add", "develop", "new feature"];

    if matches_any(prompt, &debug_words) {
        intent.intent_type = "debug".into();
        intent.skills = vec!["/debug-like-expert".into()];
        intent.agent = "ceo".into();
        intent.sub_agents = vec!["research-director".into(), "backend-engineer".into()];
        intent.confidence = "high".into();
    } else if matches_any(prompt, &perf_words) {
        intent.intent_type = "optimize".into();
        intent.skills = vec!["/dsa".into(), "/arch".into()];
        intent.agent = "ceo".into();
        intent.sub_agents = vec!["research-director".into(), "backend-engineer".into()];
        intent.confidence = "high".into();
    } else if matches_any(prompt, &refactor_words) {
        intent.intent_type = "refactor".into();
        intent.skills = vec!["/heal".into()];
        intent.agent = "ceo".into();
        intent.sub_agents = vec!["backend-engineer".into(), "aegis-guardian".into()];
    } else if matches_any(prompt, &research_words) {
        intent.intent_type = "research".into();
        intent.agent = "research-director".into();
        intent.confidence = "high".into();
    } else if matches_any(prompt, &docs_words) {
        intent.intent_type = "docs".into();
        intent.agent = "research-director".into();
        intent.sub_agents = vec!["backend-engineer".into()];
    } else if matches_any(prompt, &audit_words) {
        intent.intent_type = "audit".into();
        intent.skills = vec!["/security".into(), "/heal".into()];
        intent.agent = "ceo".into();
        intent.sub_agents = vec!["aegis-guardian".into()];
        intent.confidence = "high".into();
    } else if matches_any(prompt, &implement_words) {
        intent.intent_type = "implement".into();
        intent.agent = "ceo".into();
    }

    // Domain classification (additive)
    if matches_any(prompt, &["security", "auth", "encrypt", "oauth", "jwt"]) {
        if intent.domain.is_empty() {
            intent.domain = "security".into();
        }
        append_unique(&mut intent.skills, "/security".into());
    }
    if matches_any(prompt, &["frontend", "ui", "css", "react", "component"]) {
        if intent.domain.is_empty() {
            intent.domain = "frontend".into();
        }
        append_unique(&mut intent.skills, "/frontend".into());
        append_unique(&mut intent.sub_agents, "frontend-engineer".into());
    }
    if matches_any(prompt, &["database", "sql", "query", "migration", "postgres"]) {
        if intent.domain.is_empty() {
            intent.domain = "database".into();
        }
        append_unique(&mut intent.skills, "/sql".into());
    }
    if matches_any(prompt, &["deploy", "docker", "kubernetes", "k8s", "terraform", "infra"]) {
        if intent.domain.is_empty() {
            intent.domain = "infrastructure".into();
        }
        append_unique(&mut intent.skills, "/cloud-infrastructure-mastery".into());
    }
    if matches_any(prompt, &["test", "testing", "unit test", "integration test"]) {
        if intent.domain.is_empty() {
            intent.domain = "testing".into();
        }
        append_unique(&mut intent.skills, "/testing".into());
        append_unique(&mut intent.sub_agents, "aegis-guardian".into());
    }
    if matches_any(prompt, &["api", "endpoint", "rest", "grpc", "graphql"]) {
        if intent.domain.is_empty() {
            intent.domain = "backend".into();
        }
        append_unique(&mut intent.skills, "/api-design".into());
        append_unique(&mut intent.sub_agents, "backend-engineer".into());
    }

    // Default fallback
    if intent.intent_type.is_empty() && intent.domain.is_empty() {
        intent.intent_type = "unclassified".into();
        intent.confidence = "low".into();
    }
    if intent.agent.is_empty() {
        intent.agent = "ceo".into();
    }

    // Implementation intents require research
    if is_implementation_intent(&intent.intent_type) {
        intent.research_req = true;
    }

    intent
}

fn matches_any(text: &str, patterns: &[&str]) -> bool {
    patterns.iter().any(|p| text.contains(p))
}

fn append_unique(vec: &mut Vec<String>, item: String) {
    if !vec.contains(&item) {
        vec.push(item);
    }
}

// --- Output formatting ---

fn format_intent_directive(
    intent: &IntentClassification,
    today: &str,
    sess: &session::SessionState,
) -> String {
    let mut sb = String::new();

    // Compact header
    sb += &format!("[INTENT] type={}", intent.intent_type);
    if !intent.domain.is_empty() {
        sb += &format!(" domain={}", intent.domain);
    }
    sb += &format!(" confidence={} date={}\n", intent.confidence, today);

    // Research block (if needed and not done)
    if intent.research_req && !sess.research_done {
        sb += &format!(
            "[BLOCK:RESEARCH] BLOCKED: WebSearch required before implementation. Training weights are stale (cutoff: 2025-01). today:{}\n",
            today
        );
    }

    // Skill auto-invoke
    if !intent.skills.is_empty() {
        sb += "[SKILL:AUTO_INVOKE] MANDATORY:";
        for skill in &intent.skills {
            let name = skill.trim_start_matches('/');
            sb += &format!(" Skill(skill:\"{name}\")");
        }
        sb += "\n";
    }

    // Agent routing with delegation enforcement
    if intent.intent_type == "research" {
        sb += "[BLOCK:DELEGATION] MUST: Task(subagent_type='research-director') BEFORE any code\n";
    } else if is_delegation_required(&intent.intent_type) {
        sb += "[BLOCK:DELEGATION] MUST: Task(subagent_type='ceo') BEFORE Write/Edit\n";
    } else {
        sb += &format!("[AGENT] primary={}\n", intent.agent);
    }

    // DACE + forbidden phrases
    sb += "[DACE] max:100lines depth:5-7levels split:concern no:duplicates no:monoliths\n";

    sb
}

fn status_directive() -> String {
    "[BINARY_FIRST]\naction: IMMEDIATE\ncommand: kavach status && kavach memory bank\nFORBIDDEN: Task(Explore), Read(docs/*.md)\nreason: Memory Bank is SINGLE SOURCE OF TRUTH".into()
}

fn post_compact_recovery(sess: &session::SessionState) -> String {
    let today = hook::today();
    format!(
        "[RECOVERY] turn={} memory=kavach_memory_bank research=WebSearch_{} binary=kavach_FIRST dace=100lines_5depth",
        sess.turn_count, today
    )
}

fn periodic_reinforcement(sess: &session::SessionState) -> String {
    let today = hook::today();
    format!(
        "[REINFORCE] turn={} research={} dace=100lines_5depth fix=root_cause",
        sess.turn_count, today
    )
}
