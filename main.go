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

	// "github.com/golobby/dotenv"
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

func (s *server) CreateGame(ctx context.Context, in *pb.CreateGameRequest) (*pb.CreateGameReply, error) {
	client, _ := database.CreateRedisDatabase()
	game.CreateGames(client, 10, int(in.WrongLimit))
	return &pb.CreateGameReply{Message: "Create game success!!!"}, nil
}

// List Game
func (s *server) ListGame(ctx context.Context, in *pb.ListGameRequest) (*pb.ListGameReply, error) {
	client, _ := database.CreateRedisDatabase()

	length, Games := game.GetListGame(client)
	return &pb.ListGameReply{Length: int32(length), Games: Games}, nil
}

func (s *server) GetCurrent(ctx context.Context, in *pb.CurrentGameRequest) (*pb.CurrentGameReply, error) {
	client, _ := database.CreateRedisDatabase()
	_, IdGame := session.CreateUserSession(client, in.IdUser)
	Game := game.GetGameValue(client, IdGame)
	GameReply := pb.GameReply{XId: Game.XId, WrongLimit: Game.WrongLimit}
	return &pb.CurrentGameReply{Game: &GameReply}, nil
}

// Update Game
func (s *server) UpdateGame(ctx context.Context, in *pb.UpdateGameRequest) (*pb.UpdateGameReply, error) {
	client, _ := database.CreateRedisDatabase()
	game.UpdateGame(client, int(in.WrongLimit))
	return &pb.UpdateGameReply{Message: "Update Game Success"}, nil
}

// Play Game
func (s *server) PlayGame(ctx context.Context, in *pb.PlayGameRequest) (*pb.PlayGameReply, error) {
	// Get Game by key
	client, _ := database.CreateRedisDatabase()

	IdUser := int(in.IdUser)
	IdUserString := strconv.Itoa(IdUser)

	getSessions := "session:" + IdUserString + "*"
	keySessions := database.Scan(client, getSessions)


	var keySession string
	if len(keySessions) == 0 {
		_, IdGame := session.CreateUserSession(client, in.IdUser)
		keySession = "session:" + strconv.Itoa(IdUser) + ":" + strconv.Itoa(IdGame)
	} else {
		keySession = keySessions[0]
	}
	
	// Update history when play

	guessLeftString, _ := database.HGet(client, keySession, "guessLeft")
	guessLeft, _ := strconv.Atoi(guessLeftString)
	if guessLeft == 0 {
		return &pb.PlayGameReply{Message: "You lose!!!"}, nil
	} else {
		guessLeft--
		guessLeftString = strconv.Itoa(guessLeft)
		database.HSet(client, keySession, "guessLeft", guessLeft)
	}
	

	IdGameString := shareFunc.GetKeyElement(keySession, 2)
	IdGame, _ := strconv.Atoi(IdGameString)
	getGameValue := game.GetGameValue(client, IdGame)

	rightNumber, rightPosition := gameApp.OutputGame(in.UserGuess, getGameValue.Game)

	

	if rightNumber == rightPosition && rightNumber == 5 {
		database.HSet(client, keySession, "isWin", true)

		// Handle Time

		timeStart, _ := database.HGet_int64(client, keySession, "timeStart")
		savedTime := time.Unix(timeStart, 0)

		diffInSeconds := 5000 - time.Now().Sub(savedTime).Seconds()

		// Get right and pos
		right, _ := database.HGet(client, keySession, "sumRight")
		pos, _ := database.HGet(client, keySession, "sumPos")
		rightInt, _ := strconv.Atoi(right)
		posInt, _ := strconv.Atoi(pos)

		score := int(diffInSeconds) + guessLeft*100 + (rightInt+posInt)*2
		fmt.Println(IdUserString, "Difference Time", diffInSeconds)
		_ = leaderboard.AddScore(client, IdUserString, IdGameString, int64(score))

		return &pb.PlayGameReply{RightNumber: int32(rightNumber), RightPosition: int32(rightPosition), Message: "You win!!!"}, nil
	}
	var listHistory []*pb.ListHistory
	listHistory, _ = session.PushAndGetHistory(client, keySession, in.UserGuess, int32(rightNumber), int32(rightPosition))

	return &pb.PlayGameReply{RightNumber: int32(rightNumber), RightPosition: int32(rightPosition), Message: "Try your best !!!",GuessesLeft: guessLeftString, OldResult: listHistory}, nil
}

func (s *server) HintGame(ctx context.Context, in *pb.HintGameRequest) (*pb.HintGameReply, error) {
	client, _ := database.CreateRedisDatabase()

	keySessionPattern := "session:" + strconv.Itoa(int(in.IdUser)) + ":*"

	keySessions, _ := database.Keys(client, keySessionPattern)

	IdGame := shareFunc.GetKeyElement(keySessions[0], 2)
	key := "game:" + IdGame

	val, _ := database.Get(client, key)
	var Result *pb.Game
	_ = json.Unmarshal([]byte(val), &Result)
	res := game.GenerateHint(Result.Game, in.Type)
	return &pb.HintGameReply{GameHint: res}, nil
}

func (s *server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	client, _ := database.CreateRedisDatabase()
	Id, Name := user.CreateUser(client, in)
	return &pb.CreateUserReply{XId: Id, Message: "Welcome " + Name}, nil
}

func (s *server) GetListUser(ctx context.Context, in *pb.ListUserRequest) (*pb.ListUserReply, error) {
	client, _ := database.CreateRedisDatabase()
	Users := user.GetListUser(client)
	return &pb.ListUserReply{Length: int32(len(Users)), Users: Users}, nil
}

func (s *server) GetLeaderBoard(ctx context.Context, in *pb.LeaderBoardRequest) (*pb.LeaderBoardReply, error) {
	client, _ := database.CreateRedisDatabase()
	leaderboardData, _ := leaderboard.GetLeaderboard(client, strconv.Itoa(int(in.IdGame)), in.Size)
	return &pb.LeaderBoardReply{Data: leaderboardData}, nil
}

func main() {
	// Redis connect
	// client, _ := database.CreateRedisDatabase()

	// keys, err := database.Keys(client, "game:*")
	// fmt.Println(len(keys))
	// if len(keys) == 0 {
	// 	game.CreateGames(client, 10, 10)
	// }
	// Load environment variables from .env file
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

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
