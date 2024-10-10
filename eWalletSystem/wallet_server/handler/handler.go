package handler

import (
	"context"
	walletpb "eWalletSystem/proto/wallet/v1"
	"eWalletSystem/wallet_server/entity"
)

type WalletServiceServer struct {
	walletpb.UnimplementedWalletServiceServer
	Transactions map[string][]*entity.Transaction
}

// Inisialisasi WalletServiceServer
func NewWalletServiceServer() *WalletServiceServer {
	return &WalletServiceServer{
		Transactions: make(map[string][]*entity.Transaction),
	}
}

// Implementasi TopUp: menambahkan saldo ke user
func (s *WalletServiceServer) TopUp(ctx context.Context, req *walletpb.TopUpRequest) (*walletpb.UserResponse, error) {
	// Contoh implementasi: saldo setelah top-up
	return &walletpb.UserResponse{
		UserId:  req.UserId,
		Balance: 1500,
	}, nil
}

// Implementasi Transfer: mentransfer saldo dari satu user ke user lain
func (s *WalletServiceServer) Transfer(ctx context.Context, req *walletpb.TransferRequest) (*walletpb.UserResponse, error) {
	// Contoh implementasi: saldo setelah transfer
	return &walletpb.UserResponse{
		UserId:  req.FromUserId,
		Balance: 800,
	}, nil
}
