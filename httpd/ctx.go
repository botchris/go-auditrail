package httpd

import "context"

type httpKeyType int

const httpKey httpKeyType = 0

// FromContext extracts the HTTP details from a context.
func FromContext(ctx context.Context) Details {
	if d, ok := ctx.Value(httpKey).(Details); ok {
		return d
	}

	return Details{}
}

// AddToContext adds the HTTP details to a context.
func AddToContext(ctx context.Context, d Details) context.Context {
	return context.WithValue(ctx, httpKey, d)
}
