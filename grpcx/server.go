package grpcx

import (
	"context"

	"google.golang.org/grpc"
)

// ServerStreamWithContext decorates the given ServerStream with provided
// context. Useful when building stream server interceptors that need to
// modify the context before calling the handler.
func ServerStreamWithContext(ctx context.Context, stream grpc.ServerStream) grpc.ServerStream {
	if existing, ok := stream.(*wrappedServerStream); ok {
		existing.ctx = ctx

		return existing
	}

	return &wrappedServerStream{
		ServerStream: stream,
		ctx:          ctx,
	}
}

// wrappedServerStream is a thin wrapper around grpc.ServerStream that allows
// modifying context.
type wrappedServerStream struct {
	grpc.ServerStream

	// ctx is the wrapper's own Context. You can assign it.
	ctx context.Context
}

// Context returns the wrapper's ctx, overwriting the nested
// `grpc.ServerStream.Context()`.
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
