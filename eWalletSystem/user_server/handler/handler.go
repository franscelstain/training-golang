package handler

import (
	"context"
	walletpb "eWalletSystem/proto/wallet/v1"
	"eWalletSystem/user_server/entity"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	walletpb.UnimplementedUserServiceServer
	Users map[string]*entity.User
}

// Inisialisasi UserServiceServer dengan data pengguna
func NewUserServiceServer() *UserServiceServer {
	return &UserServiceServer{
		Users: map[string]*entity.User{
			"user1": {UserID: "user1", Name: "John Doe", Email: "john@example.com", Balance: 1000},
			"user2": {UserID: "user2", Name: "Jane Smith", Email: "jane@example.com", Balance: 1500},
		},
	}
}

// Implementasi GetUser untuk mengembalikan data user berdasarkan ID
func (s *UserServiceServer) GetUser(ctx context.Context, req *walletpb.UserRequest) (*walletpb.UserResponse, error) {
	user, exists := s.Users[req.UserId]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	return &walletpb.UserResponse{
		UserId:  user.UserID,
		Name:    user.Name,
		Email:   user.Email,
		Balance: user.Balance,
	}, nil
}
