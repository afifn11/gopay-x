package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/afifn11/gopay-x/services/user-service/config"
	"github.com/afifn11/gopay-x/services/user-service/internal/domain"
	"github.com/afifn11/gopay-x/services/user-service/internal/handler"
	"github.com/afifn11/gopay-x/services/user-service/internal/repository"
	"github.com/afifn11/gopay-x/services/user-service/internal/usecase"
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
	if err := db.AutoMigrate(&domain.UserProfile{}, &domain.KYCDocument{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	log.Println("✅ Database migrated")

	// Wire up layers
	profileRepo := repository.NewUserProfileRepository(db)
	kycRepo := repository.NewKYCDocumentRepository(db)
	userUC := usecase.NewUserUsecase(profileRepo, kycRepo)
	userHandler := handler.NewUserHandler(userUC)

	// Router
	r := handler.NewRouter(userHandler, cfg)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("🚀 %s running on %s", cfg.App.Name, addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}