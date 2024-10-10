package main

import (
	"context"
	"log"
	"net/http"

	walletpb "eWalletSystem/proto/wallet/v1"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Koneksi ke User Server di port 50051
	userConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Gagal terhubung ke User Server: %v", err)
	}
	defer userConn.Close()
	userClient := walletpb.NewUserServiceClient(userConn)

	// Koneksi ke Wallet Server di port 50052
	walletConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Gagal terhubung ke Wallet Server: %v", err)
	}
	defer walletConn.Close()
	walletClient := walletpb.NewWalletServiceClient(walletConn)

	// Endpoint untuk mendapatkan data user beserta saldo dan transaksi
	router.GET("/user/:id", func(c *gin.Context) {
		userId := c.Param("id")

		// Memanggil User Server untuk mendapatkan data user
		userResp, err := userClient.GetUser(context.Background(), &walletpb.UserRequest{UserId: userId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User tidak ditemukan"})
			return
		}

		// Memanggil Wallet Server untuk mendapatkan daftar transaksi user
		walletResp, err := walletClient.GetTransactionList(context.Background(), &walletpb.TransactionListRequest{UserId: userId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan transaksi wallet"})
			return
		}

		// Menggabungkan data user dan transaksi
		c.JSON(http.StatusOK, gin.H{
			"user":         userResp,
			"transactions": walletResp.Transactions,
		})
	})

	// Endpoint untuk top-up saldo
	router.POST("/wallet/topup", func(c *gin.Context) {
		var req walletpb.TopUpRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Memanggil Wallet Server untuk melakukan top-up
		resp, err := walletClient.TopUp(context.Background(), &walletpb.TopUpRequest{
			UserId: req.UserId,
			Amount: req.Amount,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal melakukan top-up"})
			return
		}

		// Mengembalikan saldo yang diperbarui
		c.JSON(http.StatusOK, resp)
	})

	// Endpoint untuk transfer saldo antar user
	router.POST("/wallet/transfer", func(c *gin.Context) {
		var req walletpb.TransferRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Memanggil Wallet Server untuk transfer saldo
		resp, err := walletClient.Transfer(context.Background(), &walletpb.TransferRequest{
			FromUserId: req.FromUserId,
			ToUserId:   req.ToUserId,
			Amount:     req.Amount,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal melakukan transfer saldo"})
			return
		}

		// Mengembalikan saldo terbaru pengirim
		c.JSON(http.StatusOK, resp)
	})

	// Endpoint untuk melihat daftar transaksi user
	router.GET("/wallet/transactions", func(c *gin.Context) {
		userId := c.Query("user_id")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}

		// Memanggil Wallet Server untuk mendapatkan daftar transaksi
		resp, err := walletClient.GetTransactionList(context.Background(), &walletpb.TransactionListRequest{UserId: userId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan transaksi"})
			return
		}

		// Mengembalikan daftar transaksi
		c.JSON(http.StatusOK, resp)
	})

	// Menjalankan server di port 8080
	router.Run(":8080")
}
