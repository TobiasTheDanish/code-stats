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
	Daily() ([]CodingSession, error)
}

type service struct {
	db *mongo.Client
}

// Period represents the time period for which the coding sessions have been aggregated.
type Period int8

const (
	Day Period = iota
	Week
	Month
	Year
)

// CodingSession represents a coding session that has been aggregated
// for a given time period (day, week, month, year).
type CodingSession struct {
	ID           string       `bson:"_id,omitempty"`
	Period       Period       `bson:"period"`
	EpochDateMs  int64        `bson:"date"`
	DateString   string       `bson:"date_string"`
	TotalTimeMs  int64        `bson:"total_time_ms"`
	Repositories []Repository `bson:"repositories"`
}

type Repository struct {
	Name       string `bson:"name"`
	Files      []File `bson:"files"`
	DurationMs int64  `bson:"duration_ms"`
}

type File struct {
	Name       string `bson:"name"`
	Path       string `bson:"path"`
	Filetype   string `bson:"filetype"`
	DurationMs int64  `bson:"duration_ms"`
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

func (s *service) Daily() ([]CodingSession, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	collection := s.db.Database("pulse").Collection("daily")

	data, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	result := make([]CodingSession, 0, 0)
	for data.Next(context.TODO()) {
		var res CodingSession
		err = data.Decode(&res)
		if err != nil {
			return nil, err
		}

		result = append(result, res)
	}

	return result, nil
}
