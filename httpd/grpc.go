package httpd

import (
	"context"

	"github.com/botchris/auditrail/grpcx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	_ grpc.UnaryServerInterceptor  = GRPCUnaryInterceptor
	_ grpc.StreamServerInterceptor = GRPCStreamInterceptor
)

// GRPCUnaryInterceptor is a gRPC unary call interceptor that injects
// [httpd.Details] into the context.
func GRPCUnaryInterceptor(ctx context.Context, req interface{}, usi *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	d := Details{
		Method:     md.Get(":method")[0],
		StatusCode: "",
		UserAgent:  md.Get("user-agent")[0],
		URL: URL{
			Host: md.Get(":authority")[0],
			Path: usi.FullMethod,
		},
	}

	return handler(AddToContext(ctx, d), req)
}

// GRPCStreamInterceptor is a gRPC stream server interceptor that injects
// [httpd.Details] into the context.
func GRPCStreamInterceptor(srv interface{}, ss grpc.ServerStream, ssi *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return handler(srv, ss)
	}

	d := Details{
		Method:     md.Get(":method")[0],
		StatusCode: "",
		UserAgent:  md.Get("user-agent")[0],
		URL: URL{
			Host: md.Get(":authority")[0],
			Path: ssi.FullMethod,
		},
	}

	stream := grpcx.ServerStreamWithContext(AddToContext(ctx, d), ss)

	return handler(srv, stream)
}
