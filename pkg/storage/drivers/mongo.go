package drivers

import (
	"alnovi/customizer/pkg/storage"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoConfig struct {
	Host     string
	Port     uint
	Database string
	User     string
	Password string
	Timeout  uint8
}

type Mongo struct {
	*mongo.Client
	config   MongoConfig
	database *mongo.Database
	timeout  time.Duration
}

func NewMongo(config MongoConfig) (*Mongo, error) {
	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", config.User, config.Password, config.Host, config.Port, config.Database)

	client, err := mongo.NewClient(options.Client().ApplyURI(dsn))

	if err != nil {
		return nil, err
	}

	timeout := time.Duration(config.Timeout) * time.Second

	return &Mongo{
		Client:   client,
		config:   config,
		database: client.Database(config.Database),
		timeout:  timeout,
	}, nil
}

func (m *Mongo) Health() storage.StatusConnector {
	err := m.Ping(m.newTimeoutCtx(), nil)

	if err != nil {
		return storage.StatusConnector{
			Status: storage.ErrorStatus,
		}
	}

	return storage.StatusConnector{
		Status: storage.WorkingStatus,
	}
}

func (m *Mongo) Connect() error {
	err := m.Client.Connect(m.newTimeoutCtx())

	if err != nil {
		return err
	}

	err = m.Ping(m.newTimeoutCtx(), nil)

	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) CurrentDB() *mongo.Database {
	return m.database
}

func (m *Mongo) Close() error {
	return m.Disconnect(m.newTimeoutCtx())
}

func (m *Mongo) newTimeoutCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), m.timeout)

	return ctx
}
