package session

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"intern2023/share"
	"intern2023/database"

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

func CreateUserSession(client *redis.Client, IdUser int32) ([]string, int) {
	// Check if user play any game before

	IdUserString := strconv.Itoa(int(IdUser))
	checkAnyKeys, _ := database.Keys(client, "session:"+IdUserString+"*")
	fmt.Println(checkAnyKeys)

	if len(checkAnyKeys) > 0 {
		database.Del(client, checkAnyKeys[0])
	}

	// Create new session
	GameKeys, _ := database.Keys(client, "game:*")

	randNumber := shareFunc.CreateRandomNumber(0, 9)
	GameKey := GameKeys[randNumber]
	IdGameString := shareFunc.GetKeyElement(GameKey, 1)

	fmt.Println(IdGameString)

	IdGame, _ := strconv.Atoi(IdGameString)

	SessionKey := "session:" + IdUserString + ":" + IdGameString
	field := "gameplay_0"
	playHistory := PlayHistory{
		UserRequest:   "",
		RightNumber:   0,
		RightPosition: 0,
	}

	val, _ := json.Marshal(playHistory)
	timeNow := time.Now().Unix()

	// Create Init
	database.HSet(client, SessionKey, field, val, "isWin", false,
		"guessLeft", 10, "timeStart", timeNow, "sumRight", 0, "sumPos", 0)

	return GameKeys, IdGame
}

func GetLastSession(client *redis.Client, key string) string {
	// get the latest gameplay by using the maximum nth value

	lastNthGamePlay := "gameplay_10"
	for i := 10; i >= 1; i-- {
		field := fmt.Sprintf("gameplay_%d", i)
		_, err := database.HGet(client, key, field)
		if err == nil {
			lastNthGamePlay = field
			break
		}
	}
	result, err := database.HGet(client, key, lastNthGamePlay)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	return result
}

func PushAndGetHistory(client *redis.Client, key string, userRequest string, rightNumber int32, rightPosition int32) ([]*pb.ListHistory,bool) {
	var lastNthGamePlay string
	// prevLastNthGamePlay := lastNthGamePlay
	var listHistory []*pb.ListHistory
	checkFor := time.Now()

	for i := 1; ; i++ {
		var data *pb.ListHistory
		field := fmt.Sprintf("gameplay_%d", i)
		value, err := database.HGet(client, key, field)

		if err != nil {
			lastNthGamePlay = "gameplay_" + strconv.Itoa(i)
			break
		}
		err = json.Unmarshal([]byte(value), &data)
		listHistory = append(listHistory, data)
	}

	checkFor1 := time.Since(checkFor)
	fmt.Println("Time check for", checkFor1.Milliseconds())

	playHistory := PlayHistory{
		UserRequest:   userRequest,
		RightNumber:   rightNumber,
		RightPosition: rightPosition,
	}
	val, _ := json.Marshal(playHistory)

	oldSumRightString, _ := database.HGet(client, key, "sumRight")
	oldSumPostString, _ := database.HGet(client, key, "sumPos")

	convert := time.Now() //

	oldSumRightInt, _ := strconv.Atoi(oldSumRightString)
	oldSumPostInt, _ := strconv.Atoi(oldSumPostString)

	newSumRight := int(rightNumber) + oldSumRightInt
	newSumPos := int(rightPosition) + oldSumPostInt

	convert1 := time.Since(convert) //
	fmt.Println("Time convert ", convert1.Milliseconds())

	// Push all data

	database.HSet(client, key, lastNthGamePlay, val, "sumRight", newSumRight, "sumPos", newSumPos)


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
