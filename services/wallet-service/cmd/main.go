package main

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/afifn11/gopay-x/services/wallet-service/config"
	"github.com/afifn11/gopay-x/services/wallet-service/internal/domain"
	"github.com/afifn11/gopay-x/services/wallet-service/internal/handler"
	"github.com/afifn11/gopay-x/services/wallet-service/internal/repository"
	"github.com/afifn11/gopay-x/services/wallet-service/internal/usecase"
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
	if err := db.AutoMigrate(&domain.Wallet{}, &domain.WalletTransaction{}); err != nil {
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
	walletRepo := repository.NewWalletRepository(db)
	txRepo := repository.NewWalletTransactionRepository(db)
	lockRepo := repository.NewLockRepository(rdb)
	walletUC := usecase.NewWalletUsecase(walletRepo, txRepo, lockRepo)
	walletHandler := handler.NewWalletHandler(walletUC)

	// Router
	r := handler.NewRouter(walletHandler, cfg)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("🚀 %s running on %s", cfg.App.Name, addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}