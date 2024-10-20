package auditrail

import "context"

type discardLogger struct{}

// NewDiscardLogger returns a logger that discards all log entries.
func NewDiscardLogger() Logger {
	return &discardLogger{}
}

func (n *discardLogger) Log(context.Context, *Entry) error {
	return nil
}
