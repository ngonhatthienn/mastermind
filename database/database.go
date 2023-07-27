package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/golobby/dotenv"
	"github.com/redis/go-redis/v9"
)

type database struct {
	Redis struct {
		Addr     string `env:"REDIS_ADDR"`
		Password string `env:"REDIS_PASSWORD"`
	}
	Mongodb struct {
		Port     string `env:"MONGO_PORT"`
		User     string `env:"MONGO_USER"`
		Password string `env:"MONGO_PASSWORD"`
	}
}

func ConnectRedisDatabase() (*redis.Client, error) {
	config := database{}
	file, err := os.Open("app.env")
	err = dotenv.NewDecoder(file).Decode(&config)
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password, 
		DB:       0,                     
	})
	return client, nil
}

func ConnectMongoDBConnection() *mongo.Client {
	config := database{}
	file, err := os.Open("app.env")
	err = dotenv.NewDecoder(file).Decode(&config)
	if err != nil {
		panic(err)
	}
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
