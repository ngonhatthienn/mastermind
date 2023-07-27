package service

import (
	"context"

	"intern2023/internal/auth/service"
	authpb "intern2023/pb/auth"
)

type UserService struct {
	authService *service.AuthService
}

func NewUserService() *UserService {
	authService := service.NewAuthService()

	return &UserService{
		authService: authService,
	}
}

func (us *UserService) CheckSessionId(UserId string, SessionId string) bool{
	data := &authpb.CheckUserRequest{
		UserId:    UserId,
		SessionId: SessionId,
	}

	resp, err := us.authService.AuthClient.CheckUser(context.Background(), data)
	if err != nil {
		panic(err)
	}
	us.authService.CloseConnection()
	return resp.Exists
}
