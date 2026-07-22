package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func Connect(uri, dbName string, log *zap.Logger) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	// Ping to confirm the connection is actually alive, not just "configured"
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, err
	}

	log.Info("connected to MongoDB", zap.String("database", dbName))

	db := client.Database(dbName)

	if err := ensureIndexes(ctx, db); err != nil {
		return nil, err
	}

	return db, nil
}

func ensureIndexes(ctx context.Context, db *mongo.Database) error {
	usersColl := db.Collection("users")

	indexModel := mongo.IndexModel{
		Keys:    map[string]int{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := usersColl.Indexes().CreateOne(ctx, indexModel)
	return err
}