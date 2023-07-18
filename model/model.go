package model

import (
	"context"
	"encoding/json"
	"strconv"

	"intern2023/database"
	"intern2023/token"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/metadata"

	gameApp "intern2023/app"
	game "intern2023/handler/Game"
	leaderboard "intern2023/handler/Leaderboard"
	session "intern2023/handler/Session"
	user "intern2023/handler/User"
	pb "intern2023/pb"
	share "intern2023/share"
)

type Service struct {
	redisClient *redis.Client
	mongoClient *mongo.Client
	pasetoMaker token.PasetoMaker
}

func NewService() *Service {
	redisClient, _ := database.CreateRedisDatabase()
	mongoClient := database.CreateMongoDBConnection()
	pasetoMaker, _ := token.NewPasetoMaker()

	return &Service{redisClient: redisClient, mongoClient: mongoClient, pasetoMaker: pasetoMaker}
}

// GAME
func (s *Service) CreateGame(sizeGame int, GuessLimit int) {
	game.UpdateGame(s.mongoClient, s.redisClient, GuessLimit)
}

func (s *Service) ListGame() (int, []game.Game) {
	// Check Any Games, if not, generate it
	game.CheckAndGenerateGame(s.mongoClient, s.redisClient)
	// Get list game
	length, Games := game.GetListGame(s.redisClient)
	return length, Games
}

func (s *Service) GetCurrent(IdUser int) (share.Status, *pb.GameReply) {
	// Check Any Games, if not, generate it
	game.CheckAndGenerateGame(s.mongoClient, s.redisClient)
	// Check exist user
	if user.CheckExistUser(s.redisClient, IdUser) == false {
		status := share.GenerateStatus(404, "User")
		return status, &pb.GameReply{}
	}
	_, IdGame := session.CreateUserSession(s.redisClient, int32(IdUser))
	Game := game.GetGameValue(s.redisClient, IdGame)
	GameReply := pb.GameReply{XId: int32(Game.ID), GuessLimit: int32(Game.GuessLimit)}

	status := share.GenerateStatus(200, "Get current")
	return status, &GameReply
}

// Pick one game
func (s *Service) PickGame(IdUser int, IdGame int) (share.Status, *pb.GameReply) {
	// Check Any Games, if not, generate it
	game.CheckAndGenerateGame(s.mongoClient, s.redisClient)
	// Check exist user
	if user.CheckExistUser(s.redisClient, IdUser) == false {
		status := share.GenerateStatus(404, "User")
		return status, &pb.GameReply{}
	}
	// Check exist game
	if game.CheckExistGame(s.redisClient, IdGame) == false {
		status := share.GenerateStatus(404, "Game")
		return status, &pb.GameReply{}
	}
	// Handle pick game
	session.CreateSessionWithId(s.redisClient, int32(IdUser), int32(IdGame))
	Game := game.GetGameValue(s.redisClient, IdGame)
	GameReply := pb.GameReply{XId: int32(Game.ID), GuessLimit: int32(Game.GuessLimit)}
	status := share.GenerateStatus(200, "Pick game")
	return status, &GameReply
}

// Update Game
func (s *Service) UpdateGame(GuessLimit int) share.Status {
	game.UpdateGame(s.mongoClient, s.redisClient, GuessLimit)
	status := share.GenerateStatus(200, "Update Game")
	return status
}

