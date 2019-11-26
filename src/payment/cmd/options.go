package cmd

import (
	"os"
	"strconv"
	"time"

	"github.com/isaiahwong/go-services/src/payment/pkg/log"
	"github.com/isaiahwong/go-services/src/payment/store"
	"github.com/joho/godotenv"
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

// EnvConfig Application wide env configurations
//
//
// AppEnv specifies if the app is in `development` or `production`
// Host specifies host address or dns
// Port specifies the port the server will run on
// EnableStackdriver specifies if google stackdriver will be enabled
// StripeSecret specifies Stripe api production key
// StripeSecretDev specifies Stripe api key for development
// StripeEndpointSecret specifies Stripe api key for webhook verification
// PaypalClientIDDev specifies Paypal api key for development
// PaypalSecretDev specifies Paypal api key secret for development
// PaypalClientID
// PaypalSecret         string
// PaypalURL specifies Paypal api URL for request
// DBUri
// DBUriDev
// DBUriTest
// DBName
// DBUser
// DBPassword
type EnvConfig struct {
	AppEnv               string
	Production           bool
	Host                 string
	Port                 string
	EnableStackdriver    bool
	StripeSecret         string
	StripeEndpointSecret string
	PaypalClientID       string
	PaypalSecret         string
	PaypalURL            string
	DBUri                string
	DBName               string
	DBUser               string
	DBPassword           string
	DBTimeout            time.Duration
}

// AppConfig config from EnvConfig
var config *EnvConfig

// LoadEnv loads envariables for AppConfig
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
		os.Exit(1)
	}

	// Convert to int
	sec, err := strconv.ParseInt(os.Getenv("DB_TIMEOUT"), 10, 64)
	if err != nil {
		panic(err)
	}

	dBTimeout := time.Duration(sec) * time.Second

	config = &EnvConfig{
		AppEnv:               os.Getenv("APP_ENV"),
		Production:           os.Getenv("APP_ENV") == "production",
		Host:                 os.Getenv("HOST"),
		Port:                 os.Getenv("PORT"),
		EnableStackdriver:    os.Getenv("ENABLE_STACKDRIVER") == "true",
		StripeSecret:         os.Getenv("STRIPE_SECRET"),
		StripeEndpointSecret: os.Getenv("STRIPE_ENDPOINT_SECRET"),
		PaypalClientID:       os.Getenv("PAYPAL_CLIENT_ID"),
		PaypalSecret:         os.Getenv("PAYPAL_SECRET"),
		PaypalURL:            os.Getenv("PAYPAL_URL"),
		DBUri:                os.Getenv("DB_URI"),
		DBName:               os.Getenv("DB_NAME"),
		DBUser:               os.Getenv("DB_USER"),
		DBPassword:           os.Getenv("DB_PASSWORD"),
		DBTimeout:            dBTimeout,
	}

	if !config.Production {
		config.StripeSecret = os.Getenv("STRIPE_SECRET_DEV")
		config.PaypalClientID = os.Getenv("PAYPAL_CLIENT_ID_DEV")
		config.PaypalSecret = os.Getenv("PAYPAL_SECRET_DEV")
		config.PaypalURL = os.Getenv("PAYPAL_URL_DEV")
		config.DBUri = os.Getenv("DB_URI_DEV")
	}

	if config.AppEnv == "test" {
		config.DBUri = os.Getenv("DB_URI_TEST")
	}
}
