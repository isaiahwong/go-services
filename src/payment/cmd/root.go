package cmd

import (
	"fmt"

	"github.com/isaiahwong/go-services/src/payment/pkg/log"
	"github.com/isaiahwong/go-services/src/payment/store"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger
var datastore *store.MongoStore

func init() {
	loadEnv()
	logger = log.NewLogger()

	// Create a new MongoStore
	s, me := store.NewMongoStore(
		store.ConnectionString(config.DBUri),
		store.Database(config.DBName),
		store.SetTimeout(config.DBTimeout),
		store.SetAuth(store.MongoCredential{
			Username: config.DBUser,
			Password: config.DBPassword,
		}),
		store.SetHeartbeat(config.DBTimeout),
	)
	if me != nil {
		logger.Fatalf("NewMongoStore error: %v\n", me)
	}
	datastore = s
}

// Execute payment microservice
func Execute() {
	s := NewServer(
		HostPort(fmt.Sprintf("%v:%v", config.Host, config.Port)),
		Logger(logger),
		Production(config.Production),
		Store(datastore),
	)
	if err := s.Run(config); err != nil {
		s.opts.logger.Fatalf("Error serving server: %v", err)
	}
}
