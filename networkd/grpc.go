package networkd

import (
	"context"

	"github.com/botchris/go-auditrail/grpcx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var (
	_ grpc.UnaryServerInterceptor  = GRPCUnaryInterceptor
	_ grpc.StreamServerInterceptor = GRPCStreamInterceptor
)

// GRPCUnaryInterceptor is a gRPC unary call interceptor that injects
// [networkd.Details] into the context.
func GRPCUnaryInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	d := Details{
		Client: Client{
			IP: p.Addr.String(),
		},
	}

	return handler(AddToContext(ctx, d), req)
}

// GRPCStreamInterceptor is a gRPC stream server interceptor that injects
// [networkd.Details] into the context.
func GRPCStreamInterceptor(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	p, ok := peer.FromContext(ctx)

	if !ok {
		return handler(srv, ss)
	}

	d := Details{
		Client: Client{
			IP: p.Addr.String(),
		},
	}

	stream := grpcx.ServerStreamWithContext(AddToContext(ctx, d), ss)

	return handler(srv, stream)
}
