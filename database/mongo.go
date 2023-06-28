package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
func CreateMongoDBConnection() *mongo.Client{
	// uri := "mongodb+srv://admin:Q7wpSDe11gdhIt1y@atlascluster.wp57hdf.mongodb.net/"
	url := "mongodb+srv://admin:Q7wpSDe11gdhIt1y@atlascluster.wp57hdf.mongodb.net/?retryWrites=true&w=majority"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}
	return client
}
func CreateMongoDBColumn() {
	client := CreateMongoDBConnection()
	quickstart := client.Database("quickstart")
	game := quickstart.Collection("game")
	gameResult, err := game.InsertOne(context.Background(), bson.D{
		{Key: "id", Value: 1132231},
		{Key: "game", Value: "98764"},
		{Key: "guessLimit", Value: 30},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(gameResult.InsertedID)
}

func CreateGamesCollection(client *mongo.Client) *mongo.Collection{
	quickstart := client.Database("quickstart")
	game := quickstart.Collection("games")
	return game
}

