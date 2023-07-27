package controller
import (
	"context"

	"intern2023/internal/auth/handler"
	pb "intern2023/pb/auth"
)

type AuthServer struct {
	authHandler  *handler.AuthHandler
	pb.UnimplementedAuthServer
}
func NewController() *AuthServer {
	handler := handler.NewService()
	return &AuthServer{handler, pb.UnimplementedAuthServer{}}
}

func (s *AuthServer) CheckUser(ctx context.Context, in *pb.CheckUserRequest) (*pb.CheckUserReply, error) {
	exists,err := s.authHandler.CheckUser(in.UserId, in.SessionId)
	return &pb.CheckUserReply{Exists: exists}, err
}