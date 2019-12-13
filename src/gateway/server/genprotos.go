package server

import (
	"context"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func protos() []func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error {
	return []func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error{}
}
