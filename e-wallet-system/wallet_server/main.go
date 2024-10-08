package main

import (
	"context"
	"log"
	"net"
	"sync"

	"wallet_server/proto/wallet" // Pastikan path import sesuai dengan hasil kompilasi Protobuf

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type WalletServiceServer struct {
	wallet.UnimplementedWalletServiceServer // Tambahkan ini untuk kompatibilitas
	wallets                                 map[string]float64
	transactions                            map[string][]*wallet.Transaction
	mu                                      sync.Mutex
}

// TopUp mengimplementasikan method TopUp yang didefinisikan di wallet.proto
func (s *WalletServiceServer) TopUp(ctx context.Context, req *wallet.TopUpRequest) (*wallet.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Memperbarui saldo wallet user
	s.wallets[req.UserId] += req.Amount

	// Membuat transaksi baru untuk top-up
	transaction := &wallet.Transaction{
		Id:        "txn-topup", // Sebaiknya menggunakan ID unik sebenarnya
		Type:      "topup",
		Amount:    req.Amount,
		CreatedAt: "now", // Ganti dengan timestamp aktual
	}

	// Menambahkan transaksi ke riwayat transaksi user
	s.transactions[req.UserId] = append(s.transactions[req.UserId], transaction)

	// Mengembalikan saldo wallet user yang telah diperbarui
	return &wallet.UserResponse{
		UserId:  req.UserId,
		Balance: s.wallets[req.UserId],
	}, nil
}

// Transfer mengimplementasikan method Transfer yang didefinisikan di wallet.proto
func (s *WalletServiceServer) Transfer(ctx context.Context, req *wallet.TransferRequest) (*wallet.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Memeriksa apakah pengirim memiliki saldo yang cukup
	if s.wallets[req.FromUserId] < req.Amount {
		return nil, status.Errorf(400, "Saldo tidak mencukupi")
	}

	// Mengurangi saldo pengirim
	s.wallets[req.FromUserId] -= req.Amount

	// Menambahkan saldo ke penerima
	s.wallets[req.ToUserId] += req.Amount

	// Membuat transaksi untuk pengirim dan penerima
	transaction := &wallet.Transaction{
		Id:        "txn-transfer", // Sebaiknya menggunakan ID unik sebenarnya
		Type:      "transfer",
		Amount:    req.Amount,
		CreatedAt: "now", // Ganti dengan timestamp aktual
	}

	// Menambahkan transaksi ke riwayat pengirim dan penerima
	s.transactions[req.FromUserId] = append(s.transactions[req.FromUserId], transaction)
	s.transactions[req.ToUserId] = append(s.transactions[req.ToUserId], transaction)

	// Mengembalikan saldo terbaru pengirim
	return &wallet.UserResponse{
		UserId:  req.FromUserId,
		Balance: s.wallets[req.FromUserId],
	}, nil
}

// GetTransactionList mengimplementasikan method GetTransactionList yang didefinisikan di wallet.proto
func (s *WalletServiceServer) GetTransactionList(ctx context.Context, req *wallet.TransactionListRequest) (*wallet.TransactionListResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Mendapatkan daftar transaksi user
	transactions, ok := s.transactions[req.UserId]
	if !ok {
		return &wallet.TransactionListResponse{
			Transactions: []*wallet.Transaction{},
		}, nil
	}

	// Mengembalikan daftar transaksi
	return &wallet.TransactionListResponse{
		Transactions: transactions,
	}, nil
}

func main() {
	// Membuat instance server gRPC baru
	server := grpc.NewServer()

	// Inisialisasi WalletServiceServer dengan map in-memory
	walletService := &WalletServiceServer{
		wallets:      make(map[string]float64),
		transactions: make(map[string][]*wallet.Transaction),
	}

	// Mendaftarkan WalletServiceServer ke gRPC server
	wallet.RegisterWalletServiceServer(server, walletService)

	// Mendengarkan di port 50052
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Gagal mendengarkan di port 50052: %v", err)
	}

	log.Println("Wallet server berjalan di port 50052")

	// Menjalankan gRPC server
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Gagal menjalankan server gRPC: %v", err)
	}
}
