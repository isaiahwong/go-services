package store

import (
	"context"
	"time"
)

// Store interface
type Store interface {
	Connect(context.Context, time.Duration) error
}
