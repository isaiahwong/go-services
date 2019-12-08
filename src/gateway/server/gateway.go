package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/isaiahwong/go-services/src/gateway/proto-gen/payment"
	"google.golang.org/grpc"
)

type gatewayOptions struct {
	apiGateway http.HandlerFunc
}

// GatewayOption sets options for GatewayServer.
type GatewayOption func(*gatewayOptions)

// GatewayServer encapsulates GatewayServer and Observer
type GatewayServer struct {
	Server   *http.Server
	services []APIService
}

// OnNotify receives events when being triggered
func (en *GatewayServer) OnNotify(e Event) {
	if e.Data == nil || len(e.Data) == 0 {
		return
	}

	s := &APIService{}
	// Apply routes
	constructAPI(e.Data, s)
}

func httpMiddleware(r *gin.Engine) error {
	r.Use(gin.Recovery())

	// Initialize a new gateway
	gw, err := newGateway(context.Background(), nil)
	if err != nil {
		return err
	}
	r.Any("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "success",
		})
	})
	// Proxies to gateway services
	r.Any("/api/*any", gin.WrapF(gw.ServeHTTP))
	r.NoRoute(notFound)

	return nil
}

func newRouter(attachMiddleware func(r *gin.Engine) error) (*gin.Engine, error) {
	r := gin.New()
	if attachMiddleware == nil {
		return r, nil
	}
	err := attachMiddleware(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func newGateway(ctx context.Context, opts []gwruntime.ServeMuxOption) (*gwruntime.ServeMux, error) {
	mux := gwruntime.NewServeMux(opts...)
	conn, err := grpc.DialContext(ctx, "localhost:50051", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	err = pb.RegisterPaymentServiceHandler(ctx, mux, conn)
	if err != nil {
		return nil, err
	}

	return mux, nil
}

// NewGatewayServer returns a new gin server
func NewGatewayServer(port string) (*GatewayServer, error) {
	r, err := newRouter(httpMiddleware)
	if err != nil {
		return nil, err
	}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: r,
	}

	gs := &GatewayServer{
		Server: s,
	}
	return gs, nil
}
