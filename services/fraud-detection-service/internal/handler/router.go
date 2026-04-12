package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/afifn11/gopay-x/services/fraud-detection-service/config"
	"github.com/afifn11/gopay-x/services/fraud-detection-service/internal/middleware"
)

func NewRouter(fraudHandler *FraudHandler, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "fraud-detection-service"})
	})

	v1 := r.Group("/api/v1/fraud")
	v1.Use(middleware.AuthMiddleware(cfg))
	v1.Use(middleware.AdminOnly())
	{
		v1.GET("/users/:user_id/checks", fraudHandler.GetChecks)
		v1.GET("/users/:user_id/risk-profile", fraudHandler.GetRiskProfile)
	}

	return r
}