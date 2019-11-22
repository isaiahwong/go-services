package store

import (
	"context"
	"time"

	"github.com/isaiahwong/go-services/src/payment/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoStore struct wrapper
type MongoStore struct {
	client *mongo.Client
	opts   *MongoOptions
}

type idError struct {
	s string
}

func (e *idError) Error() string {
	return e.s
}

// MongoOptions a set of mongo options
type MongoOptions struct {
	connstr  string
	database string
}

var defaultOptions = MongoOptions{
	connstr:  "mongodb://localhost:27017/",
	database: "payment",
}

// A MongoOption sets options such as hostPort; parameters, etc.
type MongoOption func(*MongoOptions)

// ConnectionString returns MongoOption; sets
func ConnectionString(connstr string) MongoOption {
	return func(o *MongoOptions) {
		o.connstr = connstr
	}
}

// Database returns MongoOption; sets default database
func Database(db string) MongoOption {
	return func(o *MongoOptions) {
		o.database = db
	}
}

// NewMongoStore provides a new MongoStore
// TODO add opts param to mongo client
func NewMongoStore(opt ...MongoOption) (*MongoStore, error) {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	// Create a new mongo client
	c, err := mongo.NewClient(options.Client().ApplyURI(opts.connstr))
	if err != nil {
		return nil, err
	}
	return &MongoStore{
		client: c,
		opts:   &opts,
	}, nil
}

// Connect connects to mongodb
func (m *MongoStore) Connect(ctx context.Context, t time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := m.client.Connect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Create creates a new payment object
func (m *MongoStore) Create(ctx context.Context, payment *model.Payment) (*primitive.ObjectID, error) {
	coll := m.client.Database(m.opts.database).Collection("payment")
	payment.ID = primitive.NewObjectID()
	payment.Updated = time.Now()
	payment.Created = time.Now()

	res, err := coll.InsertOne(ctx, payment)
	if err != nil {
		return nil, err
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, &idError{"Invalid OID"}
	}
	return &oid, nil
}
