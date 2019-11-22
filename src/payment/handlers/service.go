package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/isaiahwong/go-services/src/payment/store"
	"github.com/sirupsen/logrus"
)

// PaymentService is the server API for PaymentService service.
type PaymentService struct {
	production bool
	logger     *logrus.Logger
	store      *store.MongoStore
	v          *validator.Validate
}

// NewPaymentService creates a new PaymentService
func NewPaymentService(production bool, logger *logrus.Logger, store *store.MongoStore) *PaymentService {
	if logger == nil {
		panic("logger is nil. Please define a type log.Logger")
	}
	if store == nil {
		panic("store is nil. Please define a type log *store.MongoStore")
	}
	return &PaymentService{
		production: production,
		logger:     logger,
		store:      store,
		v:          validator.New(),
	}

}
