package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	gameApp "intern2023/app"
	"intern2023/database"
	game "intern2023/handler/Game"
	leaderboard "intern2023/handler/Leaderboard"
	session "intern2023/handler/Session"
	user "intern2023/handler/User"
	shareFunc "intern2023/share"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "intern2023/pb"
)

type server struct {
	pb.UnimplementedServicesServer
}

func NewServer() *server {
	return &server{}
}

// Init game in Mongo database
func (s *server) InitGame(ctx context.Context, in *pb.InitGameRequest) (*pb.InitGameReply, error) {
	client := database.CreateMongoDBConnection()
	game.CreateGamesMongo(client, int(in.GameSize))
	return &pb.InitGameReply{Code: 200, Message: "Init game success!!!"}, nil
}

// Create game in Redis database
func (s *server) CreateGame(ctx context.Context, in *pb.CreateGameRequest) (*pb.CreateGameReply, error) {
	client, _ := database.CreateRedisDatabase()
	game.CreateGames(client, 10, int(in.GuessLimit))

	return &pb.CreateGameReply{Code: 200, Message: "Create game success!!!"}, nil
}

// List Game
func (s *server) ListGame(ctx context.Context, in *pb.ListGameRequest) (*pb.ListGameReply, error) {
	client, _ := database.CreateRedisDatabase()
	length, Games := game.GetListGame(client)
	if length == 0 {
		game.CreateGames(client, 10, 30)
		length, Games = game.GetListGame(client)
	}
	return &pb.ListGameReply{Code: 200, Length: int32(length), Games: Games}, nil
}

func (s *server) GetCurrent(ctx context.Context, in *pb.CurrentGameRequest) (*pb.CurrentGameReply, error) {
	client, _ := database.CreateRedisDatabase()
	// Check exist user
	IdUser := int(in.IdUser)
	if user.CheckExistUser(client, IdUser) == false {
		status := shareFunc.GenerateStatus(404, "User")
		return &pb.CurrentGameReply{Code: status.Code, Message: status.Message}, nil
	}
	_, IdGame := session.CreateUserSession(client, in.IdUser)
	Game := game.GetGameValue(client, IdGame)
	// GameReply := pb.GameReply{XId: Game.XId, GuessLimit: Game.GuessLimit}
	GameReply := pb.GameReply{XId: int32(Game.ID), GuessLimit: int32(Game.GuessLimit)}

	status := shareFunc.GenerateStatus(200, "Get current")
	return &pb.CurrentGameReply{Code: status.Code, Message: status.Message, Game: &GameReply}, nil
}

func (s *server) PickGame(ctx context.Context, in *pb.PickGameRequest) (*pb.PickGameReply, error) {
	client, _ := database.CreateRedisDatabase()
	// Check exist user
	IdUser := int(in.IdUser)
	if user.CheckExistUser(client, IdUser) == false {
		status := shareFunc.GenerateStatus(404, "User")
		return &pb.PickGameReply{Code: status.Code, Message: status.Message}, nil
	}
	// Check exist game
	IdGame := int(in.IdGame)
	if game.CheckExistGame(client, IdGame) == false {
		status := shareFunc.GenerateStatus(404, "Game")
		return &pb.PickGameReply{Code: status.Code, Message: status.Message}, nil
	}
	// Handle pick game
	session.CreateSessionWithId(client, in.IdUser, in.IdGame)
	Game := game.GetGameValue(client, int(in.IdGame))
	// GameReply := pb.GameReply{XId: Game.XId, GuessLimit: Game.GuessLimit}
	GameReply := pb.GameReply{XId: int32(Game.ID), GuessLimit: int32(Game.GuessLimit)}
	status := shareFunc.GenerateStatus(200, "Pick game")

	return &pb.PickGameReply{Code: status.Code, Message: status.Message, Game: &GameReply}, nil
}

// Update Game
func (s *server) UpdateGame(ctx context.Context, in *pb.UpdateGameRequest) (*pb.UpdateGameReply, error) {
	client, _ := database.CreateRedisDatabase()
	game.UpdateGame(client, int(in.GuessLimit))
	return &pb.UpdateGameReply{Message: "Update Game Success"}, nil
}

