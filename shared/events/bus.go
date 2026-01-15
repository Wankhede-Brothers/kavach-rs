package events

import (
	"sync"
	"time"
)

// EventType represents types of events
type EventType string

const (
	EventMemoryWrite  EventType = "memory_write"
	EventSessionStart EventType = "session_start"
	EventAgentInvoke  EventType = "agent_invoke"
	EventSkillInvoke  EventType = "skill_invoke"
	EventError        EventType = "error"
)

// Event represents a system event
type Event struct {
	Type      EventType   `json:"type"`
	Source    string      `json:"source"` // "opencode" or "claude-code"
	SessionID string      `json:"session_id,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

// EventBus manages event subscriptions and publishing
type EventBus struct {
	subscribers map[EventType][]chan Event
	mu          sync.RWMutex
}

var (
	globalBus *EventBus
	busOnce   sync.Once
)

// GetEventBus returns singleton event bus
func GetEventBus() *EventBus {
	busOnce.Do(func() {
		globalBus = &EventBus{
			subscribers: make(map[EventType][]chan Event),
			mu:          sync.RWMutex{},
		}
	})
	return globalBus
}

// Subscribe subscribes to events of a given type
func (bus *EventBus) Subscribe(eventType EventType) <-chan Event {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	ch := make(chan Event, 100)
	bus.subscribers[eventType] = append(bus.subscribers[eventType], ch)
	return ch
}

// maxPublishWorkers limits concurrent publish goroutines to prevent unbounded growth.
const maxPublishWorkers = 10

// Publish publishes an event to all subscribers.
// P0 FIX: Bounded goroutines with timeout to prevent leaks.
func (bus *EventBus) Publish(eventType EventType, source string, payload interface{}) {
	bus.mu.RLock()
	subs := bus.subscribers[eventType]
	bus.mu.RUnlock()

	if len(subs) == 0 {
		return
	}

	event := Event{
		Type:      eventType,
		Source:    source,
		Timestamp: time.Now(),
		Payload:   payload,
	}

	// Use semaphore to limit concurrent goroutines
	sem := make(chan struct{}, maxPublishWorkers)

	for _, ch := range subs {
		sem <- struct{}{} // Acquire semaphore slot
		go func(c chan Event) {
			defer func() { <-sem }() // Release semaphore slot
			// Non-blocking send with timeout to prevent goroutine leak
			select {
			case c <- event:
				// Success
			case <-time.After(100 * time.Millisecond):
				// Drop event if subscriber is blocked
			}
		}(ch)
	}
}
