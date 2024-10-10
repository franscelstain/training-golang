package main

import (
	walletpb "eWalletSystem/proto/wallet/v1"
	"eWalletSystem/wallet_server/handler"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052: %v", err)
	}

	grpcServer := grpc.NewServer()

	walletService := handler.NewWalletServiceServer()
	walletpb.RegisterWalletServiceServer(grpcServer, walletService)

	log.Println("Wallet server is running on port 50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
