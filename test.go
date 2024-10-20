package auditrail

import (
	"context"
	"sync"
)

var _ Logger = (*TestLogger)(nil)

// TestLogger is a test spy that records all the logs it receives for later
// inspection. Useful for testing.
//
// This logger is safe for concurrent use and can be shared across multiple
// goroutines.
type TestLogger struct {
	logs []Entry
	mu   sync.RWMutex
}

// Log records the log event.
func (s *TestLogger) Log(_ context.Context, event *Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logs = append(s.logs, *event)

	return nil
}

// Size returns the number of logs recorded.
func (s *TestLogger) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.logs)
}

// Trail returns all the logs recorded.
func (s *TestLogger) Trail() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Entry, len(s.logs))
	copy(out, s.logs)

	return out
}

// Flush clears all the logs recorded.
func (s *TestLogger) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logs = nil
}
