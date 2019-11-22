package model

// Stripe struct
type Stripe struct {
	Customer             string `bson:"customer" json:"customer"`
	DefaultPaymentMethod string `bson:"default_payment_method" json:"default_payment_method"`
}
