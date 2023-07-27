package service

import (
	authpb "intern2023/pb/auth"

	"google.golang.org/grpc"
)

type AuthService struct {
	AuthClient authpb.AuthClient
	Connection *grpc.ClientConn
}

func newAuthClient() (authpb.AuthClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	// defer conn.Close()
	return authpb.NewAuthClient(conn), conn, nil
}

func NewAuthService() *AuthService {
	client, conn, err := newAuthClient()
	if err != nil {
		panic(err)
	}

	return &AuthService{
		AuthClient: client,
		Connection: conn,
	}
}

func (as *AuthService) CloseConnection() {
	if as.Connection != nil {
		as.Connection.Close()
	}
}
