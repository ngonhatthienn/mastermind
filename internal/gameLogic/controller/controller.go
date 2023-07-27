package controller

import (
	"context"

	"intern2023/handler/ToProto"
	"intern2023/internal/gameLogic/handler"
	pb "intern2023/pb/game"
	"intern2023/share"

	"google.golang.org/grpc/metadata"
)

type Controller struct {
	service *handler.GameLogicHandler
	pb.UnimplementedServicesServer
}

func NewController() *Controller {
	Service := handler.NewService()
	return &Controller{Service, pb.UnimplementedServicesServer{}}
}

// GAME
// Create game in Redis database
func (c *Controller) CreateGame(ctx context.Context, in *pb.CreateGameRequest) (*pb.CreateGameReply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	status, _ := c.service.AuthorAndAuthn(md, "admin")

	if status.Code != 200 {
		return &pb.CreateGameReply{
			Code: status.Code, Message: status.Message,
		}, nil
	}
	sizeGame := 10
	GuessLimit := int(in.GuessLimit)
	c.service.CreateGame(sizeGame, GuessLimit)

	return &pb.CreateGameReply{Code: 200, Message: "Create game success!!!"}, nil
}

func (c *Controller) ListGame(ctx context.Context, in *pb.ListGameRequest) (*pb.ListGameReply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	status, _ := c.service.AuthorAndAuthn(md, "none")
	if status.Code != 200 {
		return &pb.ListGameReply{
			Code: status.Code, Message: status.Message,
		}, nil
	}

	// Check it's admin or not
	var isAdmin bool
	status, _ = c.service.AuthorAndAuthn(md, "admin")
	if status.Code == 200 {
		isAdmin = true
	} else {
		isAdmin = false
	}

	length, Games := c.service.ListGame()
	gameProtos := ToProto.ToListGameProto(Games, isAdmin)
	status = share.GenerateStatus(200, "")

	return &pb.ListGameReply{Code: status.Code, Message: status.Message, Length: int32(length), Games: gameProtos}, nil
}

// Get Random game
func (c *Controller) GetCurrent(ctx context.Context, in *pb.CurrentGameRequest) (*pb.CurrentGameReply, error) {
	// Check Auth
	md, _ := metadata.FromIncomingContext(ctx)
	status, IdUser := c.service.AuthorAndAuthn(md, "user")
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
	status, IdUser := c.service.AuthorAndAuthn(md, "user")
	if status.Code != 200 {
		return &pb.PickGameReply{Code: status.Code, Message: status.Message}, nil
	}
	//
	status, GameReply := c.service.PickGame(IdUser, int(in.IdGame))
	return &pb.PickGameReply{Code: status.Code, Message: status.Message, Game: GameReply}, nil
}

// Update Game
func (c *Controller) UpdateGame(ctx context.Context, in *pb.UpdateGameRequest) (*pb.UpdateGameReply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	status, _ := c.service.AuthorAndAuthn(md, "admin")
	if status.Code != 200 {
		return &pb.UpdateGameReply{
			Code: status.Code, Message: status.Message,
		}, nil
	}
	status = c.service.UpdateGame(30)
	return &pb.UpdateGameReply{Code: status.Code, Message: status.Message}, nil
}

// Play Game
func (c *Controller) PlayGame(ctx context.Context, in *pb.PlayGameRequest) (*pb.PlayGameReply, error) {
	// Check Auth
	md, _ := metadata.FromIncomingContext(ctx)
	status, IdUser := c.service.AuthorAndAuthn(md, "user")

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
	status, IdUser := c.service.AuthorAndAuthn(md, "user")
	if status.Code != 200 {
		return &pb.HintGameReply{Code: status.Code, Message: status.Message}, nil
	}
	// Check exist user
	status, res := c.service.HintGame(IdUser, in.Type)
	return &pb.HintGameReply{Code: status.Code, Message: status.Message, GameHint: res}, nil
}

// USER
func (c *Controller) LogIn(ctx context.Context, in *pb.LogInRequest) (*pb.LogInReply, error) {
	status, IdUser, userRole, ok := c.service.LogIn(in.Username, in.Password)
	if !ok {
		status := share.GenerateStatus(404, "User")
		return &pb.LogInReply{Code: status.Code, Message: status.Message}, nil
	}
	token := c.service.CreateToken(IdUser, userRole)
	return &pb.LogInReply{Code: status.Code, Message: status.Message, Token: token}, nil
}

func (c *Controller) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	Id, _ := c.service.CreateUser(in.Fullname, in.Username, in.Password, in.Email, in.Role)
	return &pb.CreateUserReply{XId: Id, Message: "Welcome " + in.Username}, nil
}

func (c *Controller) GetListUser(ctx context.Context, in *pb.ListUserRequest) (*pb.ListUserReply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	status, _ := c.service.AuthorAndAuthn(md, "admin")
	if status.Code != 200 {
		return &pb.ListUserReply{
			Code: status.Code, Message: status.Message,
		}, nil
	}
	Users, err := c.service.ListUsers()
	if err != nil {
		return &pb.ListUserReply{Code: 404, Message: "User Not Found"}, nil
	}
	Length := len(Users)
	userProtos := ToProto.ToListUserProto(Users)
	status = share.GenerateStatus(200, "")
	return &pb.ListUserReply{Code: status.Code, Message: status.Message, Length: int32(Length), Users: userProtos}, nil
}

// LEADERBOARD
func (c *Controller) GetLeaderBoard(ctx context.Context, in *pb.LeaderBoardRequest) (*pb.LeaderBoardReply, error) {
	// Check Auth
	md, _ := metadata.FromIncomingContext(ctx)
	status, IdUser := c.service.AuthorAndAuthn(md, "none")
	if status.Code != 200 {
		return &pb.LeaderBoardReply{
			Code: status.Code, Message: status.Message,
		}, nil
	}

	// Check it's admin or not
	var isAdmin bool
	status, _ = c.service.AuthorAndAuthn(md, "admin")
	if status.Code == 200 {
		isAdmin = true
	} else {
		isAdmin = false
	}

	status, leaderboard, UserRank, UserScore := c.service.GetLeaderBoard(int(in.IdGame), IdUser, int(in.Size), isAdmin)

	leaderboardProto := ToProto.ToLeaderBoardProto(leaderboard)
	return &pb.LeaderBoardReply{Code: status.Code, Message: status.Message, Ranks: leaderboardProto, UserRank: UserRank, UserScore: UserScore}, nil
}
