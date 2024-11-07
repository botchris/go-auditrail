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
	logs         map[string]*Entry
	closeChannel chan struct{}
	closed       bool
	mu           sync.RWMutex
}

// NewMemoryLogger creates a new MemoryLogger.
func NewMemoryLogger() *MemoryLogger {
	return &MemoryLogger{
		logs:         make(map[string]*Entry),
		closeChannel: make(chan struct{}),
	}
}

// Log records the log event.
func (s *MemoryLogger) Log(_ context.Context, event *Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrTrailClosed
	}

	s.logs[event.GetIdempotencyID()] = event

	return nil
}

func (s *MemoryLogger) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true

	close(s.closeChannel)

	return nil
}

func (s *MemoryLogger) Closed() <-chan struct{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.closeChannel
}

func (s *MemoryLogger) IsClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.closed
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
		if _, ok := s.logs[id]; !ok {
			return false
		}
	}

	return true
}

// Trail returns all the logs recorded.
func (s *MemoryLogger) Trail() []*Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]*Entry, 0)

	for _, lb := range s.logs {
		out = append(out, lb)
	}

	return out
}

// Flush clears all the logs recorded.
func (s *MemoryLogger) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logs = nil
	s.logs = make(map[string]*Entry)
}
