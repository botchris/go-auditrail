package networkd

import (
	"context"

	"github.com/botchris/auditrail"
)

type clientDecorator struct {
	inner auditrail.Logger
	ipr   IPResolver
}

// Decorator returns a new audit.Logger that appends client details to the log
// entry before logging it.
//
// This decorator assumes that client data was previously added to the context
// using the AddToContext function.
//
// If the IPResolver is provided (not nil), it will be used to enrich the client
// details with GeoIP information.
func Decorator(inner auditrail.Logger, ipr IPResolver) auditrail.Logger {
	return clientDecorator{
		inner: inner,
		ipr:   ipr,
	}
}

func (h clientDecorator) Log(ctx context.Context, entry *auditrail.Entry) error {
	d := FromContext(ctx)
	if d == (Details{}) {
		return h.inner.Log(ctx, entry)
	}

	if d.Client.IP == "" {
		return h.inner.Log(ctx, entry)
	}

	if d.Client.GeoIP == nil && h.ipr != nil {
		gd := h.ipr.Resolve(d.Client.IP)
		d.Client.GeoIP = &gd

		AddToContext(ctx, d)
	}

	entry.AppendDetails("client", d)

	return h.inner.Log(ctx, entry)
}

func (h clientDecorator) Close() error {
	return h.inner.Close()
}

func (h clientDecorator) Closed() <-chan struct{} {
	return h.inner.Closed()
}

func (h clientDecorator) IsClosed() bool {
	return h.inner.IsClosed()
}
