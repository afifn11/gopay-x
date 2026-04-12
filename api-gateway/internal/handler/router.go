package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/afifn11/gopay-x/api-gateway/config"
	"github.com/afifn11/gopay-x/api-gateway/internal/middleware"
	"github.com/afifn11/gopay-x/api-gateway/internal/proxy"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.RateLimiter())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "api-gateway",
			"version": "1.0.0",
		})
	})

	// Services health
	r.GET("/health/services", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"services": gin.H{
				"auth-service":        cfg.Services.AuthService,
				"user-service":        cfg.Services.UserService,
				"wallet-service":      cfg.Services.WalletService,
				"payment-service":     cfg.Services.PaymentService,
				"transaction-service": cfg.Services.TransactionService,
				"notification-service": cfg.Services.NotificationService,
				"fraud-service":       cfg.Services.FraudService,
				"audit-service":       cfg.Services.AuditService,
			},
		})
	})

	// ─── Auth Service (public) ────────────────────────────────
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", proxy.ReverseProxy(cfg.Services.AuthService))
		auth.POST("/login", proxy.ReverseProxy(cfg.Services.AuthService))
		auth.POST("/refresh", proxy.ReverseProxy(cfg.Services.AuthService))
		auth.POST("/logout", proxy.ReverseProxy(cfg.Services.AuthService))
		auth.GET("/validate", proxy.ReverseProxy(cfg.Services.AuthService))
	}

	// ─── Protected Routes ─────────────────────────────────────
	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))

	// User Service
	users := protected.Group("/api/v1/users")
	{
		users.GET("/me", proxy.ReverseProxy(cfg.Services.UserService))
		users.PUT("/me", proxy.ReverseProxy(cfg.Services.UserService))
		users.POST("/me/kyc", proxy.ReverseProxy(cfg.Services.UserService))
	}

	// Wallet Service
	wallets := protected.Group("/api/v1/wallets")
	{
		wallets.POST("", proxy.ReverseProxy(cfg.Services.WalletService))
		wallets.GET("", proxy.ReverseProxy(cfg.Services.WalletService))
		wallets.POST("/topup", proxy.ReverseProxy(cfg.Services.WalletService))
		wallets.GET("/transactions", proxy.ReverseProxy(cfg.Services.WalletService))
	}

	// Payment Service
	payments := protected.Group("/api/v1/payments")
	{
		payments.POST("/transfer", proxy.ReverseProxy(cfg.Services.PaymentService))
		payments.POST("/topup", proxy.ReverseProxy(cfg.Services.PaymentService))
		payments.GET("/:id", proxy.ReverseProxy(cfg.Services.PaymentService))
		payments.GET("", proxy.ReverseProxy(cfg.Services.PaymentService))
	}
	// Payment callback (public — dari payment gateway)
	r.POST("/api/v1/payments/callback", proxy.ReverseProxy(cfg.Services.PaymentService))

	// Transaction Service
	transactions := protected.Group("/api/v1/transactions")
	{
		transactions.GET("", proxy.ReverseProxy(cfg.Services.TransactionService))
		transactions.GET("/summary", proxy.ReverseProxy(cfg.Services.TransactionService))
		transactions.GET("/:id", proxy.ReverseProxy(cfg.Services.TransactionService))
	}

	// ─── Admin Routes ─────────────────────────────────────────
	admin := protected.Group("")
	admin.Use(middleware.AdminOnly())

	// Fraud Service (admin only)
	fraud := admin.Group("/api/v1/fraud")
	{
		fraud.GET("/users/:user_id/checks", proxy.ReverseProxy(cfg.Services.FraudService))
		fraud.GET("/users/:user_id/risk-profile", proxy.ReverseProxy(cfg.Services.FraudService))
	}

	// Audit Service (admin only)
	audit := admin.Group("/api/v1/audit")
	{
		audit.GET("/logs", proxy.ReverseProxy(cfg.Services.AuditService))
		audit.GET("/actors/:actor_id", proxy.ReverseProxy(cfg.Services.AuditService))
		audit.GET("/resources/:resource_id", proxy.ReverseProxy(cfg.Services.AuditService))
	}

	// User management (admin only)
	adminUsers := admin.Group("/api/v1/users")
	{
		adminUsers.GET("/:user_id", proxy.ReverseProxy(cfg.Services.UserService))
		adminUsers.PUT("/:user_id/kyc-status", proxy.ReverseProxy(cfg.Services.UserService))
	}

	return r
}