// Play Game
func (s *Service) PlayGame(IdUser int, UserGuess string) (share.Status, int, []*pb.ListHistory) {
	// Check Any Games, if not, generate it
	game.CheckAndGenerateGame(s.mongoClient, s.redisClient)
	// Check exist user
	if user.CheckExistUser(s.redisClient, IdUser) == false {
		status := share.GenerateStatus(404, "User")
		return status, 0, nil
	}
	IdUserString := strconv.Itoa(IdUser)

	// Check Exist session or not
	keySessions, _ := s.redisClient.Keys(context.Background(), share.AllSessionPatterns(IdUserString)).Result()
	var keySession string
	if len(keySessions) == 0 {
		_, IdGame := session.CreateUserSession(s.redisClient, int32(IdUser))
		keySession = share.SessionPattern(strconv.Itoa(IdGame), strconv.Itoa(IdUser))
	} else {
		keySession = keySessions[0]
	}
	// Check if user already win or not
	isWin, _ := s.redisClient.HGet(context.Background(), keySession, "isWin").Bool()
	if isWin {
		status := share.GenerateStatus(200, "")
		status.Message = "You'd already won, please get another game"
		return status, 0, nil
	}
	guessLeft, _ := s.redisClient.HGet(context.Background(), keySession, "guessLeft").Int()

	if guessLeft == 0 {
		status := share.GenerateStatus(200, "")
		status.Message = "You lose!!!"
		return status, 0, nil
	} else {
		guessLeft--
		s.redisClient.HSet(context.Background(), keySession, "guessLeft", guessLeft)
	}

	IdGameString := share.GetKeyElement(keySession, 2)
	IdGame, _ := strconv.Atoi(IdGameString)
	getGameValue := game.GetGameValue(s.redisClient, IdGame)

	rightNumber, rightPosition := gameApp.OutputGame(UserGuess, getGameValue.Game)
	// If user win this game
	if rightNumber == rightPosition && rightNumber == 5 {
		s.redisClient.HSet(context.Background(), keySession, "isWin", true)
		// Handle Time
		score := share.CalcScore(s.redisClient, keySession, guessLeft)
		_ = leaderboard.AddScore(s.redisClient, IdUserString, IdGameString, int64(score))
		status := share.GenerateStatus(200, "")
		status.Message = "You win!!!"
		return status, guessLeft, nil
	}

	var listHistory []*pb.ListHistory
	listHistory, _ = session.PushAndGetHistory(s.redisClient, keySession, UserGuess, int32(rightNumber), int32(rightPosition))
	status := share.GenerateStatus(200, "")
	status.Message = "Try your best !!!"
	return status, guessLeft, listHistory
}

// Hint Game
func (s *Service) HintGame(IdUser int, Type string) (share.Status, string) {
	// Check Any Games, if not, generate it
	game.CheckAndGenerateGame(s.mongoClient, s.redisClient)
	if user.CheckExistUser(s.redisClient, IdUser) == false {
		status := share.GenerateStatus(404, "User")
		return status, ""
	}

	keySessionPattern := share.AllSessionPatterns(strconv.Itoa(IdUser))

	keySessions, _ := s.redisClient.Keys(context.Background(), keySessionPattern).Result()
	// The session of this user is not exists or expired
	if len(keySessions) == 0 {
		status := share.GenerateStatus(404, "Session")
		return status, ""
	}

	IdGame := share.GetKeyElement(keySessions[0], 2)
	key := "game:" + IdGame
	val, _ := s.redisClient.Get(context.Background(), key).Result()

	// The game in this session is not exists or expired
	if val == "" {
		status := share.GenerateStatus(404, "Game")
		return status, ""
	}
	var Result game.Game

	_ = json.Unmarshal([]byte(val), &Result)

	res, isSuccess := game.GenerateHint(Result.Game, Type)
	if !isSuccess {
		status := share.GenerateStatus(400, "Get hint")
		return status, ""
	}
	status := share.GenerateStatus(200, "")
	return status, res
}

// USER
func (s *Service) LogIn(Name string, Password string) (share.Status, int, string, bool) {
	IdUser, userRole, ok := user.LogIn(s.redisClient, Name, Password)
	if ok {
		status := share.GenerateStatus(200, "LogIn")
		return status, IdUser, userRole, ok
	}
	status := share.GenerateStatus(404, "User")
	return status, IdUser, userRole, ok
}

