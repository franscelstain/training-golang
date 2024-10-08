package main

import (
	"context"
	"log"
	"net"

	"user_server/proto/wallet" // Pastikan path import sesuai dengan hasil kompilasi Protobuf

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	wallet.UnimplementedUserServiceServer // Tambahkan ini untuk kompatibilitas
	users                                 map[string]*wallet.UserResponse
}

// GetUser mengimplementasikan method GetUser yang didefinisikan di wallet.proto
func (s *UserServiceServer) GetUser(ctx context.Context, req *wallet.UserRequest) (*wallet.UserResponse, error) {
	user, exists := s.users[req.UserId]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "User tidak ditemukan")
	}
	return user, nil
}

func main() {
	// Membuat instance server gRPC baru
	server := grpc.NewServer()

	// Inisialisasi UserServiceServer dengan data in-memory
	userService := &UserServiceServer{
		users: map[string]*wallet.UserResponse{
			"user1": {UserId: "user1", Name: "John Doe", Email: "john@example.com", Balance: 1000},
			"user2": {UserId: "user2", Name: "Jane Smith", Email: "jane@example.com", Balance: 1500},
		},
	}

	// Mendaftarkan UserServiceServer ke gRPC server
	wallet.RegisterUserServiceServer(server, userService)

	// Mendengarkan di port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Gagal mendengarkan di port 50051: %v", err)
	}

	log.Println("User server berjalan di port 50051")

	// Menjalankan gRPC server
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Gagal menjalankan server gRPC: %v", err)
	}
}
