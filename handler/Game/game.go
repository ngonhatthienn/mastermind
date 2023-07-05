package game

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	gameApp "intern2023/app"
	"intern2023/database"
	pb "intern2023/pb"
	"intern2023/share"
)

// Declare 10 game

type GameItem struct {
	ID         int    `json:"_id"`
	Game       string `json:"game"`
	GuessLimit int    `json:"guessLimit"`
}
type Game struct {
	ID         int    `bson:"id"`
	Game       string `bson:"game"`
	GuessLimit int    `bson:"guessLimit"`
}



// Check the database for any games or not
func CheckAnyGames(client *redis.Client) bool {
	keyGames, _ := client.Keys(context.Background(), share.AllGamePattern()).Result()
	if len(keyGames) == 0 {
		return false
	}
	return true
}

// Check if the database exists game with a given id or not
func CheckExistGame(client *redis.Client, IdGame int) bool {
	IdGameString := strconv.Itoa(IdGame)
	valGame, _ := client.Get(context.Background(), share.GamePattern(IdGameString)).Result()
	if valGame == "" {
		return false
	}
	return true
}

// Help for Create Game
func CreateGameHelper(sizeGame int) []string {
	res := []string{}
	for i := 1; i <= sizeGame; i++ {
		randoms := [5]int{}
		share.CreateArrRand(&randoms)
		string := gameApp.ConvertArrString(&randoms)
		res = append(res, string)
	}
	return res
}

// Check the database for any games or not, if not, generate 10 games
func CheckAndGenerateGame(mongoClient *mongo.Client, redisClient *redis.Client) {
	if !CheckAnyGames(redisClient) {
		CacheGame(mongoClient, redisClient, 30)
	}
}

// for Create Game
func CreateGamesMongo(client *mongo.Client, sizeGame int) {
	arr := CreateGameHelper(sizeGame)
	gameCollection := database.CreateGamesCollection(client)
	var ui []interface{}
	for _, v := range arr {
		// generate a random 8-digit number
		min := 10000000
		max := 99999999
		randId := share.CreateRandomNumber(min, max)
		item := GameItem{ID: randId, Game: v}
		ui = append(ui, item)

	}
	_, err := gameCollection.InsertMany(context.Background(), ui)
	if err != nil {
		panic(err)
	}
}

func CacheGame(mongoClient *mongo.Client, redisClient *redis.Client, guessLimit int) {
	gameCollection := database.CreateGamesCollection(mongoClient)

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
		var gameItem GameItem
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

	var records []GameItem
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

// func CreateGames(client *redis.Client, sizeGame int, guessLimit int) {
// 	arr := CreateGameHelper(sizeGame)
// 	// seed the random number generator
// 	items := make([]GameItem, len(arr))
// 	for i, v := range arr {
// 		// generate a random 8-digit number
// 		min := 10000000
// 		max := 99999999
// 		randId := share.CreateRandomNumber(min, max)
// 		items[i] = GameItem{ID: randId, Game: v, GuessLimit: guessLimit}
// 		val, _ := json.Marshal(items[i])
// 		_, err := client.Set(context.Background(), share.GamePattern(strconv.Itoa(randId)), val, 24*7*time.Hour).Result() //
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

// to get data of game
func GetGameValue(client *redis.Client, IdGame int) GameItem {
	var Game GameItem
	getGameString, _ := client.Get(context.Background(), share.GamePattern(strconv.Itoa(IdGame))).Result()
	_ = json.Unmarshal([]byte(getGameString), &Game)
	return Game
}

// to get list of games
func GetListGame(client *redis.Client) (int, []*pb.Game) {
	keys, _ := client.Keys(context.Background(), share.AllGamePattern()).Result()

	cmdS, _ := client.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.Get(context.Background(), key).Result()
		}
		return nil
	})

	var Games []*pb.Game
	for _, cmd := range cmdS {
		val := cmd.(*redis.StringCmd).Val()
		var data *pb.Game
		_ = json.Unmarshal([]byte(val), &data)
		Games = append(Games, data)
	}

	return len(Games), Games
}

// to update list of games
func UpdateGame(redisClient *redis.Client, mongoClient *mongo.Client, wrongLimit int) {
	DeleteGames(redisClient)
	CacheGame(mongoClient, redisClient, 30)
}

func DeleteGames(client *redis.Client) int {
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

func haveResults(result []int) []int {
	res := []int{1, 2, 3, 4, 5}

	for i := 0; i < len(result); i++ {
		for i := 1; i < 9; i++ {
			if i == result[i] {
				res[i] = i
			}
		}
	}
	return res
}

func GenerateHint(result string, types string) (string, bool) {
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
