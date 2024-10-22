package auditrail

import (
	"context"
	"sync"
)

type discardLogger struct {
	closed       bool
	closeChannel chan struct{}
	mu           sync.RWMutex
}

// NewDiscardLogger returns a logger that discards all log entries.
func NewDiscardLogger() Logger {
	return &discardLogger{
		closeChannel: make(chan struct{}),
	}
}

func (n *discardLogger) Log(context.Context, *Entry) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.closed {
		return ErrTrailClosed
	}

	return nil
}

func (n *discardLogger) Close() error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if n.closed {
		return nil
	}

	n.closed = true

	close(n.closeChannel)

	return nil

}

func (n *discardLogger) Closed() <-chan struct{} {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.closeChannel
}

func (n *discardLogger) IsClosed() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.closed
}
