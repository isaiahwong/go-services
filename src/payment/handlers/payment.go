package handlers

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/isaiahwong/go-services/src/payment/model"
	pb "github.com/isaiahwong/go-services/src/payment/proto-gen/payment"
)

// CreatePayment TODO
func (p *PaymentService) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	user := strings.TrimSpace(req.GetUser())
	email := strings.TrimSpace(req.GetEmail())

	us := p.v.Var(user, "required,max=30")
	if us != nil {
		return nil, status.Errorf(codes.InvalidArgument, "User id required")
	}

	ee := p.v.Var(email, "required,email,max=64")
	if ee != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid email")
	}

	oid, err := p.store.Create(context.Background(), &model.Payment{
		User:  user,
		Email: email,
	})

	if err != nil {
		p.logger.Error(err)
		return nil, status.Error(codes.Internal, "An Internal error has occurred")
	}

	return &pb.CreatePaymentResponse{Success: true, Payment: &pb.Payment{Id: oid.String()}}, nil
}

// RetrievePayment TODO
func (*PaymentService) RetrievePayment(context.Context, *pb.RetrievePaymentRequest) (*pb.RetrievePaymentResponse, error) {
	return &pb.RetrievePaymentResponse{}, nil
}
func (*PaymentService) Refund(context.Context, *pb.RefundRequest) (*pb.RefundResponse, error) {
	return &pb.RefundResponse{}, nil
}
