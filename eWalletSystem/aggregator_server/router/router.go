package router

import (
	"eWalletSystem/aggregator_server/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	userHandler := handler.NewUserHandler()
	transactionHandler := handler.NewTransactionHandler()

	// Route untuk mendapatkan data user berdasarkan user_id
	router.GET("/user/:user_id", userHandler.GetUser)

	// Route untuk melakukan top-up saldo
	router.POST("/wallet/topup", transactionHandler.TopUp)

	// Route untuk melakukan transfer saldo antar pengguna
	router.POST("/wallet/transfer", transactionHandler.Transfer)

	// Route untuk mendapatkan daftar transaksi user berdasarkan user_id
	router.GET("/wallet/transactions/:user_id", transactionHandler.GetTransactionList)

	return router
}
