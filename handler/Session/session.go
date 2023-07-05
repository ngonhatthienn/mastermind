package session

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	game "intern2023/handler/Game"
	pb "intern2023/pb"
	"intern2023/share"
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

func CreateUserSession(client *redis.Client, IdUser int32) ([]string, int) {
	// Check if user play any game before
	IdUserString := strconv.Itoa(int(IdUser))
	checkAnyKeys, _ := client.Keys(context.Background(), share.AllSessionPatterns(IdUserString)).Result() //

	for _, key := range checkAnyKeys {
		_, _ = client.Del(context.Background(), key).Result()
	}
	// Create new session
	GameKeys, _ := client.Keys(context.Background(), share.AllGamePattern()).Result() //

	randNumber := share.CreateRandomNumber(0, 9)
	GameKey := GameKeys[randNumber]
	IdGameString := share.GetKeyElement(GameKey, 1)

	// Get guessLimit
	var game game.GameItem
	gameData, _ := client.Get(context.Background(), GameKey).Result()
	_ = json.Unmarshal([]byte(gameData), &game)
	// Set session
	IdGame, _ := strconv.Atoi(IdGameString)

	SessionKey := "session:" + IdUserString + ":" + IdGameString

	timeNow := time.Now().Unix()

	// Create Init
	client.HSet(context.Background(), SessionKey, "isWin", false,
		"guessLeft", game.GuessLimit, "timeStart", timeNow, "sumRight", 0, "sumPos", 0)
	return GameKeys, IdGame
}

func CreateSessionWithId(client *redis.Client, IdUser int32, IdGame int32) {
	// Check if user play any game before
	IdUserString := strconv.Itoa(int(IdUser))
	checkAnyKeys, _ := client.Keys(context.Background(), share.AllSessionPatterns(IdUserString)).Result() //
	fmt.Print("keys", checkAnyKeys)

	for _, key := range checkAnyKeys {
		_, _ = client.Del(context.Background(), key).Result()
	}
	// Create new session
	// Get guessLimit
	var game game.GameItem
	IdGameString := strconv.Itoa(int(IdGame))
	gameData, _ := client.Get(context.Background(), share.GamePattern(IdGameString)).Result()
	_ = json.Unmarshal([]byte(gameData), &game)

	// Set session

	SessionKey := share.SessionPattern(IdGameString, IdUserString)

	// val, _ := json.Marshal(playHistory)
	timeNow := time.Now().Unix()

	// Create Init
	client.HSet(context.Background(), SessionKey, "isWin", false,
		"guessLeft", game.GuessLimit, "timeStart", timeNow, "sumRight", 0, "sumPos", 0)
}

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
	client.HSet(context.Background(), key, "gameplay", data, "sumRight", newSumRight, "sumPos", newSumPos)

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
