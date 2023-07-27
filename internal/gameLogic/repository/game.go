package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"intern2023/database"
	"intern2023/share"
)

type Game struct {
	ID         int    `json:"_id"`
	Game       string `json:"game"`
	GuessLimit int    `json:"guessLimit"`
}

type GameRepositoryImpl struct {
	redisClient *redis.Client
	mongoClient *mongo.Client
}

func NewGameRepositoryImpl() *GameRepositoryImpl {
	redisClient, _ := database.ConnectRedisDatabase()
	mongoClient := database.ConnectMongoDBConnection()

	return &GameRepositoryImpl{
		redisClient: redisClient,
		mongoClient: mongoClient,
	}
}

// Cache game from mongodb to redis
func CacheGameFromDB(mongoClient *mongo.Client, redisClient *redis.Client, guessLimit int) {
	gameCollection := database.ConnectGamesCollection(mongoClient)

	// Update games that already played
	keyPattern := "session:*"
	sessionKeys, _ := redisClient.Keys(context.Background(), keyPattern).Result()
	mp := make(map[string]int)

	for _, sessionKey := range sessionKeys {
		IdGame := share.GetKeyElement(sessionKey, 2)

		if mp[IdGame] == 1 {
			continue
		}
		mp[IdGame] = 1
		IdGameInt, _ := strconv.Atoi(IdGame)
		filter := bson.D{{Key: "id", Value: IdGameInt}}
		var gameItem Game
		gameCollection.FindOne(context.Background(), filter).Decode(&gameItem)
		gameItem.GuessLimit = guessLimit

		gameData, _ := json.Marshal(gameItem)
		_ = redisClient.Set(context.Background(), share.GamePattern(IdGame), gameData, 10*time.Minute)
	}

	gameSize := 10 - len(mp)

	filter := bson.A{
		bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: gameSize}}}},
	}

	cursor, err := gameCollection.Aggregate(context.Background(), filter)
	if err != nil {
		panic(err)
	}
	// Decode the resulting cursor into a slice of Record structs
	var records []Game
	if err := cursor.All(context.Background(), &records); err != nil {
		panic(err)
	}

	// Set the records into redis: can use pipeline
	_, err = redisClient.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, record := range records {
			record.GuessLimit = guessLimit
			gameData, _ := json.Marshal(record)
			_ = pipe.Set(context.Background(), share.GamePattern(strconv.Itoa(record.ID)), gameData, 10*time.Minute)
		}
		return nil
	})
	if err != nil && err != redis.Nil {
		panic(err)
	}
}

// Check the database if any games exist or not
func CheckAnyGames(client *redis.Client) bool {
	keyGames, _ := client.Keys(context.Background(), share.AllGamePattern()).Result()
	if len(keyGames) == 0 {
		return false
	}
	return true
}

// If not exists any game in cache, generate 10 games
func (r *GameRepositoryImpl) CacheGameFromDB() {
	if !CheckAnyGames(r.redisClient) {
		CacheGameFromDB(r.mongoClient, r.redisClient, 30)
	}
}

// Delete All game in cache
func DeleteGameCache(client *redis.Client) int {
	keys, _ := client.Keys(context.Background(), share.AllGamePattern()).Result() //

	_, err := client.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.Del(context.Background(), key)
		}
		return nil
	})
	if err != nil && err != redis.Nil {
		return 1
	}

	return 0
}

// to get list of games
func (r *GameRepositoryImpl) GetListGame() (int, []Game) {
	keys, _ := r.redisClient.Keys(context.Background(), share.AllGamePattern()).Result()

	cmdS, _ := r.redisClient.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.Get(context.Background(), key).Result()
		}
		return nil
	})

	var Games []Game
	for _, cmd := range cmdS {
		val := cmd.(*redis.StringCmd).Val()
		var game Game
		_ = json.Unmarshal([]byte(val), &game)
		Games = append(Games, game)
	}

	return len(Games), Games
}

// to update list of games
func (r *GameRepositoryImpl) UpdateListGame(wrongLimit int) {
	DeleteGameCache(r.redisClient)
	CacheGameFromDB(r.mongoClient, r.redisClient, wrongLimit)
}

// to get data of game
func (r *GameRepositoryImpl) GetGame(IdGame int) (Game, bool) {
	var Result Game
	gameValueString, err := r.redisClient.Get(context.Background(), share.GamePattern(strconv.Itoa(IdGame))).Result()
	if gameValueString == "" || err != nil {
		return Result , false
	}
	_ = json.Unmarshal([]byte(gameValueString), &Result)
	return Result, true
}

// Check if the database exists game with a given id or not
func  (r *GameRepositoryImpl) CheckExistGame(IdGame int) bool {
	IdGameString := strconv.Itoa(IdGame)
	valGame, _ := r.redisClient.Get(context.Background(), share.GamePattern(IdGameString)).Result()
	if valGame == "" {
		return false
	}
	return true
}

func (r *GameRepositoryImpl) GenerateHint(result string, types string) (string, bool) {
	resultBytes := []byte(result)
	switch types {
	case "3begin":
		for i := 3; i < len(result); i++ {
			resultBytes[i] = '*'
		}
		result = string(resultBytes)
	case "3final":
		for i := 0; i < len(result)-3; i++ {
			resultBytes[i] = '*'
		}
		result = string(resultBytes)
	case "3random":
		check := [5]int{}
		randoms := [2]int{}

		for i := 0; i < 2; i++ {
			create := share.CreateRandomNumber(0, 4) // random position
			if check[create] == 0 {
				randoms[i] = create
				check[create]++
			} else {
				i--
			}
		}
		for i := 0; i < len(randoms); i++ {
			resultBytes[randoms[i]] = '*'
		}
		result = string(resultBytes)
	}
	if result == "" {
		return "", false
	}
	return result, true
}
