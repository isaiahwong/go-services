package main

import (
	"html/template"
	"log"
	"os"
)

type data struct {
	Name string
}

const prototemplate = `
package server

import (
	"context"

	pb "github.com/isaiahwong/go-services/src/gateway/proto-gen/payment"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)


func getProtos() []func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error {
	return []func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error{
		pb.Register{{.Name}}Handler,
	}
}
`

func main() {
	var d data = data{"PaymentService"}
	f, err := os.Create("./internal/server/protos.go")
	if err != nil {
		log.Fatalf("failed with %s\n", err)
	}
	t := template.Must(template.New("server").Parse(prototemplate))
	t.Execute(f, d)
}
