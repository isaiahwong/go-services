package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/isaiahwong/go-services/src/gateway/server"
	"github.com/isaiahwong/go-services/src/gateway/util/log"
	"github.com/sirupsen/logrus"
)

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
func Execute() {
	gs, err := server.NewGatewayServer(config.Port)
	ws, werr := server.NewWebhook(config.WebhookPort)
	if err != nil || werr != nil {
		// End the application
		panic(err)
	}

	// Registers gateway as an observer
	ws.Notifier.Register(gs)

	srvs := []*http.Server{
		gs.Server,
		ws.Server,
	}

	go func() {
		// service connections
		if err := gs.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Gateway Server: %s\n", err)
		}
	}()

	go func() {
		// service connections
		if err := ws.Server.ListenAndServeTLS(config.WebhookCertDir, config.WebhookKeyDir); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Webhook server: %s\n", err)
		}
	}()

	gracefully(srvs)
}
