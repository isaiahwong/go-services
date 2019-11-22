package cmd

import (
	"github.com/isaiahwong/go-services/src/payment/pkg/log"
	"github.com/isaiahwong/go-services/src/payment/store"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials"
)

// A serverOptions sets options
type serverOptions struct {
	hostPort   string
	production bool
	logger     *logrus.Logger
	creds      credentials.TransportCredentials
	store      *store.MongoStore
}

var defaultServerOptions = serverOptions{
	hostPort:   "0.0.0.0:50051",
	production: false,
	logger:     log.NewLogger(),
}

// A ServerOption sets options such as hostPort; parameters, etc.
type ServerOption func(*serverOptions)

// HostPort returns a ServerOption that sets the address of the server
func HostPort(hostPort string) ServerOption {
	return func(o *serverOptions) {
		o.hostPort = hostPort
	}
}

// Production returns a ServerOption; determines if the server runs in production
func Production(p bool) ServerOption {
	return func(o *serverOptions) {
		o.production = p
	}
}

// Logger returns a ServerOption that sets the server logger
func Logger(l *logrus.Logger) ServerOption {
	return func(o *serverOptions) {
		o.logger = l
	}
}

// Store returns a ServerOption that sets the store
func Store(store *store.MongoStore) ServerOption {
	return func(o *serverOptions) {
		o.store = store
	}
}
