package auditrail

import (
	"context"
)

// Logger represents a component capable of logging audit entries.
type Logger interface {
	// Log writes the given log entry to the audit log, returning an error
	// indicates that the log entry could not be written and should be retried.
	Log(context.Context, *Entry) error
}
