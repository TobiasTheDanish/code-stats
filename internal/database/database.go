package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	GetAllDaily() (*mongo.Cursor, error)
	GetAllWeekly() (*mongo.Cursor, error)
	GetAllMonthly() (*mongo.Cursor, error)
	GetAllYearly() (*mongo.Cursor, error)
}

type service struct {
	db *mongo.Client
}

var (
	connectionString = os.Getenv("MONGO_CONNECTION_STRING")
)

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString))

	if err != nil {
		log.Fatal(err)

	}
	return &service{
		db: client,
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) GetAllDaily() (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	collection := s.db.Database("pulse").Collection("daily")

	return collection.Find(ctx, bson.D{})
}

func (s *service) GetAllWeekly() (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	collection := s.db.Database("pulse").Collection("weekly")

	return collection.Find(ctx, bson.D{})
}

func (s *service) GetAllMonthly() (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	collection := s.db.Database("pulse").Collection("monthly")

	return collection.Find(ctx, bson.D{})
}

func (s *service) GetAllYearly() (*mongo.Cursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	collection := s.db.Database("pulse").Collection("yearly")

	return collection.Find(ctx, bson.D{})
}
