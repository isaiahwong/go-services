package handlers

import (
	"context"
	"encoding/json"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/isaiahwong/go-services/src/payment/model"
	"github.com/isaiahwong/go-services/src/payment/pkg/validator"
	pb "github.com/isaiahwong/go-services/src/payment/proto-gen/payment"
)

// CreatePayment TODO
func (p *PaymentService) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	user := strings.TrimSpace(req.GetUser())
	email := strings.TrimSpace(req.GetEmail())
	md := metadata.Pairs()

	errors := validator.Val(
		validator.Field{
			Param:   "email",
			Message: "Invalid email",
			Value:   email,
			Tag:     "required,email,max=64",
		},
		validator.Field{
			Param:   "user",
			Message: "User required",
			Value:   user,
			Tag:     "required,max=30",
		},
	)

	if errors != nil {
		json, jerr := json.Marshal(errors)
		if jerr != nil {
			return nil, status.Error(codes.Internal, "Unexpected error")
		}
		md.Append("errors-bin", string(json))
		grpc.SetTrailer(ctx, md)
		return nil, status.Error(codes.InvalidArgument, "Invalid arguments")
	}

	pay := &model.Payment{
		User:  user,
		Email: email,
	}

	// perr := p.val.Struct(pay)
	// if perr != nil {
	// 	for _, err := range p.valcast(perr) {

	// 		fmt.Println("\nFIELD")
	// 		fmt.Println(err.Namespace())
	// 		fmt.Println(err.Field())
	// 		fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
	// 		fmt.Println(err.StructField())     // by passing alt name to ReportError like below
	// 		fmt.Println(err.Tag())
	// 		fmt.Println(err.ActualTag())
	// 		fmt.Println(err.Kind())
	// 		fmt.Println(err.Type())
	// 		fmt.Println(err.StructField(), err.Value(), err.Param())
	// 	}
	// }

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
