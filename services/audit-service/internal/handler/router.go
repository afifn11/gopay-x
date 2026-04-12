package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/afifn11/gopay-x/services/audit-service/config"
	"github.com/afifn11/gopay-x/services/audit-service/internal/middleware"
)

func NewRouter(auditHandler *AuditHandler, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "audit-service"})
	})

	v1 := r.Group("/api/v1/audit")
	v1.Use(middleware.AuthMiddleware(cfg))
	v1.Use(middleware.AdminOnly())
	{
		v1.GET("/logs", auditHandler.QueryLogs)
		v1.GET("/actors/:actor_id", auditHandler.GetByActor)
		v1.GET("/resources/:resource_id", auditHandler.GetByResource)
	}

	return r
}