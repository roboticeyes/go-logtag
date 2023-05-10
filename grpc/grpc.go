package grpc

import (
	"context"

	"github.com/roboticeyes/go-logtag/logtag"
	"google.golang.org/grpc"
)

func GrpcLogTagServerInterceptor(logTag string) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		logtag.Printf(logTag, "↘️ %s: %s", info.FullMethod, req)
		// Calls the handler
		h, err := handler(ctx, req)

		if err != nil {
			logtag.Errorf(logTag, "↗️ %s: %s", info.FullMethod, logtag.ToColoredText(logtag.Red, err.Error()))
		} else {
			logtag.Printf(logTag, "↗️ %s: %s", info.FullMethod, h)
		}

		return h, err
	}
}

func GrpcLogTagClientInterceptor(logTag string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		logtag.Printf(logTag, "↗️ %s: %s", method, req)
		// Calls the handler
		err := invoker(ctx, method, req, reply, cc, opts...)

		if err != nil {
			logtag.Errorf(logTag, "↘️ %s: %s", method, logtag.ToColoredText(logtag.Red, err.Error()))
		} else {
			logtag.Printf(logTag, "↘️ %s: %s", method, reply)
		}

		return err
	}
}
