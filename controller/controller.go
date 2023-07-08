package controller

import (
	"context"

	"google.golang.org/grpc/metadata"

	"intern2023/model"
	pb "intern2023/pb"
	"intern2023/share"
)

type Controller struct {
	service *model.Service
	pb.UnimplementedServicesServer
}

func NewController(Service *model.Service) *Controller {
	return &Controller{Service, pb.UnimplementedServicesServer{}}
}

// GAME

// Create game in Redis database
func (c *Controller) CreateGame(ctx context.Context, in *pb.CreateGameRequest) (*pb.CreateGameReply, error) {
	sizeGame := 10
	GuessLimit := int(in.GuessLimit)
	c.service.CreateGame(sizeGame, GuessLimit)

	return &pb.CreateGameReply{Code: 200, Message: "Create game success!!!"}, nil
}

func (c *Controller) ListGame(ctx context.Context, in *pb.ListGameRequest) (*pb.ListGameReply, error) {
	length, Games := c.service.ListGame()
	return &pb.ListGameReply{Code: 200, Length: int32(length), Games: Games}, nil
}

// Get Random game
func (c *Controller) GetCurrent(ctx context.Context, in *pb.CurrentGameRequest) (*pb.CurrentGameReply, error) {
	// Check Auth
	md, _ := metadata.FromIncomingContext(ctx)
	// bearerToken := md.Get("authorization")
	// if len(bearerToken) <= 0 {
	// 	status := share.GenerateStatus(401, "")
	// 	return &pb.CurrentGameReply{Code: status.Code, Message: status.Message}, nil
	// }
	status, IdUser := c.service.Authorization(md)
	if status.Code != 200 {
		return &pb.CurrentGameReply{Code: status.Code, Message: status.Message}, nil
	}
	//
	status, GameReply := c.service.GetCurrent(IdUser)
	return &pb.CurrentGameReply{Code: status.Code, Message: status.Message, Game: GameReply}, nil
}

func (c *Controller) PickGame(ctx context.Context, in *pb.PickGameRequest) (*pb.PickGameReply, error) {
	// Check Auth
	md, _ := metadata.FromIncomingContext(ctx)
	// bearerToken := md.Get("authorization")
	// if len(bearerToken) <= 0 {
	// 	status := share.GenerateStatus(401, "")
	// 	return &pb.PickGameReply{Code: status.Code, Message: status.Message}, nil
	// }
	status, IdUser := c.service.Authorization(md)
	if status.Code != 200 {
		return &pb.PickGameReply{Code: status.Code, Message: status.Message}, nil
	}
	//
	status, GameReply := c.service.PickGame(IdUser, int(in.IdGame))
	return &pb.PickGameReply{Code: status.Code, Message: status.Message, Game: GameReply}, nil
}

// Update Game
func (c *Controller) UpdateGame(ctx context.Context, in *pb.UpdateGameRequest) (*pb.UpdateGameReply, error) {
	status := c.service.UpdateGame(int(in.GuessLimit))
	return &pb.UpdateGameReply{Code: status.Code, Message: status.Message}, nil
}

// Play Game
func (c *Controller) PlayGame(ctx context.Context, in *pb.PlayGameRequest) (*pb.PlayGameReply, error) {
	// Check Auth
	md, _ := metadata.FromIncomingContext(ctx)
	// bearerToken := md.Get("authorization")
	// if len(bearerToken) <= 0 {
	// 	status := share.GenerateStatus(401, "")
	// 	return &pb.PlayGameReply{Code: status.Code, Message: status.Message}, nil
	// }
	status, IdUser := c.service.Authorization(md)
	if status.Code != 200 {
		return &pb.PlayGameReply{Code: status.Code, Message: status.Message}, nil
	}
	status, guessLeft, listHistory := c.service.PlayGame(IdUser, in.UserGuess)
	return &pb.PlayGameReply{Code: status.Code, Message: status.Message, GuessesLeft: int32(guessLeft), Result: listHistory}, nil
}

// Hint Game

func (c *Controller) HintGame(ctx context.Context, in *pb.HintGameRequest) (*pb.HintGameReply, error) {
	// Check Auth
	md, _ := metadata.FromIncomingContext(ctx)
	// bearerToken := md.Get("authorization")
	// if len(bearerToken) <= 0 {
	// 	status := share.GenerateStatus(401, "")
	// 	return &pb.HintGameReply{Code: status.Code, Message: status.Message}, nil
	// }
	status, IdUser := c.service.Authorization(md)
	if status.Code != 200 {
		return &pb.HintGameReply{Code: status.Code, Message: status.Message}, nil
	}
	// Check exist user
	status, res := c.service.HintGame(IdUser, in.Type)
	return &pb.HintGameReply{Code: status.Code, Message: status.Message, GameHint: res}, nil
}

// USER
func (c *Controller) LogIn(ctx context.Context, in *pb.LogInRequest) (*pb.LogInReply, error) {
	status, IdUser, ok := c.service.LogIn(in.Username, in.Password)
	if !ok {
		status := share.GenerateStatus(404, "User")
		return &pb.LogInReply{Code: status.Code, Message: status.Message}, nil
	}
	token := c.service.CreateToken(IdUser)
	return &pb.LogInReply{Code: status.Code, Message: status.Message, Token: token}, nil
}

func (c *Controller) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	Id, _ := c.service.CreateUser(in.Name, in.Password)
	return &pb.CreateUserReply{XId: Id, Message: "Welcome " + in.Name}, nil
}

func (c *Controller) GetListUser(ctx context.Context, in *pb.ListUserRequest) (*pb.ListUserReply, error) {
	Users, _ := c.service.ListUsers()
	Length := len(Users)
	return &pb.ListUserReply{Length: int32(Length), Users: Users}, nil
}

// LEADERBOARD
func (c *Controller) GetLeaderBoard(ctx context.Context, in *pb.LeaderBoardRequest) (*pb.LeaderBoardReply, error) {
	// Check Auth
	md, _ := metadata.FromIncomingContext(ctx)
	// bearerToken := md.Get("authorization")
	// if len(bearerToken) <= 0 {
	// 	status := share.GenerateStatus(401, "")
	// 	return &pb.LeaderBoardReply{Code: status.Code, Message: status.Message}, nil
	// }
	status, IdUser := c.service.Authorization(md)
	if status.Code != 200 {
		return &pb.LeaderBoardReply{Code: status.Code, Message: status.Message}, nil
	}
	status, leaderboardData, UserRank, UserScore := c.service.GetLeaderBoard(int(in.IdGame), IdUser, int(in.Size))

	return &pb.LeaderBoardReply{Code: status.Code, Message: status.Message, Ranks: leaderboardData, UserRank: UserRank, UserScore: UserScore}, nil
}
