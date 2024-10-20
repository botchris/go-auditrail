package httpd

import (
	"context"

	"github.com/botchris/auditrail"
)

type httpDecorator struct {
	inner auditrail.Logger
}

// Decorator returns a new audit.Logger that appends http details to the log
// entry before logging it.
//
// This decorator assumes that http data was previously added to the context
// using the [AddToContext] function, either directly or through some of the
// provided middlewares (Gin, Echo, etc).
func Decorator(inner auditrail.Logger) auditrail.Logger {
	return httpDecorator{inner: inner}
}

func (h httpDecorator) Log(ctx context.Context, entry *auditrail.Entry) error {
	d := FromContext(ctx)
	if d.IsEmpty() {
		return h.inner.Log(ctx, entry)
	}

	entry.AppendDetails("http", d)

	return h.inner.Log(ctx, entry)
}
