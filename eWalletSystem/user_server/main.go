package main

import (
	walletpb "eWalletSystem/proto/wallet/v1"
	"eWalletSystem/user_server/handler"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register UserServiceServer di gRPC server
	userService := handler.NewUserServiceServer()
	walletpb.RegisterUserServiceServer(grpcServer, userService)

	log.Println("User server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
