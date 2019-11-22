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
	Object          string             `bson:"object" json:"object" validate:"eq=payment,required"`
	User            string             `bson:"user" json:"user" validate:"required"`
	Email           string             `bson:"email" json:"email"  validate:"required,email"`
	DefaultProvider PaymentProvider    `bson:"default_provider" json:"default_provider" validate:"required"`
	Stripe          Stripe             `bson:"stripe" json:"stripe"`
	Updated         time.Time          `bson:"updated" json:"updated" validate:"required"`
	Created         time.Time          `bson:"created" json:"created" validate:"required"`
}

// NewPayment returns payment with default values
func NewPayment() *Payment {
	return &Payment{
		Object: "payment",
	}
}
