package grpc

import (
	"context"
	"path"

	"google.golang.org/grpc"

	"gitlab.ozon.dev/iTukaev/homework/internal/counter"
)

func MetricsUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	method := path.Base(info.FullMethod)
	counter.Request.Inc(method)

	resp, err := handler(ctx, req)
	if err != nil {
		counter.Errors.Inc(method)
	} else {
		counter.Success.Inc(method)
	}
	counter.Response.Inc(method)
	return resp, err
}

func MetricsStreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	method := path.Base(info.FullMethod)
	counter.Request.Inc(method)

	err := handler(srv, ss)
	if err != nil {
		counter.Errors.Inc(method)
	} else {
		counter.Success.Inc(method)
	}
	counter.Response.Inc(method)
	return err
}
