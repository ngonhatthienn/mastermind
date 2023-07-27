package auth

import (
	"context"
	"log"
	"net"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"intern2023/database"
	pb "intern2023/pb/auth"
	"intern2023/share"
)

type AuthServer struct {
	redisClient *redis.Client
	pb.UnimplementedAuthServer
}

func NewAuthService() *AuthServer {
	redisClient, _ := database.ConnectRedisDatabase()

	return &AuthServer{redisClient, pb.UnimplementedAuthServer{}}
}

func (s *AuthServer) CheckUser(ctx context.Context, in *pb.CheckUserRequest) (*pb.CheckUserReply, error) {
	if in.UserId == "" || in.SessionId == "" {
		return &pb.CheckUserReply{Exists: false}, nil
	}
	// Check if the user exists in the Redis database
	IdSession, err := s.redisClient.Get(context.Background(), share.UserPatternSession(in.UserId)).Result()
	if err != nil || IdSession == "" {
		return &pb.CheckUserReply{Exists: false}, nil
	}
	if IdSession != in.SessionId {
		return &pb.CheckUserReply{Exists: false}, nil
	}
	return &pb.CheckUserReply{Exists: true}, nil
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	// Create a gRPC server object
	s := grpc.NewServer()
	Controller := NewAuthService()
	pb.RegisterAuthServer(s, Controller)
	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:9080")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
