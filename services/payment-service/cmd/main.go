package main

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/afifn11/gopay-x/services/payment-service/config"
	"github.com/afifn11/gopay-x/services/payment-service/internal/domain"
	"github.com/afifn11/gopay-x/services/payment-service/internal/gateway"
	"github.com/afifn11/gopay-x/services/payment-service/internal/handler"
	"github.com/afifn11/gopay-x/services/payment-service/internal/repository"
	"github.com/afifn11/gopay-x/services/payment-service/internal/usecase"
)

func main() {
	cfg := config.Load()

	// Connect PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.DBName, cfg.Database.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("✅ PostgreSQL connected")

	// Auto migrate
	if err := db.AutoMigrate(&domain.Payment{}, &domain.PaymentCallback{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	log.Println("✅ Database migrated")

	// Connect Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	log.Println("✅ Redis connected")

	// Wire up layers
	paymentRepo := repository.NewPaymentRepository(db)
	callbackRepo := repository.NewPaymentCallbackRepository(db)
	lockRepo := repository.NewLockRepository(rdb)
	gw := gateway.NewMockGateway()
	paymentUC := usecase.NewPaymentUsecase(paymentRepo, callbackRepo, lockRepo, gw)
	paymentHandler := handler.NewPaymentHandler(paymentUC)

	// Router
	r := handler.NewRouter(paymentHandler, cfg)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("🚀 %s running on %s", cfg.App.Name, addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}