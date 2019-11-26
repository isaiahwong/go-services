package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentProvider type
type PaymentProvider string

//  Payment payment provider types
const (
	StripeProvider PaymentProvider = "stripe"
	PaypalProvider                 = "paypal"
)

// Payment struct
type Payment struct {
	ID              primitive.ObjectID `bson:"id,omitempty" json:"id"`
	Object          string             `validate:"eq=payment,required" bson:"object" json:"object" `
	User            string             `validate:"required" bson:"user" json:"user"`
	Email           string             `validate:"required,email" bson:"email" json:"email"  `
	DefaultProvider PaymentProvider    `validate:"required" bson:"default_provider" json:"default_provider"`
	Stripe          Stripe             `bson:"stripe" json:"stripe"`
	Updated         time.Time          `validate:"required" bson:"updated" json:"updated"`
	Created         time.Time          `validate:"required" bson:"created" json:"created"`
}

// NewPayment returns payment with default values
func NewPayment() *Payment {
	return &Payment{
		Object: "payment",
	}
}
