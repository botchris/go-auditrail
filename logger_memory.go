package auditrail

import (
	"context"
	"sync"
)

var _ Logger = (*MemoryLogger)(nil)

// MemoryLogger is a test spy that records all the logs it receives for later
// inspection. Useful for testing.
//
// This logger is safe for concurrent use and can be shared across multiple
// goroutines.
type MemoryLogger struct {
	logs []Entry
	ids  map[string]struct{}
	mu   sync.RWMutex
}

// Log records the log event.
func (s *MemoryLogger) Log(_ context.Context, event *Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ids == nil {
		s.ids = make(map[string]struct{})
	}

	s.logs = append(s.logs, *event)
	s.ids[event.IdempotencyID] = struct{}{}

	return nil
}

// Size returns the number of logs recorded.
func (s *MemoryLogger) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.logs)
}

// Has whether the logger has recorded all the logs with the given idempotency
// IDs.
func (s *MemoryLogger) Has(idempotencyID ...string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, id := range idempotencyID {
		if _, ok := s.ids[id]; !ok {
			return false
		}
	}

	return true
}

// Trail returns all the logs recorded.
func (s *MemoryLogger) Trail() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Entry, len(s.logs))
	copy(out, s.logs)

	return out
}

// Flush clears all the logs recorded.
func (s *MemoryLogger) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logs = nil
	s.ids = make(map[string]struct{})
}
