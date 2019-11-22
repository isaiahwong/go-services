package cmd

import (
	"context"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/isaiahwong/go-services/src/payment/handlers"
	pb "github.com/isaiahwong/go-services/src/payment/proto-gen/payment"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var appEnv string
var production bool

// Server encapsulates payment service
type Server struct {
	opts serverOptions
	gs   *grpc.Server
	lis  net.Listener
}

// NewServer Creates a new
func NewServer(opt ...ServerOption) *Server {
	opts := defaultServerOptions
	for _, o := range opt {
		o(&opts)
	}

	lis, err := net.Listen("tcp", opts.hostPort)
	if err != nil {
		opts.logger.Fatalf("Listening err: %v/n", err)
		return nil
	}

	// grpc.Creds()
	gs := grpc.NewServer(
		// grpc.Creds()
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(opts.logger)),
		)),
	)

	svc := handlers.NewPaymentService(opts.production, opts.logger, opts.store)

	// Registers gRPC services
	pb.RegisterPaymentServiceServer(gs, svc)

	return &Server{
		opts: opts,
		gs:   gs,
		lis:  lis,
	}
}

// Run Starts Payment Service
func (s *Server) Run() error {
	ce := s.opts.store.Connect(context.Background(), 60*time.Second)
	if ce != nil {
		return ce
	}

	s.opts.logger.Info("Starting Payment service")
	s.opts.logger.Infof("Production: %v", s.opts.production)
	if err := s.gs.Serve(s.lis); err != nil {
		return err
	}
	return nil

	// ch := make(chan os.Signal, 1)
	// signal.Notify(ch, os.Interrupt)

	// // Block until signal received
	// <-ch
	// logger.Println("Stopping server")
	// s.Stop()
	// lis.Close()
	// logger.Println("Exited")
}
