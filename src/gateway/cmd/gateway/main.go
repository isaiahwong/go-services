package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/isaiahwong/go-services/src/gateway/server"
	"github.com/isaiahwong/go-services/src/gateway/server/k8s"
	"github.com/isaiahwong/go-services/src/gateway/util/log"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// EnvConfig Application wide env configurations
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
	AppEnv           string
	Production       bool
	Host             string
	Port             string
	WebhookPort      string
	WebhookSecretKey string
	DisableK8S       bool

	WebhookKeyDir  string
	WebhookCertDir string

	EnableStackdriver bool

	DBUri      string
	DBName     string
	DBUser     string
	DBPassword string
	DBTimeout  time.Duration
}

// AppConfig config from EnvConfig
var config *EnvConfig

func mapEnvWithDefaults(envKey string, defaults string) string {
	v := os.Getenv(envKey)
	if v == "" {
		if defaults == "" {
			panic("defaults is not specified")
		}
		return defaults
	}
	return v
}

// LoadEnv loads envariables for AppConfig
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not loaded", err)
	}

	config = &EnvConfig{
		AppEnv:            mapEnvWithDefaults("APP_ENV", "development"),
		Production:        mapEnvWithDefaults("APP_ENV", "development") == "true",
		DisableK8S:        mapEnvWithDefaults("DISABLE_K8S_CLIENT", "false") == "true",
		Port:              mapEnvWithDefaults("PORT", "5000"),
		WebhookPort:       mapEnvWithDefaults("WEBHOOK_PORT", "8443"),
		WebhookKeyDir:     mapEnvWithDefaults("WEBHOOK_KEY_DIR", "/run/secrets/tls/tls.key"),
		WebhookCertDir:    mapEnvWithDefaults("WEBHOOK_CERT_DIR", "/run/secrets/tls/tls.crt"),
		EnableStackdriver: mapEnvWithDefaults("ENABLE_STACKDRIVER", "true") == "true",
	}
}

var logger *logrus.Logger

func init() {
	loadEnv()
	logger = log.NewLogger()
}

// Kills server gracefully
func gracefully(srvs []*http.Server) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutdown Servers ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for _, srv := range srvs {
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatal("Server Shutdown:", err)
		}
		// catching ctx.Done(). timeout of 5 seconds.
		select {
		case <-ctx.Done():
			logger.Println("timeout of 5 seconds.")
		}
		logger.Println("Server exiting")
	}
}

// Execute the entry point for gateway
func main() {
	var k *k8s.Client
	var err error

	if !config.DisableK8S {
		k, err = k8s.NewClient()
		if err != nil {
			logger.Fatalf("K8SClient Error: %v", err)
		}
	}

	gs, err := server.NewGatewayServer(config.Port,
		server.Logger(logger),
		server.K8SClient(k),
	)
	if err != nil {
		logger.Fatalf("New Gateway error: %v", err)
	}
	ws, err := server.NewWebhook(config.WebhookPort)
	if err != nil {
		logger.Fatalf("New Webhook error: %v", err)
	}

	// Registers gateway as an observer
	ws.Notifier.Register(gs)

	srvs := []*http.Server{
		gs.Server,
		// ws.Server,
	}

	go func() {
		// service connections
		if err := gs.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Gateway Server: %s\n", err)
		}
	}()

	// go func() {
	// 	// service connections
	// 	if err := ws.Server.ListenAndServeTLS(config.WebhookCertDir, config.WebhookKeyDir); err != nil && err != http.ErrServerClosed {
	// 		logger.Fatalf("Webhook server: %s\n", err)
	// 	}
	// }()

	gracefully(srvs)
}
