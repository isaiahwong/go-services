package handlers

import (
	"context"
	"fmt"
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

	us := p.val.Var(user, "required,max=30")
	if us != nil {
		return nil, status.Errorf(codes.InvalidArgument, "User id required")
	}

	ee := p.val.Var(email, "required,email,max=64")
	if ee != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid email")
	}

	pay := &model.Payment{
		User:  user,
		Email: email,
	}

	errors := []model.Error{}
	perr := p.val.Struct(pay)

	if perr != nil {
		for _, err := range p.valcast(perr) {
			fmt.Println("\nFIELD")
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
			fmt.Println(err.StructField())     // by passing alt name to ReportError like below
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.StructField(), err.Value(), err.Param())
			append(errors, model.Error{})
		}
	}

	oid, err := p.store.Create(ctx, pay)

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
