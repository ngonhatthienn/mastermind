package main

import (
	// "context"
	"log"
	"net"

	"google.golang.org/grpc"

	"intern2023/internal/auth/controller"
	pb "intern2023/pb/auth"
)

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	// Create a gRPC server object
	s := grpc.NewServer()
	Controller := controller.NewController()
	pb.RegisterAuthServer(s, Controller)
	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:9080")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
