package logtag_grpc

import (
	"context"
	"io"

	"github.com/roboticeyes/go-logtag/logtag"
	"google.golang.org/grpc"
)

func GrpcLogTagServerUnaryInterceptor(logTag string, logPayload bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		logtag.Printf(logTag, "↘️ %s: %s", info.FullMethod, req)
		// Calls the handler
		h, err := handler(ctx, req)

		if err != nil {
			logtag.Errorf(logTag, "↗️ %s: %s", info.FullMethod, logtag.ToColoredText(logtag.Red, err.Error()))
		} else if logPayload {
			logtag.Printf(logTag, "↘️ %s: %.500s", info.FullMethod, h)
		} else {
			logtag.Printf(logTag, "↘️ %s: <payload truncated>", info.FullMethod)
		}

		return h, err
	}
}

func GrpcLogTagClientUnaryInterceptor(logTag string, logPayload bool) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		logtag.Printf(logTag, "↗️ %s: %s", method, req)
		// Calls the handler
		err := invoker(ctx, method, req, reply, cc, opts...)

		if err != nil {
			logtag.Errorf(logTag, "↘️ %s: %s", method, logtag.ToColoredText(logtag.Red, err.Error()))
		} else if logPayload {
			logtag.Printf(logTag, "↘️ %s: %.500s", method, reply)
		} else {
			logtag.Printf(logTag, "↘️ %s: <payload truncated>", method)
		}

		return err
	}
}

func GrpcLogTagServerStreamInterceptor(logTag string) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		logtag.Printf(logTag, "↘️ %s: streaming started (client streaming: %t, server streaming: %t)", info.FullMethod, info.IsClientStream, info.IsServerStream)
		// Calls the handler
		err := handler(srv, &serverStreamMsgInterceptor{ServerStream: ss, tag: logTag, info: info})
		if err == io.EOF {
			return err
		}
		if err != nil {
			logtag.Errorf(logTag, "↗️ %s: %s", info.FullMethod, logtag.ToColoredText(logtag.Red, err.Error()))
		} else {
			logtag.Printf(logTag, "↗️ %s: streaming closed", info.FullMethod)
		}

		return err
	}
}

type serverStreamMsgInterceptor struct {
	grpc.ServerStream
	tag  string
	info *grpc.StreamServerInfo
}

func (s *serverStreamMsgInterceptor) SendMsg(m any) error {
	err := s.ServerStream.SendMsg(m)
	if err == io.EOF {
		return err
	}
	if err != nil {
		logtag.Errorf(s.tag, "↗️ %s: %s", s.info.FullMethod, logtag.ToColoredText(logtag.Red, err.Error()))
	} else if s.info.IsServerStream {
		logtag.Printf(s.tag, "↗️ %s: %s", s.info.FullMethod, m)
	}

	return err
}

func (s *serverStreamMsgInterceptor) RecvMsg(m any) error {
	err := s.ServerStream.RecvMsg(m)
	if err == io.EOF {
		return err
	}
	if err != nil {
		logtag.Errorf(s.tag, "↘️ %s: %s", s.info.FullMethod, logtag.ToColoredText(logtag.Red, err.Error()))
	} else if s.info.IsClientStream {
		logtag.Printf(s.tag, "↘️ %s: %s", s.info.FullMethod, m)
	}

	return err
}

// TODO GrpcLogTagClientStreamInterceptor Not tested, might need some love
func GrpcLogTagClientStreamInterceptor(logTag string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {

		logtag.Printf(logTag, "↘️ %s: streaming started  (client streaming: %t, server streaming: %t)", method, desc.ClientStreams, desc.ServerStreams)

		// Calls the handler
		clientStream, err := streamer(ctx, desc, cc, method, opts...)

		if err != nil {
			logtag.Errorf(logTag, "↗️ %s: %s", method, logtag.ToColoredText(logtag.Red, err.Error()))
			return nil, err
		} else {
			logtag.Printf(logTag, "↗️ %s: streaming closed", method)
		}
		return &clientStreamMsgInterceptor{ClientStream: clientStream, tag: logTag, desc: desc, method: method}, nil
	}
}

type clientStreamMsgInterceptor struct {
	grpc.ClientStream
	desc   *grpc.StreamDesc
	tag    string
	method string
}

func (c *clientStreamMsgInterceptor) SendMsg(m any) error {
	err := c.ClientStream.SendMsg(m)

	if err != nil {
		logtag.Errorf(c.tag, "%s: %s", c.method, logtag.ToColoredText(logtag.Red, err.Error()))
	} else if c.desc.ClientStreams {
		logtag.Printf(c.tag, "↗️ %s: %s", c.method, m)
	}

	return err
}

func (c *clientStreamMsgInterceptor) RecvMsg(m any) error {
	err := c.ClientStream.RecvMsg(m)

	if err != nil && err != io.EOF {
		logtag.Errorf(c.tag, "%s: %s", c.method, logtag.ToColoredText(logtag.Red, err.Error()))
	} else if c.desc.ServerStreams {
		logtag.Printf(c.tag, "↘️ %s: %s", c.method, m)
	}

	return err
}
