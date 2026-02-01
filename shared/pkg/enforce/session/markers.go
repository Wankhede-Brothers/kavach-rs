// Package session provides session state management.
// markers.go: State marking functions.
// DACE: Single responsibility - marker functions only.
package session

import "time"

// MarkResearchDone marks that WebSearch was performed (backward compat).
func (s *SessionState) MarkResearchDone() {
	s.ResearchDone = true
	s.Save()
}

// MarkResearchDoneWithTopic marks research done and records the topic.
func (s *SessionState) MarkResearchDoneWithTopic(topic string) {
	s.ResearchDone = true
	if topic != "" {
		s.ResearchTopics = append(s.ResearchTopics, topic)
	}
	s.Save()
}

// MarkMemoryQueried marks that memory bank was queried.
func (s *SessionState) MarkMemoryQueried() {
	s.MemoryQueried = true
	s.Save()
}

// MarkCEOInvoked marks that CEO was invoked.
func (s *SessionState) MarkCEOInvoked() {
	s.CEOInvoked = true
	s.Save()
}

// MarkNLUParsed marks that NLU parsing was done.
// Called by: intent gate after classifyIntentFromConfig().
func (s *SessionState) MarkNLUParsed() {
	s.NLUParsed = true
	s.Save()
}

// MarkAegisVerified marks that Aegis verification passed.
// Called by: aegis gate after verification passes.
func (s *SessionState) MarkAegisVerified() {
	s.AegisVerified = true
	s.Save()
}

// MarkPostCompact marks that context was compacted.
func (s *SessionState) MarkPostCompact() {
	s.PostCompact = true
	s.CompactedAt = time.Now().Format("2006-01-02T15:04:05")
	s.CompactCount++
	s.Save()
}

// ClearPostCompact clears post-compact mode after resume.
func (s *SessionState) ClearPostCompact() {
	s.PostCompact = false
	s.Save()
}

// IsPostCompact returns true if in post-compact mode.
func (s *SessionState) IsPostCompact() bool {
	return s.PostCompact
}

// IncrementTurn increments the conversation turn counter.
// Called on each UserPromptSubmit.
func (s *SessionState) IncrementTurn() {
	s.TurnCount++
	s.Save()
}

// ResetResearchForNewPrompt resets ResearchDone and CEOInvoked on each new user prompt,
// but only when no explicit task is active (preserves task-scoped behavior).
func (s *SessionState) ResetResearchForNewPrompt() {
	if !s.HasTask() {
		s.ResearchDone = false
		s.ResearchTopics = nil
		s.CEOInvoked = false
		s.Save()
	}
}

// NeedsReinforcement returns true if context reinforcement is needed.
// This mitigates "Lost in the Middle" attention decay in long conversations.
// Default: reinforce every 15 turns (~50k tokens estimated).
func (s *SessionState) NeedsReinforcement() bool {
	threshold := s.ReinforceEveryN
	if threshold == 0 {
		threshold = 15 // Default: every 15 turns
	}
	return s.TurnCount-s.LastReinforceTurn >= threshold
}

// MarkReinforcementDone marks that context was reinforced.
func (s *SessionState) MarkReinforcementDone() {
	s.LastReinforceTurn = s.TurnCount
	s.Save()
}

// SetCurrentTask sets a new task and resets task-scoped research state.
// Called by: task gate on TaskCreate to scope research per task.
func (s *SessionState) SetCurrentTask(task string) {
	if s.CurrentTask != task && task != "" {
		s.CurrentTask = task
		s.ResearchDone = false
		s.AegisVerified = false
		s.TaskStatus = "in_progress"
		s.Save()
	}
}

// MarkSpecInjected records that a spec file was injected this session.
func (s *SessionState) MarkSpecInjected(name string) {
	for _, n := range s.SpecsInjected {
		if n == name {
			return
		}
	}
	s.SpecsInjected = append(s.SpecsInjected, name)
	s.Save()
}

// WasSpecInjected returns true if the named spec was already injected.
func (s *SessionState) WasSpecInjected(name string) bool {
	for _, n := range s.SpecsInjected {
		if n == name {
			return true
		}
	}
	return false
}

// StoreIntent persists intent classification for the CEO gate to read.
func (s *SessionState) StoreIntent(intentType, domain string, subAgents, skills []string) {
	s.IntentType = intentType
	s.IntentDomain = domain
	s.IntentSubAgents = subAgents
	s.IntentSkills = skills
	s.Save()
}
