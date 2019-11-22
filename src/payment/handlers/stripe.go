package handlers

import (
	"context"
	pb "github.com/isaiahwong/go-services/src/payment/proto-gen/payment"
)

func (*PaymentService) StripeSetupIntent(context.Context, *pb.SetupIntentRequest) (*pb.SetupIntentResponse, error) {
	return &pb.SetupIntentResponse{}, nil
}
func (*PaymentService) StripeAddCard(context.Context, *pb.AddCardRequest) (*pb.AddCardResponse, error) {
	return &pb.AddCardResponse{}, nil
}
func (*PaymentService) StripeCharge(context.Context, *pb.OnStripeChargeRequest) (*pb.OnStripeChargeResponse, error) {
	return &pb.OnStripeChargeResponse{}, nil
}
func (*PaymentService) StripeTestWebhook(context.Context, *pb.StripeWebhook) (*pb.Response, error) {
	return &pb.Response{}, nil
}
func (*PaymentService) StripePaymentIntentWebhook(context.Context, *pb.StripeWebhook) (*pb.Response, error) {
	return &pb.Response{}, nil
}
