package model

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	gameApp "intern2023/app"
	"intern2023/database"
	shareFunc "intern2023/share"

	game "intern2023/handler/Game"
	leaderboard "intern2023/handler/Leaderboard"
	session "intern2023/handler/Session"
	user "intern2023/handler/User"

	pb "intern2023/pb"
)

type Service struct {
	redisClient *redis.Client
	mongoClient *mongo.Client
}

type UserItem struct {
	ID       int32  `json:"_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewService() *Service {
	redisClient, _ := database.CreateRedisDatabase()
	mongoClient := database.CreateMongoDBConnection()

	return &Service{redisClient: redisClient, mongoClient: mongoClient}
}

// GAME
func (s *Service) InitGame() {
	game.CacheGame(s.mongoClient, s.redisClient, 10)
}

func (s *Service) CreateGame(sizeGame int, GuessLimit int) {
	game.CreateGames(s.redisClient, sizeGame, GuessLimit)
}

func (s *Service) ListGame() (int, []*pb.Game) {
	length, Games := game.GetListGame(s.redisClient)
	if length == 0 {
		game.CreateGames(s.redisClient, 10, 30)
		length, Games = game.GetListGame(s.redisClient)
	}
	return length, Games
}

func (s *Service) GetCurrent(IdUser int) (shareFunc.Status, *pb.GameReply) {
	// Check exist user
	// IdUser := int(in.IdUser)
	if user.CheckExistUser(s.redisClient, IdUser) == false {
		status := shareFunc.GenerateStatus(404, "User")
		return status, &pb.GameReply{}
	}
	_, IdGame := session.CreateUserSession(s.redisClient, int32(IdUser))
	Game := game.GetGameValue(s.redisClient, IdGame)
	GameReply := pb.GameReply{XId: int32(Game.ID), GuessLimit: int32(Game.GuessLimit)}

	status := shareFunc.GenerateStatus(200, "Get current")
	return status, &GameReply
}

// Pick one game
func (s *Service) PickGame(IdUser int, IdGame int) (shareFunc.Status, *pb.GameReply) {
	// Check exist user
	if user.CheckExistUser(s.redisClient, IdUser) == false {
		status := shareFunc.GenerateStatus(404, "User")
		return status, &pb.GameReply{}
	}
	// Check exist game
	if game.CheckExistGame(s.redisClient, IdGame) == false {
		status := shareFunc.GenerateStatus(404, "Game")
		return status, &pb.GameReply{}
	}
	// Handle pick game
	session.CreateSessionWithId(s.redisClient, int32(IdUser), int32(IdGame))
	Game := game.GetGameValue(s.redisClient, IdGame)
	GameReply := pb.GameReply{XId: int32(Game.ID), GuessLimit: int32(Game.GuessLimit)}
	status := shareFunc.GenerateStatus(200, "Pick game")
	return status, &GameReply
}

// Update Game
func (s *Service) UpdateGame(GuessLimit int) shareFunc.Status {
	// Check exist user
	game.UpdateGame(s.redisClient, GuessLimit)
	status := shareFunc.GenerateStatus(200, "Update Game")
	return status
}

// Play Game
func (s *Service) PlayGame(IdUser int, UserGuess string) (shareFunc.Status, int, []*pb.ListHistory) {
	if user.CheckExistUser(s.redisClient, IdUser) == false {
		status := shareFunc.GenerateStatus(404, "User")
		return status, 0, nil
	}
	IdUserString := strconv.Itoa(IdUser)

	// Check Exist session or not
	getSessions := "session:" + IdUserString + ":*"
	keySessions, _ := s.redisClient.Keys(context.Background(), getSessions).Result()
	var keySession string
	if len(keySessions) == 0 {
		_, IdGame := session.CreateUserSession(s.redisClient, int32(IdUser))
		keySession = "session:" + strconv.Itoa(IdUser) + ":" + strconv.Itoa(IdGame)
	} else {
		keySession = keySessions[0]
	}

	// Update history when play

	isWin, _ := s.redisClient.HGet(context.Background(), keySession, "isWin").Bool()
	if isWin {
		status := shareFunc.GenerateStatus(200, "")
		status.Message = "You'd already won, please get another game"
		return status, 0, nil
	}
	guessLeft, _ := s.redisClient.HGet(context.Background(), keySession, "guessLeft").Int()

	if guessLeft == 0 {
		status := shareFunc.GenerateStatus(200, "")
		status.Message = "You lose!!!"
		return status, 0, nil
	} else {
		guessLeft--
		s.redisClient.HSet(context.Background(), keySession, "guessLeft", guessLeft)
	}

	IdGameString := shareFunc.GetKeyElement(keySession, 2)
	IdGame, _ := strconv.Atoi(IdGameString)
	getGameValue := game.GetGameValue(s.redisClient, IdGame)

	rightNumber, rightPosition := gameApp.OutputGame(UserGuess, getGameValue.Game)

	// If user win this game
	if rightNumber == rightPosition && rightNumber == 5 {
		s.redisClient.HSet(context.Background(), keySession, "isWin", true)

		// Handle Time

		timeStart, _ := s.redisClient.HGet(context.Background(), keySession, "timeStart").Int64()
		savedTime := time.Unix(timeStart, 0)

		diffInSeconds := 5000 - time.Now().Sub(savedTime).Seconds()

		// Get right and pos
		right, _ := s.redisClient.HGet(context.Background(), keySession, "sumRight").Int()
		pos, _ := s.redisClient.HGet(context.Background(), keySession, "sumPos").Int()

		score := int(diffInSeconds) + guessLeft*100 + (right+pos)*2
		_ = leaderboard.AddScore(s.redisClient, IdUserString, IdGameString, int64(score))
		status := shareFunc.GenerateStatus(200, "")
		status.Message = "You win!!!"
		return status, guessLeft, nil
	}

	var listHistory []*pb.ListHistory
	listHistory, _ = session.PushAndGetHistory(s.redisClient, keySession, UserGuess, int32(rightNumber), int32(rightPosition))
	status := shareFunc.GenerateStatus(200, "")
	status.Message = "Try your best !!!"
	return status, guessLeft, listHistory
}

// Hint Game
func (s *Service) HintGame(IdUser int, Type string) (shareFunc.Status, string) {
	// IdUser := int(in.IdUser)
	if user.CheckExistUser(s.redisClient, IdUser) == false {
		status := shareFunc.GenerateStatus(404, "User")
		return status, ""
	}

	keySessionPattern := "session:" + strconv.Itoa(IdUser) + ":*"

	keySessions, _ := s.redisClient.Keys(context.Background(), keySessionPattern).Result()
	// The session of this user is not exists or expired
	if len(keySessions) == 0 {
		status := shareFunc.GenerateStatus(404, "Session")
		return status, ""
	}

	IdGame := shareFunc.GetKeyElement(keySessions[0], 2)
	key := "game:" + IdGame
	val, _ := s.redisClient.Get(context.Background(), key).Result()

	// The game in this session is not exists or expired
	if val == "" {
		status := shareFunc.GenerateStatus(404, "Game")
		return status, ""
	}
	var Result game.GameItem

	_ = json.Unmarshal([]byte(val), &Result)

	res, isSuccess := game.GenerateHint(Result.Game, Type)
	if !isSuccess {
		status := shareFunc.GenerateStatus(400, "")
		return status, ""
	}
	status := shareFunc.GenerateStatus(400, "Get hint")
	return status, res
}

// USER
func (s *Service) CreateUser(Name string, Password string) (int32, error) {
	Id := user.CreateUser(s.redisClient, Name, Password) // Not in best practices
	return Id, nil
}

func (s *Service) ListUsers() ([]*pb.User, error) {
	// redisClient, _ := database.CreateRedisDatabase()

	keys, _ := s.redisClient.Keys(context.Background(), "user:*").Result()

	cmdS, _ := s.redisClient.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.Get(context.Background(), key).Result()
		}
		return nil
	})

	var Users []*pb.User
	for _, cmd := range cmdS {
		val := cmd.(*redis.StringCmd).Val()
		var data *pb.User
		_ = json.Unmarshal([]byte(val), &data)
		Users = append(Users, data)
	}
	// convert data into proto in controller
	return Users, nil
}

// LEADERBOARD
func (s *Service) GetLeaderBoard(IdGame int, IdUser int, Size int) (shareFunc.Status,[]*pb.LeaderBoardRank,int32, string) {
	// Check exist game
	// IdGame := int(in.IdGame)
	if game.CheckExistGame(s.redisClient, IdGame) == false {
		status := shareFunc.GenerateStatus(404, "Get LeaderBoard")
		return status, nil,0,""

		// return &pb.LeaderBoardReply{Code: 404}, nil
	}
	// Check exits user
	// IdUser := int(in.IdUser)
	if user.CheckExistUser(s.redisClient, IdUser) == false {
		status := shareFunc.GenerateStatus(404, "User")
		return status, nil,0,""
	}

	IdUserString := strconv.Itoa(IdUser)
	leaderboardData, _ := leaderboard.GetLeaderboard(s.redisClient, strconv.Itoa(IdGame), int64(Size), IdUserString)
	status := shareFunc.GenerateStatus(200, "Get LeaderBoard")
	UserRank, UserScore := leaderboard.GetUserRank(s.redisClient, strconv.Itoa(IdGame), strconv.Itoa(IdUser))
	return status, leaderboardData,UserRank,UserScore

}
