package mongodb

import (
	"context"
	"os"
	"time"

	"go-clean-architecture/pkg/logger"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitMongoDB - initialize mongo
func InitMongoDB() (context.Context, func(), *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_URL")))
	if err != nil {
		logger.Error(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		logger.Error(err)
	}

	// Checking the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logger.Error(err)
	}
	logrus.Println("Database connected")

	return ctx, cancel, client
}
