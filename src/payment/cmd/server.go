package cmd

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/isaiahwong/go-services/src/payment/handlers"
	pb "github.com/isaiahwong/go-services/src/payment/proto-gen/payment"
	"github.com/isaiahwong/go-services/src/payment/store"
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
			store.UnaryServerInterceptor(opts.store),
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
func (s *Server) Run(config *EnvConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DBTimeout)
	defer cancel()

	ce := s.opts.store.Connect(ctx)
	if ce != nil {
		return ce
	}

	s.opts.logger.Infof("Starting Payment service on %v", s.opts.hostPort)
	s.opts.logger.Infof("Production: %v", s.opts.production)

	if err := s.gs.Serve(s.lis); err != nil {
		return err
	}

	return nil
}

// Stop stops all services
func (s *Server) Stop() {
	s.gs.Stop()
	s.lis.Close()
	s.opts.store.Disconnect(context.Background())
}
