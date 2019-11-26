package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/isaiahwong/go-services/src/payment/proto-gen/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	opts := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	c := pb.NewPaymentServiceClient(conn)
	var trailer metadata.MD // variable to store header and trailer

	res, err := c.CreatePayment(context.Background(), &pb.CreatePaymentRequest{
		Email: "isaiah@jirehsoho.com",
		User:  "s",
	}, grpc.Trailer(&trailer))

	e := trailer.Get("errors-bin")

	log.Println(e)

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Read block %v\n", res)
}
