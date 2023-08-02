package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"intern2023/internal/auth/controller"
	pb "intern2023/pb/auth"
	config "intern2023/handler/Config"
)

func main() {
	config := config.GetConfig()
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", config.Auth.GRPC_URL)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	// Create a gRPC server object
	s := grpc.NewServer()
	Controller := controller.NewController()
	pb.RegisterAuthServer(s, Controller)
	// Serve gRPC server
	log.Println("Serving gRPC on ", config.Auth.GRPC_URL)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
