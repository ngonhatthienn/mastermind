package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/redis/go-redis/v9"

	config "intern2023/handler/Config"

)


func ConnectRedisDatabase() (*redis.Client, error) {
	config := config.GetConfig()

	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password, 
		DB:       0,                     
	})
	return client, nil
}

func ConnectMongoDBConnection() *mongo.Client {
	config := config.GetConfig()
	url := "mongodb+srv://" + config.Mongodb.User + ":" + config.Mongodb.Password + config.Mongodb.Port
	opts := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}
	return client
}

func ConnectGamesCollection(client *mongo.Client) *mongo.Collection {
	quickstart := client.Database("quickstart")
	game := quickstart.Collection("games")
	return game
}
