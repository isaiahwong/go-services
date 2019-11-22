package handlers

import (
	"context"
	"github.com/isaiahwong/go-services/src/payment/model"
	pb "github.com/isaiahwong/go-services/src/payment/proto-gen/payment"
)

func (*PaymentService) PaypalCreateOrder(context.Context, *pb.PaypalCreateOrderRequest) (*pb.PaypalCreateOrderResponse, error) {
	return &pb.PaypalCreateOrderResponse{}, nil
}
func (*PaymentService) PaypalProcessOrder(context.Context, *pb.PaypalProcessOrderRequest) (*pb.Response, error) {
	return &pb.Response{}, nil
}
func (*PaymentService) PaypalOrderWebhook(context.Context, *pb.PaypalWebhook) (*pb.Response, error) {
	return &pb.Response{}, nil
}
func (p *PaymentService) PaypalTestWebhook(context.Context, *pb.PaypalWebhook) (*pb.Response, error) {
	pay := model.Payment{
		Object: "payment",
		User:   "123",
		Email:  "isaiah@jirehsoho.com",
	}
	p.store.Create(context.Background(), &pay)

	return &pb.Response{Success: true}, nil
}
