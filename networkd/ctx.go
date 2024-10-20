package networkd

import "context"

type networkKeyType int

const networkKey networkKeyType = 0

// FromContext extracts the network details from the context.
func FromContext(ctx context.Context) Details {
	if d, ok := ctx.Value(networkKey).(Details); ok {
		return d
	}

	return Details{}
}

// AddToContext adds the network details to the context.
func AddToContext(ctx context.Context, d Details) context.Context {
	return context.WithValue(ctx, networkKey, d)
}
