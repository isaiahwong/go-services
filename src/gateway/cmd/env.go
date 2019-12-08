package cmd

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
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
		log.Println(".env not loaded", err)
	}

	config = &EnvConfig{
		AppEnv:            mapEnvWithDefaults("APP_ENV", "development"),
		Production:        mapEnvWithDefaults("APP_ENV", "development") == "true",
		Port:              mapEnvWithDefaults("PORT", "5000"),
		WebhookPort:       mapEnvWithDefaults("WEBHOOK_PORT", "8443"),
		WebhookKeyDir:     mapEnvWithDefaults("WEBHOOK_KEY_DIR", "/run/secrets/tls/tls.key"),
		WebhookCertDir:    mapEnvWithDefaults("WEBHOOK_CERT_DIR", "/run/secrets/tls/tls.crt"),
		EnableStackdriver: mapEnvWithDefaults("ENABLE_STACKDRIVER", "true") == "true",
	}
	if !config.Production {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
}
