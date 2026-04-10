package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/afifn11/gopay-x/services/wallet-service/config"
	"github.com/afifn11/gopay-x/services/wallet-service/internal/middleware"
)

func NewRouter(walletHandler *WalletHandler, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "wallet-service"})
	})

	v1 := r.Group("/api/v1/wallets")

	// Internal route
	v1.POST("/internal/create", walletHandler.CreateWalletInternal)

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.POST("", walletHandler.CreateWallet)
		protected.GET("", walletHandler.GetWallet)
		protected.POST("/topup", walletHandler.TopUp)
		protected.GET("/transactions", walletHandler.GetTransactionHistory)
	}

	return r
}