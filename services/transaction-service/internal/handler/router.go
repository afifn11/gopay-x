package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/afifn11/gopay-x/services/transaction-service/config"
	"github.com/afifn11/gopay-x/services/transaction-service/internal/middleware"
)

func NewRouter(txHandler *TransactionHandler, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "transaction-service"})
	})

	v1 := r.Group("/api/v1/transactions")
	v1.Use(middleware.AuthMiddleware(cfg))
	{
		v1.GET("", txHandler.GetHistory)
		v1.GET("/summary", txHandler.GetSummary)
		v1.GET("/:id", txHandler.GetTransaction)
	}

	return r
}