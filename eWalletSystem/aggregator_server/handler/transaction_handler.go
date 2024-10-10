package handler

import (
	"context"
	walletpb "eWalletSystem/proto/wallet/v1"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TransactionHandler struct {
	walletClient walletpb.WalletServiceClient
}

func NewTransactionHandler() *TransactionHandler {
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Failed to connect to Wallet Server")
	}
	client := walletpb.NewWalletServiceClient(conn)
	return &TransactionHandler{
		walletClient: client,
	}
}

// Fungsi untuk top-up saldo ke wallet
func (h *TransactionHandler) TopUp(c *gin.Context) {
	var req walletpb.TopUpRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	res, err := h.walletClient.TopUp(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// Fungsi untuk transfer saldo antar user
func (h *TransactionHandler) Transfer(c *gin.Context) {
	var req walletpb.TransferRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	res, err := h.walletClient.Transfer(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// Fungsi untuk mendapatkan daftar transaksi berdasarkan user_id
func (h *TransactionHandler) GetTransactionList(c *gin.Context) {
	userID := c.Param("user_id")
	req := &walletpb.TransactionListRequest{UserId: userID}
	res, err := h.walletClient.GetTransactionList(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
