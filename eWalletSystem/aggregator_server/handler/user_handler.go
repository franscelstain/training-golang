package handler

import (
	"context"
	walletpb "eWalletSystem/proto/wallet/v1"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserHandler struct {
	userClient walletpb.UserServiceClient
}

// Membuat instance baru dari UserHandler
func NewUserHandler() *UserHandler {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Failed to connect to User Server")
	}
	client := walletpb.NewUserServiceClient(conn)
	return &UserHandler{
		userClient: client,
	}
}

// Mengambil data user dari user server
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("user_id")
	req := &walletpb.UserRequest{UserId: userID}
	res, err := h.userClient.GetUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
