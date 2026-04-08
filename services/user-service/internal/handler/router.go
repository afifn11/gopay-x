package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/afifn11/gopay-x/services/user-service/internal/middleware"
	"github.com/afifn11/gopay-x/services/user-service/config"
)

func NewRouter(userHandler *UserHandler, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "user-service"})
	})

	v1 := r.Group("/api/v1/users")

	// Internal route — dipanggil oleh auth-service saat register (tanpa auth)
	v1.POST("/internal/create", userHandler.CreateProfile)

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/me", userHandler.GetMyProfile)
		protected.PUT("/me", userHandler.UpdateProfile)
		protected.POST("/me/kyc", userHandler.SubmitKYC)

		// Admin only
		admin := protected.Group("")
		admin.Use(middleware.AdminOnly())
		{
			admin.GET("/:user_id", userHandler.GetProfileByID)
			admin.PUT("/:user_id/kyc-status", userHandler.UpdateKYCStatus)
		}
	}

	return r
}