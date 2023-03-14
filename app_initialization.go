package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

const (
	// ExecutionMode is the environment variable that determines the execution mode of the application.
	// Possible values are "local" and "lambda".
	// if the value is "local", the application will run locally and not in lambda.
	ExecutionMode = "EXECUTION_MODE"
	MongoDBURI    = "MONGODB_URI"
)

type InitializationResponse struct {
	environmentIsLocal bool
	mongoDbClient      *mongo.Client
}

func InitializationHandler() (*InitializationResponse, error) {
	setViperDefaults()
	response := &InitializationResponse{}

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	shouldRunLocally := viper.GetString(ExecutionMode) == "local"
	response.environmentIsLocal = shouldRunLocally

	mongoURI := viper.GetString(MongoDBURI)

	mongoClient, err := getMongoClient(mongoURI)
	if err != nil {
		return nil, err
	}

	response.mongoDbClient = mongoClient
	return response, nil

}

func getMongoClient(uri string) (*mongo.Client, error) {
	log.Println("Connecting to MongoDB: ", uri, "")
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

func setViperDefaults() {
	viper.SetDefault(ExecutionMode, "local")
	viper.AddConfigPath(".") // optionally look for config in the working directory
}

//useinvoicewise
//1ONzze7cQOCGq9HS
