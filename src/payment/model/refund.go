package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Refund struct
type Refund struct {
	ID            primitive.ObjectID `bson:"id,omitempty" json:"id"`
	Object        string             `bson:"object" json:"object" validate:"eq=refund,required"`
	Transaction   primitive.ObjectID `bson:"transaction" json:"transaction"`
	Amount        float64            `bson:"amount" json:"amount" validate:"required"`
	Currency      Currency           `bson:"currency" json:"currency" validate:"required"`
	StripeRefund  string             `bson:"stripe_refund" json:"stripe_refund"`
	PaypalRefund  string             `bson:"paypal_refund" json:"paypal_refund"`
	Reason        RefundReason       `bson:"reason" json:"reason" validate:"required"`
	Status        RefundStatus       `bson:"status" json:"status" validate:"required"`
	FailureRefund FailureRefund      `bson:"failure_refund" json:"failure_refund"`
	FailureReason string             `bson:"failure_reason" json:"failure_reason"`
	Updated       time.Time          `bson:"updated" json:"updated" validate:"required"`
	Created       time.Time          `bson:"created" json:"created" validate:"required"`
}

type FailureRefund struct {
	Ref      string   `bson:"ref" json:"ref"`
	Currency Currency `bson:"currency" json:"currency" validate:"required"`
	Amount   float64  `bson:"amount" json:"amount" validate:"required"`
	Fee      float64  `bson:"fee" json:"fee"`
	Net      float64  `bson:"net" json:"net" validate:"required"`
}

// RefundReason reason for refund
type RefundReason string

// Reasons for refund
const (
	RequestedByCustomer = "requested_by_customer"
	Fraudulent          = "fraudulent"
	Admin               = "admin"
)

// RefundStatus status of refund
type RefundStatus string

// Types of refund status
const (
	RefundSucceeded = "succeeded"
	RefundDeclined  = "declined"
	RefundFailed    = "failed"
	RefundPending   = "pending"
)
