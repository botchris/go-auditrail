package auditrail

import (
	"context"
	"fmt"
)

// ErrTrailClosed is returned when the queue is closed.
var ErrTrailClosed = fmt.Errorf("trail is closed")

// Logger trail logger to which audit logs are written.
type Logger interface {
	// Log writes the given log entry to the audit log, returning an error
	// indicates that the log entry could not be written and should be retried.
	Log(context.Context, *Entry) error
}

// Closable is a trail logger that can be closed/stopped and waited for.
type Closable interface {
	Logger

	// Close closes the retryer and waits for all messages to be processed.
	Close() error

	// Closed returns a channel that is closed when the retryer is closed.
	Closed() <-chan struct{}

	// IsClosed returns true if the retryer is closed.
	IsClosed() bool
}
