package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/afifn11/gopay-x/services/payment-service/config"
	"github.com/afifn11/gopay-x/services/payment-service/internal/middleware"
)

func NewRouter(paymentHandler *PaymentHandler, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "payment-service"})
	})

	// Webhook callback dari payment gateway (no auth)
	r.POST("/api/v1/payments/callback", paymentHandler.HandleCallback)

	v1 := r.Group("/api/v1/payments")
	v1.Use(middleware.AuthMiddleware(cfg))
	{
		v1.POST("/transfer", paymentHandler.Transfer)
		v1.POST("/topup", paymentHandler.TopUpViaGateway)
		v1.GET("/:id", paymentHandler.GetPayment)
		v1.GET("", paymentHandler.GetHistory)
	}

	return r
}