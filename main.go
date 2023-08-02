package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	config "intern2023/handler/Config"
	"intern2023/internal/gameLogic/controller"
	pb "intern2023/pb/game"
)

func main() {
	config := config.GetConfig()
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", config.GameLogic.GRPC_URL)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	Controller := controller.NewController()
	pb.RegisterServicesServer(s, Controller)
	// Serve gRPC server
	log.Println("Serving gRPC on ", config.GameLogic.GRPC_URL)
	go func() {
		log.Fatalln(s.Serve(lis))
	}()
	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		config.GameLogic.GRPC_URL,
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
		Addr:    config.GameLogic.GRPC_GATEWAY_URL,
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on ", config.GameLogic.GRPC_GATEWAY_URL)
	log.Fatalln(gwServer.ListenAndServe())
}
