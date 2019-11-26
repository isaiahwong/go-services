package store

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a new unary server interceptors that checks connection to MongoDB
func UnaryServerInterceptor(m *MongoStore) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		pe := m.Ping()
		if pe != nil {
			// Reconnect
			return nil, status.Error(codes.Internal, "Ping pong ping ring")
		}
		resp, err := handler(ctx, req)
		return resp, err
	}
}
