package main

import (
	"context"
	"log"
	"net"

	pb "github.com/isaiahwong/go-services/src/test_/health" // Update
	"google.golang.org/grpc"


const (
	port = ":9090"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
}

// SayHello implements helloworld.GreeterServer
func (s *server) Echo(ctx context.Context, in *pb.Test) (*pb.Test, error) {
	log.Printf("Received: %v", in.GetValue())
	return &pb.Test{Value: "Hello " + in.GetValue()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
