package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/isaiahwong/go-services/src/payment/proto-gen/payment"

	"google.golang.org/grpc"
)

func main() {
	opts := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	c := pb.NewPaymentServiceClient(conn)

	for i := 0; i < 5; i++ {
		res, err := c.CreatePayment(context.Background(), &pb.CreatePaymentRequest{
			Email: "isaiah@jirehsoho.com",
			User:  "12313123",
		})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Read block %v\n", res)
	}
}
