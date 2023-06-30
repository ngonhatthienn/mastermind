package session

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"intern2023/database"
	game "intern2023/handler/Game"
	shareFunc "intern2023/share"

	pb "intern2023/pb"

	"github.com/redis/go-redis/v9"
)

type PlayHistory struct {
	UserRequest   string `json:"userRequest"`
	RightNumber   int32  `json:"rightNumber"`
	RightPosition int32  `json:"rightPosition"`
}
type UserSessionItem struct {
	IdUser        int32         `json:"id_user"`
	IdGame        int32         `json:"id_game"`
	GuessHistory  []PlayHistory `json:"guessHistory"`
	GuessTimes    int32         `json:"guessTimes"`
	GuessRequests string        `json:"guessRequests"`
}
type GameSessionItem struct {
	IdGame   int32    `json:"id_game"`
	UserPlay []string `json:"userPlay"`
}

// Can't use userPlay because when the session is expired => it don't automatically delete session
func AllSessionPatterns(IdUser string) string {
	return "session:" + IdUser + "*"
}

func SessionPattern(IdGame string, IdUser string) string {
	return "session:" + IdUser + ":" + IdGame
}

func CreateUserSession(client *redis.Client, IdUser int32) ([]string, int) {
	// Check if user play any game before
	IdUserString := strconv.Itoa(int(IdUser))
	checkAnyKeys, _ := client.Keys(context.Background(), "session:"+IdUserString+"*").Result() //

	for _, key := range checkAnyKeys {
		_, _ = client.Del(context.Background(), key).Result()
	}
	// Create new session
	GameKeys, _ := database.Keys(client, "game:*") //

	randNumber := shareFunc.CreateRandomNumber(0, 9)
	GameKey := GameKeys[randNumber]
	IdGameString := shareFunc.GetKeyElement(GameKey, 1)

	// Get guessLimit
	var game game.GameItem
	gameData, _ := client.Get(context.Background(), GameKey).Result()
	_ = json.Unmarshal([]byte(gameData), &game)
	// Set session
	IdGame, _ := strconv.Atoi(IdGameString)

	SessionKey := "session:" + IdUserString + ":" + IdGameString
	

	timeNow := time.Now().Unix()

	// Create Init
	// database.HSet(client, SessionKey, field, val, "isWin", false,
	// 	"guessLeft", 10, "timeStart", timeNow, "sumRight", 0, "sumPos", 0)
	client.HSet(context.Background(), SessionKey, "isWin", false,
		"guessLeft", game.GuessLimit, "timeStart", timeNow, "sumRight", 0, "sumPos", 0)
	return GameKeys, IdGame
}

func CreateSessionWithId(client *redis.Client, IdUser int32, IdGame int32) {
	// Check if user play any game before
	// pipe := client.Pipeline()

	IdUserString := strconv.Itoa(int(IdUser))
	checkAnyKeys, _ := client.Keys(context.Background(), "session:"+IdUserString+"*").Result() //
	fmt.Print("keys", checkAnyKeys)

	for _, key := range checkAnyKeys {
		_, _ = client.Del(context.Background(), key).Result()
	}
	// Create new session
	IdGameString := strconv.Itoa(int(IdGame))
	// Get guessLimit
	var game game.GameItem
	gameData, _ := client.Get(context.Background(), "game:" +IdGameString).Result()
	_ = json.Unmarshal([]byte(gameData), &game)

	// Set session

	SessionKey := "session:" + IdUserString + ":" + IdGameString


	// val, _ := json.Marshal(playHistory)
	timeNow := time.Now().Unix()

	// Create Init
	// database.HSet(client, SessionKey, field, val, "isWin", false,
	// 	"guessLeft", 10, "timeStart", timeNow, "sumRight", 0, "sumPos", 0)
	client.HSet(context.Background(), SessionKey, "isWin", false,
		"guessLeft", game.GuessLimit, "timeStart", timeNow, "sumRight", 0, "sumPos", 0)
}

// func GetLastSession(client *redis.Client, key string) string {
// 	// get the latest gameplay by using the maximum nth value

// 	lastNthGamePlay := "gameplay_10"
// 	for i := 10; i >= 1; i-- {
// 		field := fmt.Sprintf("gameplay_%d", i)
// 		_, err := database.HGet(client, key, field)
// 		if err == nil {
// 			lastNthGamePlay = field
// 			break
// 		}
// 	}
// 	result, err := database.HGet(client, key, lastNthGamePlay)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(result)
// 	return result
// }

func PushAndGetHistory(client *redis.Client, key string, userRequest string, rightNumber int32, rightPosition int32) ([]*pb.ListHistory, bool) {
	var listHistory []*pb.ListHistory
	
	listHistoryString, _ := client.HGet(context.Background(), key, "gameplay").Result()
	
		_ = json.Unmarshal([]byte(listHistoryString), &listHistory)
	playHistory := pb.ListHistory{
		UserRequest:   userRequest,
		RightNumber:   rightNumber,
		RightPosition: rightPosition,
	}
	listHistory = append(listHistory, &playHistory)

	
	data, _ := json.Marshal(listHistory)
	oldSumRightInt, _ := client.HGet(context.Background(), key, "sumRight").Int()
	oldSumPostInt, _ := client.HGet(context.Background(), key, "sumPos").Int()

	newSumRight := int(rightNumber) + oldSumRightInt
	newSumPos := int(rightPosition) + oldSumPostInt
	client.HSet(context.Background(), key, "gameplay",data ,"sumRight", newSumRight, "sumPos", newSumPos)

	return listHistory, true
}

func DecodeSession(result string) PlayHistory {
	var resultGamePlay PlayHistory
	err := json.Unmarshal([]byte(result), &resultGamePlay)
	if err != nil {
		panic(err)
	}
	return resultGamePlay
}