// Play Game
func (s *server) PlayGame(ctx context.Context, in *pb.PlayGameRequest) (*pb.PlayGameReply, error) {
	client, _ := database.CreateRedisDatabase()

	IdUser := int(in.IdUser)
	if user.CheckExistUser(client, IdUser) == false {
		status := shareFunc.GenerateStatus(404, "User")
		return &pb.PlayGameReply{Code: status.Code, Message: status.Message}, nil
	}
	IdUserString := strconv.Itoa(IdUser)

	// Check Exist session or not
	getSessions := "session:" + IdUserString + ":*"
	keySessions, _ := client.Keys(context.Background(), getSessions).Result()
	var keySession string
	if len(keySessions) == 0 {
		_, IdGame := session.CreateUserSession(client, in.IdUser)
		keySession = "session:" + strconv.Itoa(IdUser) + ":" + strconv.Itoa(IdGame)
	} else {
		keySession = keySessions[0]
	}

	// Update history when play
	checkHistory := time.Now()

	isWin, _ := client.HGet(context.Background(), keySession, "isWin").Bool()
	if isWin {
		return &pb.PlayGameReply{Message: "You'd already won, please get another game"}, nil
	}
	guessLeft, _ := client.HGet(context.Background(), keySession, "guessLeft").Int()

	if guessLeft == 0 {
		return &pb.PlayGameReply{Message: "You lose!!!"}, nil
	} else {
		guessLeft--
		client.HSet(context.Background(), keySession, "guessLeft", guessLeft)
	}

	IdGameString := shareFunc.GetKeyElement(keySession, 2)
	IdGame, _ := strconv.Atoi(IdGameString)
	getGameValue := game.GetGameValue(client, IdGame)

	rightNumber, rightPosition := gameApp.OutputGame(in.UserGuess, getGameValue.Game)

	checkHistory1 := time.Since(checkHistory)
	fmt.Println("Update history when play", checkHistory1.Milliseconds())

	// If user win this game
	if rightNumber == rightPosition && rightNumber == 5 {
		client.HSet(context.Background(), keySession, "isWin", true)

		// Handle Time

		timeStart, _ := client.HGet(context.Background(), keySession, "timeStart").Int64()
		savedTime := time.Unix(timeStart, 0)

		diffInSeconds := 5000 - time.Now().Sub(savedTime).Seconds()

		// Get right and pos
		right, _ := client.HGet(context.Background(), keySession, "sumRight").Int()
		pos, _ := client.HGet(context.Background(), keySession, "sumPos").Int()

		score := int(diffInSeconds) + guessLeft*100 + (right+pos)*2
		_ = leaderboard.AddScore(client, IdUserString, IdGameString, int64(score))

		return &pb.PlayGameReply{Code: 200, Message: "You win!!!"}, nil
	}
	check := time.Now()

	var listHistory []*pb.ListHistory
	listHistory, _ = session.PushAndGetHistory(client, keySession, in.UserGuess, int32(rightNumber), int32(rightPosition))

	check1 := time.Since(check)
	fmt.Println("Time check PushAndGetHistory", check1.Milliseconds())

	return &pb.PlayGameReply{Code: 200, Message: "Try your best !!!", GuessesLeft: int32(guessLeft), Result: listHistory}, nil
}

func (s *server) HintGame(ctx context.Context, in *pb.HintGameRequest) (*pb.HintGameReply, error) {
	client, _ := database.CreateRedisDatabase()
	// Check exist user
	IdUser := int(in.IdUser)
	if user.CheckExistUser(client, IdUser) == false {
		return &pb.HintGameReply{Code: 404}, nil
	}

	keySessionPattern := "session:" + strconv.Itoa(int(in.IdUser)) + ":*"

	keySessions, _ := database.Keys(client, keySessionPattern)
	// The session of this user is not exists or expired
	if len(keySessions) == 0 {
		return &pb.HintGameReply{Code: 404}, nil
	}

	IdGame := shareFunc.GetKeyElement(keySessions[0], 2)
	key := "game:" + IdGame

	val, _ := database.Get(client, key)

	// The game in this session is not exists or expired
	if val == "" {
		return &pb.HintGameReply{Code: 404}, nil
	}
	// var Result *pb.Game
	var Result game.GameItem

	_ = json.Unmarshal([]byte(val), &Result)

	res, isSuccess := game.GenerateHint(Result.Game, in.Type)
	code := int32(200)
	if !isSuccess {
		code = 400
	}
	return &pb.HintGameReply{Code: code, GameHint: res}, nil
}

func (s *server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	client, _ := database.CreateRedisDatabase()
	Id, Name := user.CreateUser(client, in)
	return &pb.CreateUserReply{XId: Id, Message: "Welcome " + Name}, nil
}

func (s *server) GetListUser(ctx context.Context, in *pb.ListUserRequest) (*pb.ListUserReply, error) {
	client, _ := database.CreateRedisDatabase()

	Length, Users := user.GetListUser(client)
	return &pb.ListUserReply{Length: int32(Length), Users: Users}, nil
}

func (s *server) GetLeaderBoard(ctx context.Context, in *pb.LeaderBoardRequest) (*pb.LeaderBoardReply, error) {
	client, _ := database.CreateRedisDatabase()
	// Check exist game
	IdGame := int(in.IdGame)
	if game.CheckExistGame(client, IdGame) == false {
		return &pb.LeaderBoardReply{Code: 404}, nil
	}
	leaderboardData, _ := leaderboard.GetLeaderboard(client, strconv.Itoa(int(in.IdGame)), in.Size)
	return &pb.LeaderBoardReply{Data: leaderboardData}, nil
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	// client, _ := database.CreateRedisDatabase()
	// game.DeleteGames(client)

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	pb.RegisterServicesServer(s, &server{})
	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = pb.RegisterServicesHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
