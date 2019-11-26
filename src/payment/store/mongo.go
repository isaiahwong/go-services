package store

import (
	"context"
	"fmt"
	"time"

	"github.com/isaiahwong/go-services/src/payment/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoStore struct wrapper
type MongoStore struct {
	client *mongo.Client
	opts   *mongoOptions
}

// mongoOptions a set of mongo options
type mongoOptions struct {
	connstr  string
	database string
	auth     *options.Credential
	timeout  time.Duration
}

var defaultOptions = mongoOptions{
	connstr:  "mongodb://localhost:27017/",
	database: "payment",
}

// taken from go.mongodb.org/mongo-driver/mongo/options
// MongoCredential holds auth options.
//
// AuthMechanism indicates the mechanism to use for authentication.
// Supported values include "SCRAM-SHA-256", "SCRAM-SHA-1", "MONGODB-CR", "PLAIN", "GSSAPI", and "MONGODB-X509".
//
// AuthMechanismProperties specifies additional configuration options which may be used by certain
// authentication mechanisms. Supported properties are:
// SERVICE_NAME: Specifies the name of the service. Defaults to mongodb.
// CANONICALIZE_HOST_NAME: If true, tells the driver to canonicalize the given hostname. Defaults to false. This
// property may not be used on Linux and Darwin systems and may not be used at the same time as SERVICE_HOST.
// SERVICE_REALM: Specifies the realm of the service.
// SERVICE_HOST: Specifies a hostname for GSSAPI authentication if it is different from the server's address. For
// authentication mechanisms besides GSSAPI, this property is ignored.
//
// AuthSource specifies the database to authenticate against.
//
// Username specifies the username that will be authenticated.
//
// Password specifies the password used for authentication.
//
// PasswordSet specifies if the password is actually set, since an empty password is a valid password.
type MongoCredential struct {
	AuthMechanism           string
	AuthMechanismProperties map[string]string
	AuthSource              string
	Username                string
	Password                string
	PasswordSet             bool
}

// A MongoOption sets options such as hostPort; parameters, etc.
type MongoOption func(*mongoOptions, *options.ClientOptions)

// ConnectionString returns MongoOption; sets
func ConnectionString(connstr string) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		o.connstr = connstr
		m.ApplyURI(connstr)
	}
}

// Database returns MongoOption; sets default database
func Database(db string) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		o.database = db
	}
}

// SetTimeout specifies the timeout for an initial connection to a server.
// If a custom Dialer is used, this method won't be set and the user is
// responsible for setting the ConnectTimeout for connections on the dialer
// themselves.
func SetTimeout(t time.Duration) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		o.timeout = t
		m.SetConnectTimeout(1 * time.Second)
	}
}

// SetAuth Authentication for mongodb
func SetAuth(credential MongoCredential) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		m.SetAuth(options.Credential{
			AuthMechanism: credential.AuthMechanism,
			AuthSource:    credential.AuthSource,
			Username:      credential.Username,
			Password:      credential.Password,
			PasswordSet:   credential.PasswordSet,
		})
	}
}

// SetHeartbeat TODO
func SetHeartbeat(t time.Duration) MongoOption {
	return func(o *mongoOptions, m *options.ClientOptions) {
		m.SetHeartbeatInterval(t)
	}
}

// NewMongoStore provides a new MongoStore
func NewMongoStore(opt ...MongoOption) (*MongoStore, error) {
	opts := defaultOptions
	// Get the initial mongo settings
	mongoOpt := options.Client()

	// Apply options
	for _, o := range opt {
		o(&opts, mongoOpt)
	}

	// Create a new mongo client
	c, err := mongo.NewClient(
		mongoOpt,
	)
	if err != nil {
		return nil, err
	}
	return &MongoStore{
		client: c,
		opts:   &opts,
	}, nil
}

// Connect connects to mongodb
func (m *MongoStore) Connect(ctx context.Context) error {
	err := m.client.Connect(ctx)
	if err != nil {
		return &connectError{fmt.Sprintf("Mongo Connection %v", err)}
	}

	pe := m.Ping()
	if pe != nil {
		return pe
	}

	return nil
}

// Disconnect Disconnects Mongo client
func (m *MongoStore) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// Ping verifies that the client can connect to the topology.
func (m *MongoStore) Ping() error {
	ctx, _ := context.WithTimeout(context.Background(), m.opts.timeout)
	// Test connection
	err := m.client.Ping(ctx, readpref.Primary())
	if err != nil {
		return &connectError{fmt.Sprintf("Mongo Connection %v", err)}
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
