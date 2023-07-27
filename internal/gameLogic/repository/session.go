package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"intern2023/database"
	pb "intern2023/pb/game"
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

type SessionRepositoryImpl struct {
	redisClient *redis.Client
	mongoClient *mongo.Client
}

func NewSessionRepositoryImpl() *SessionRepositoryImpl {
	fmt.Println("Create new session repository")
	redisClient, _ := database.ConnectRedisDatabase()
	mongoClient := database.ConnectMongoDBConnection()

	return &SessionRepositoryImpl{
		redisClient: redisClient,
		mongoClient: mongoClient,
	}
}

func (r *SessionRepositoryImpl) GetKeySessionByUserID(IdUser int) (string, bool) {
	IdUserString := strconv.Itoa(IdUser)
	keySessions, _ := r.redisClient.Keys(context.Background(), share.AllSessionPatterns(IdUserString)).Result()
	var keySession string
	if len(keySessions) == 0 {
		return "", false
	} else {
		keySession = keySessions[0]
	}
	return keySession, true
}

func (r *SessionRepositoryImpl) GetSessionValue(hashKey string, fieldKey string) *redis.StringCmd {
	value := r.redisClient.HGet(context.Background(), hashKey, fieldKey)
	return value
}

func (r *SessionRepositoryImpl) SetSessionValue(hashKey string, values ...interface{}) {
	r.redisClient.HSet(context.Background(), hashKey, values)
}

func (r *SessionRepositoryImpl) SetSessionID(hashKey string, IdSession string) {
	_, _ = r.redisClient.Set(context.Background(), hashKey, IdSession, 0).Result()
}

func (r *SessionRepositoryImpl) CreateNewSession(IdUser int32) ([]string, int) {
	// Check if user play any game before
	IdUserString := strconv.Itoa(int(IdUser))
	checkAnyKeys, _ := r.redisClient.Keys(context.Background(), share.AllSessionPatterns(IdUserString)).Result() //

	for _, key := range checkAnyKeys {
		_, _ = r.redisClient.Del(context.Background(), key).Result()
	}
	// Create new session
	GameKeys, _ := r.redisClient.Keys(context.Background(), share.AllGamePattern()).Result() //

	randNumber := share.CreateRandomNumber(0, 9)
	GameKey := GameKeys[randNumber]
	IdGameString := share.GetKeyElement(GameKey, 1)

	// Get guessLimit
	var game Game
	gameData, _ := r.redisClient.Get(context.Background(), GameKey).Result()
	_ = json.Unmarshal([]byte(gameData), &game)
	// Set session
	IdGame, _ := strconv.Atoi(IdGameString)

	SessionKey := "session:" + IdUserString + ":" + IdGameString

	timeNow := time.Now().Unix()

	// Create Init
	r.redisClient.HSet(context.Background(), SessionKey, "isWin", false,
		"guessLeft", game.GuessLimit, "timeStart", timeNow, "sumRight", 0, "sumPos", 0)
	return GameKeys, IdGame
}

func (r *SessionRepositoryImpl) CreateSessionWithId(IdUser int32, IdGame int32) {
	// Check if user play any game before
	IdUserString := strconv.Itoa(int(IdUser))
	checkAnyKeys, _ := r.redisClient.Keys(context.Background(), share.AllSessionPatterns(IdUserString)).Result() //

	for _, key := range checkAnyKeys {
		_, _ = r.redisClient.Del(context.Background(), key).Result()
	}
	var game Game
	IdGameString := strconv.Itoa(int(IdGame))
	gameData, _ := r.redisClient.Get(context.Background(), share.GamePattern(IdGameString)).Result()
	_ = json.Unmarshal([]byte(gameData), &game)

	// Set session
	SessionKey := share.SessionPattern(IdGameString, IdUserString)

	timeNow := time.Now().Unix()

	// Create Init
	r.redisClient.HSet(context.Background(), SessionKey, "isWin", false,
		"guessLeft", game.GuessLimit, "timeStart", timeNow, "sumRight", 0, "sumPos", 0)
}

func (r *SessionRepositoryImpl) PushAndGetHistory(key string, userRequest string, rightNumber int32, rightPosition int32) ([]*pb.ListHistory, bool) {
	var listHistory []*pb.ListHistory
	pipe := r.redisClient.Pipeline()
	pipe.HGet(context.Background(), key, "gameplay")
	pipe.HGet(context.Background(), key, "sumRight")
	pipe.HGet(context.Background(), key, "sumPos")
	result, _ := pipe.Exec(context.Background())

	listHistoryString,_ := result[0].(*redis.StringCmd).Result()
	oldSumRightInt, _ := result[1].(*redis.StringCmd).Int()
	oldSumPostInt, _ := result[2].(*redis.StringCmd).Int()

	_ = json.Unmarshal([]byte(listHistoryString), &listHistory)
	playHistory := pb.ListHistory{
		UserRequest:   userRequest,
		RightNumber:   rightNumber,
		RightPosition: rightPosition,
	}
	listHistory = append(listHistory, &playHistory)

	data, _ := json.Marshal(listHistory)

	newSumRight := int(rightNumber) + oldSumRightInt
	newSumPos := int(rightPosition) + oldSumPostInt
	r.redisClient.HSet(context.Background(), key, "gameplay", data, "sumRight", newSumRight, "sumPos", newSumPos)

	return listHistory, true
}

func (r *SessionRepositoryImpl) DecodeSession(result string) PlayHistory {
	var resultGamePlay PlayHistory
	err := json.Unmarshal([]byte(result), &resultGamePlay)
	if err != nil {
		panic(err)
	}
	return resultGamePlay
}
