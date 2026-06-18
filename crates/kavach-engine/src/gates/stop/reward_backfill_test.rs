//! Proof that the back-fill trigger is total and fail-closed: an empty session
//! id is a clean no-op (no RPC, no panic), and a populated one never panics even
//! when the daemon is absent (fire-and-forget). E5: a reward is banked ONLY on an
//! observed status transition (the internal `transition_observed` guard), never
//! on an allow-stop skip (user_focus supremacy or a foreign-tree turn).

use super::backfill_session_rewards;
use kavach_session::SessionState;

#[test]
fn an_empty_session_id_is_a_clean_noop() {
    let mut session = SessionState::default();
    assert!(session.session_id.is_empty(), "default has no session id");
    // Must not panic and must not attempt the RPC join (nothing to key on).
    backfill_session_rewards(&mut session);
}

#[test]
fn a_populated_session_never_panics_even_with_no_daemon() {
    let mut session = SessionState::default();
    session.session_id = "sess_test".to_owned();
    session.goal_receipt_pass = true;
    // Default has empty user_focus + empty work_dir ⇒ transition_observed = true.
    // Fire-and-forget: a down daemon is a silent no-op, never a panic.
    backfill_session_rewards(&mut session);
}

#[test]
fn an_allow_stop_skip_via_user_focus_banks_no_reward() {
    // E5: a verified-clean receipt is present, BUT the user pinned a scope
    // (user_focus supremacy) so the kanban was NOT drained — no transition. The
    // reward grader must NOT fire; the session total stays 0 (the bug banked +1).
    let mut session = SessionState::default();
    session.session_id = "sess_test".to_owned();
    session.goal_receipt_pass = true; // a clean receipt exists this turn
    session.user_focus = "fix the login bug".to_owned(); // pinned scope ⇒ allow-stop skip
    backfill_session_rewards(&mut session);
    assert_eq!(
        session.reward_session_total, 0,
        "user_focus allow-stop must bank no reward — gated on the status delta"
    );
    assert!(
        session.last_reward_summary.is_empty(),
        "no reward summary written on an allow-stop skip"
    );
}

#[test]
fn an_allow_stop_skip_via_foreign_tree_banks_no_reward() {
    // E5: every edit this turn is OUT-OF-TREE (e.g. editing harness source while
    // the card is rooted in the project) ⇒ the card cannot own the work ⇒ no
    // transition. card_owns_any_turn_file is false, so no reward is banked.
    let mut session = SessionState::default();
    session.session_id = "sess_test".to_owned();
    session.goal_receipt_pass = true;
    session.work_dir = "/proj/root".to_owned();
    session.files_modified_this_turn = vec!["/elsewhere/harness/file.rs".to_owned()];
    backfill_session_rewards(&mut session);
    assert_eq!(
        session.reward_session_total, 0,
        "foreign-tree allow-stop must bank no reward — the card cannot own the work"
    );
}

#[test]
fn an_observed_transition_with_clean_receipt_banks_a_reward() {
    // The mirror: a real transition (no pinned focus + an in-tree edit) WITH a
    // clean receipt grades the session — total increments to 1, pass to 1.
    let mut session = SessionState::default();
    session.session_id = "sess_test".to_owned();
    session.goal_receipt_pass = true;
    session.work_dir = "/proj/root".to_owned();
    session.files_modified_this_turn = vec!["/proj/root/src/lib.rs".to_owned()];
    backfill_session_rewards(&mut session);
    assert_eq!(session.reward_session_total, 1, "a real transition is a graded sample");
    assert_eq!(session.reward_session_pass, 1, "a clean receipt passes the sample");
}

#[test]
fn a_transition_without_a_clean_receipt_abstains() {
    // Boundary: transition observed but NO clean receipt ⇒ abstain (None),
    // never a -1. The total stays 0 (abstention is not a graded sample).
    let mut session = SessionState::default();
    session.session_id = "sess_test".to_owned();
    // empty user_focus + empty work_dir ⇒ transition_observed = true;
    // goal_receipt_pass defaults to false ⇒ outcome None ⇒ abstain.
    backfill_session_rewards(&mut session);
    assert_eq!(
        session.reward_session_total, 0,
        "no clean receipt ⇒ abstain, not a graded sample"
    );
}
