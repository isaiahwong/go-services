package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TransactionStatus status for transaction
type TransactionStatus string

// Types for Transaction status
const (
	Succeeded  TransactionStatus = "succeeded"
	Declined                     = "declined"
	Refunded                     = "refunded"
	Pending                      = "pending"
	Transitory                   = "transitory"
)

// Transaction struct
type Transaction struct {
	ID                  primitive.ObjectID `bson:"id,omitempty" json:"id"`
	Payment             primitive.ObjectID `bson:"payment" json:"payment" validate:"required"`
	Coupon              primitive.ObjectID `bson:"coupon" json:"coupon" `
	Refund              primitive.ObjectID `bson:"refund" json:"refund" `
	Object              string             `bson:"object" json:"object" validate:"eq=transaction"`
	User                string             `bson:"user" json:"user" validate:"required"`
	Email               string             `bson:"email" json:"email" validate:"required,email"`
	Provider            PaymentProvider    `bson:"provider" json:"provider" validate:"required"`
	PaypalOrderID       string             `bson:"paypal_order_id" json:"paypal_order_id"`
	StripePaymentIntent string             `bson:"stripe_payment_intent" json:"stripe_payment_intent"`
	Currency            Currency           `bson:"currency" json:"currency" validate:"required"`
	Items               Items              `bson:"items" json:"items" validate:"required"`
	Total               float64            `bson:"total" json:"total" `
	Status              TransactionStatus  `bson:"status" json:"status" validate:"required"`
	TransitoryExpires   time.Time          `bson:"transitory_expires" json:"transitory_expires"`
	Paid                bool               `bson:"paid" json:"paid" validate:"required"`
	TransactionError    TransactionError   `bson:"transaction_error" json:"transaction_error"`
	IP                  string             `bson:"ip" json:"ip" `
	Updated             time.Time          `bson:"updated" json:"updated" validate:"required"`
	Created             time.Time          `bson:"created" json:"created" validate:"required"`
}

// Items transaction Items
type Items struct {
	ID               string `validate:"required"`
	Description      string
	Metadata         map[string]string
	Data             []Data
	TotalItems       int64 `validate:"required" bson:"total_items"`
	Subtotal         float64
	Shipping         float64
	Tax              float64
	ShippingDiscount float64
	Discount         float64
	Currency         Currency
}

// Data stores transaction details
type Data struct {
	ID          string `validate:"required"`
	Name        string `validate:"required"`
	Description string
	Amount      float64 `validate:"required"`
	Quantity    int64   `validate:"required"`
	Metadata    map[string]string
	Currency    Currency
}

// TransactionError details of transaction errors
type TransactionError struct {
	Error           string
	Type            string
	Message         string
	StripeErrorCode string `bson:"stripe_error_code"`
}
