package game

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	gameApp "intern2023/app"
	"intern2023/database"
	pb "intern2023/pb"
	shareFunc "intern2023/share"
)

// Declare 10 game

type GameItem struct {
	ID         int    `json:"_id"`
	Game       string `json:"game"`
	GuessLimit int    `json:"guessLimit"`
}

func AllGamePatterns() string {
	return "game:*"
}

func GamePattern(IdGame string) string {
	return "game:" + IdGame
}

func CheckExistGame(client *redis.Client, IdGame int) bool {
	IdGameString := strconv.Itoa(IdGame)
	valGame, _ := client.Get(context.Background(), GamePattern(IdGameString)).Result()
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
		shareFunc.CreateArrRand(&randoms)
		string := gameApp.ConvertArrString(&randoms)
		res = append(res, string)
	}
	return res
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
		randId := shareFunc.CreateRandomNumber(min, max)
		item := GameItem{ID: randId, Game: v}
		ui = append(ui, item)

	}
	_, err := gameCollection.InsertMany(context.Background(), ui)
	if err != nil {
		panic(err)
	}
}

func CreateGames(client *redis.Client, sizeGame int, guessLimit int) {
	arr := CreateGameHelper(sizeGame)
	// seed the random number generator
	items := make([]GameItem, len(arr))
	for i, v := range arr {
		// generate a random 8-digit number
		min := 10000000
		max := 99999999
		randId := shareFunc.CreateRandomNumber(min, max)
		items[i] = GameItem{ID: randId, Game: v, GuessLimit: guessLimit}
		val, _ := json.Marshal(items[i])
		_, err := database.Set(client, "game:"+strconv.Itoa(randId), val, 24*7*time.Hour) //
		if err != nil {
			panic(err)
		}
	}
}

// to get data of game
func GetGameValue(client *redis.Client, IdGame int) GameItem {
	// var Game *pb.Game
	var Game GameItem
	getGameString, _ := database.Get(client, "game:"+strconv.Itoa(IdGame))
	_ = json.Unmarshal([]byte(getGameString), &Game)
	return Game
}

// to get list of games
func GetListGame(client *redis.Client) (int, []*pb.Game) {
	keys, _ := client.Keys(context.Background(), "game:*").Result()

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
func UpdateGame(client *redis.Client, wrongLimit int) {
	gameKeys, _ := database.Keys(client, "game:*")
	for _, gameKey := range gameKeys {
		IdGame := shareFunc.GetKeyElement(gameKey, 1)
		keyPattern := "session:*:" + IdGame
		sessionKeys, _ := database.Keys(client, keyPattern)
		if len(sessionKeys) == 0 {
			database.Del(client, gameKey)
			CreateGames(client, 1, wrongLimit)
		}
	}
}

func DeleteGames(client *redis.Client) (int, []*pb.Game) {
	// pipe := client.Pipeline()
	keys, _ := client.Keys(context.Background(), "game:*").Result() //

	var Games []*pb.Game
	for _, key := range keys {
		_, _ = client.Del(context.Background(), key).Result()
	}
	return len(Games), Games
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
			create := shareFunc.CreateRandomNumber(0, 4) // random position
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
