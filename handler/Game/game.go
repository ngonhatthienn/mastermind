package game

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	

	"github.com/redis/go-redis/v9"

	"intern2023/app"
	"intern2023/share"
	"intern2023/database"
	pb "intern2023/pb"

)

// Declare 10 game

type GameItem struct {
	ID         int     `json:"_id"`
	Game       string  `json:"game"`
	WrongLimit int    `json:"wrongLimit"`
}

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

func CreateGames(client *redis.Client, sizeGame int, wrongLimit int) {
	arr := CreateGameHelper(sizeGame)
	// seed the random number generator

	// print the number
	items := make([]GameItem, len(arr))
	for i, v := range arr {
		// generate a random 8-digit number
		min := 10000000
		max := 99999999
		randId := shareFunc.CreateRandomNumber(min, max)
		items[i] = GameItem{ID: randId, Game: v, WrongLimit: wrongLimit}
		// Add Game -----
		val, _ := json.Marshal(items[i])
		_, err := database.Set(client, "game:"+strconv.Itoa(randId), val, 24*7*time.Hour)
		if err != nil {
			panic(err)
		}
	}
}

func GetGameValue(client *redis.Client, IdGame int) *pb.Game {
	var Game *pb.Game
	getGameString, _ := database.Get(client, "game:"+strconv.Itoa(IdGame))
	_ = json.Unmarshal([]byte(getGameString), &Game)
	return Game
}

func GetListGame(client *redis.Client) (int, []*pb.Game) {

	keys, _ := database.Keys(client, "game:*")


	var Games []*pb.Game
	for _, key := range keys {
		val, _ := database.Get(client, key)
		var data *pb.Game
		err := json.Unmarshal([]byte(val), &data)
		if err != nil {
			// Handle the error here
		}
		Games = append(Games, data)
	}
	return len(Games), Games
}

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

func GenerateHint(result string, types string) string {
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
		fmt.Println(randoms)
		for i := 0; i < len(randoms); i++ {
			resultBytes[randoms[i]] = '*'
		}
		result = string(resultBytes)
	}
	return result
}
