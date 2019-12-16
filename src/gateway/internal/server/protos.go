package server

import (
	"context"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/isaiahwong/go-services/src/gateway/proto-gen/payment"
	"google.golang.org/grpc"
)

func getProtos() []func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error {
	return []func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error{
		pb.RegisterPaymentServiceHandler,
	}
}
