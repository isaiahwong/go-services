package cmd

import (
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/isaiahwong/go-services/src/payment/pkg/log"
	"github.com/isaiahwong/go-services/src/payment/store"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger
var datastore *store.MongoStore

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
		os.Exit(1)
	}
}

func init() {
	loadEnv()
	logger = log.NewLogger()
	s, me := store.NewMongoStore(
		store.ConnectionString(os.Getenv("DB_URI_DEV")),
		store.Database("payment"),
	)
	if me != nil {
		logger.Fatalf("NewMongoStore error: %v\n", me)
	}
	datastore = s
}

// Execute payment microservice
func Execute() {
	s := NewServer(
		HostPort(net.JoinHostPort("0.0.0.0", strconv.Itoa(50051))),
		Logger(logger),
		Production(strings.ToLower(os.Getenv("APP_ENV")) == "production"),
		Store(datastore),
	)
	if err := s.Run(); err != nil {
		s.opts.logger.Fatalf("Error serving server: %v", err)
	}
}