func (s *Service) CreateToken(IdUser int, userRole string) string {
	IdUserString := strconv.Itoa(IdUser)
	token, IdSession := s.pasetoMaker.CreateToken(IdUserString, userRole)
	_, _ = s.redisClient.Set(context.Background(), share.UserPatternSession(IdUserString), IdSession, 0).Result()
	// IdSession
	return token
}

func (s *Service) CreateUser(Fullname string, Username string, Password string, Email string, Role string) (int32, error) {
	Id := user.CreateUser(s.redisClient, Fullname, Username, Password, Email, Role) // Not in best practices
	return Id, nil
}

func (s *Service) ListUsers() ([]user.User, error) {
	keys, _ := s.redisClient.Keys(context.Background(), share.AllUserPattern()).Result()

	cmdS, _ := s.redisClient.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.Get(context.Background(), key).Result()
		}
		return nil
	})

	var Users []user.User
	for _, cmd := range cmdS {
		val := cmd.(*redis.StringCmd).Val()
		var user user.User
		_ = json.Unmarshal([]byte(val), &user)
		Users = append(Users, user)
	}
	// convert data into proto in controller
	return Users, nil
}

// LEADERBOARD
func (s *Service) GetLeaderBoard(IdGame int, IdUser int, Size int, isAdmin bool) (share.Status, []leaderboard.LeaderBoard, int32, string) {
	// Check Any Games, if not, generate it
	game.CheckAndGenerateGame(s.mongoClient, s.redisClient)
	// Check exist game
	if game.CheckExistGame(s.redisClient, IdGame) == false {
		status := share.GenerateStatus(404, "Id Game")
		return status, nil, 0, ""
	}
	// Check exits user
	if user.CheckExistUser(s.redisClient, IdUser) == false {
		status := share.GenerateStatus(404, "User")
		return status, nil, 0, ""
	}

	IdUserString := strconv.Itoa(IdUser)
	var UserRank int32
	var UserScore string
	leaderboardData, err := leaderboard.GetLeaderboard(s.redisClient, strconv.Itoa(IdGame), int64(Size), IdUserString)
	if err != nil || leaderboardData == nil {
		status := share.GenerateStatus(200, "")
		status.Message = "No user has won this game yet"
		return status, leaderboardData, UserRank, UserScore
	}
	status := share.GenerateStatus(200, "Get LeaderBoard")
	if !isAdmin {
		UserRank, UserScore = leaderboard.GetUserRank(s.redisClient, strconv.Itoa(IdGame), strconv.Itoa(IdUser))
	}

	return status, leaderboardData, UserRank, UserScore
}

// AUTHORIZATION
func (s *Service) AuthorAndAuthn(md metadata.MD, permission string) (share.Status, int) {
	bearerToken := md.Get("authorization")
	if len(bearerToken) <= 0 {
		status := share.GenerateStatus(401, "")
		return status, 0
	}
	reqToken := share.GetTokenElement(bearerToken[0], 1)
	decryptedToken, decryptedOk := s.pasetoMaker.DecryptedToken(reqToken)
	if !decryptedOk {
		status := share.GenerateStatus(401, "Token")
		status.Message = "Invalid or Expired token "
		return status, 0
	}

	IdUserString, ok := s.pasetoMaker.CheckExistUser(decryptedToken, s.redisClient)
	if !ok {
		status := share.GenerateStatus(404, "User")
		return status, 0
	}
	_, ok = s.pasetoMaker.CheckExactSession(decryptedToken, s.redisClient)
	if !ok {
		status := share.GenerateStatus(401, "Token")
		status.Message = "Expired token "
		return status, 0
	}

	isAuthn := s.pasetoMaker.Authentication(decryptedToken, permission)
	if !isAuthn {
		status := share.GenerateStatus(403, "")
		return status, 0
	}

	IdUser, _ := strconv.Atoi(IdUserString)
	status := share.GenerateStatus(200, "")
	return status, IdUser
}
