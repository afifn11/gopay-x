package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/afifn11/gopay-x/services/fraud-detection-service/config"
	"github.com/afifn11/gopay-x/services/fraud-detection-service/internal/domain"
	"github.com/afifn11/gopay-x/services/fraud-detection-service/internal/handler"
	kafkaconsumer "github.com/afifn11/gopay-x/services/fraud-detection-service/internal/kafka"
	"github.com/afifn11/gopay-x/services/fraud-detection-service/internal/repository"
	"github.com/afifn11/gopay-x/services/fraud-detection-service/internal/rules"
	"github.com/afifn11/gopay-x/services/fraud-detection-service/internal/usecase"
)

func main() {
	cfg := config.Load()

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
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Println("✅ PostgreSQL connected")

	if err := db.AutoMigrate(&domain.FraudCheck{}, &domain.UserRiskProfile{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	log.Println("✅ Database migrated")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	log.Println("✅ Redis connected")
	_ = rdb

	fraudRepo := repository.NewFraudCheckRepository(db)
	profileRepo := repository.NewUserRiskProfileRepository(db)
	engine := rules.NewRuleEngine(fraudRepo)
	fraudUC := usecase.NewFraudUsecase(fraudRepo, profileRepo, engine)
	fraudHandler := handler.NewFraudHandler(fraudUC)

	consumer := kafkaconsumer.NewConsumer(
		cfg.Kafka.Brokers,
		"payment.created",
		"fraud-detection-group",
		fraudUC,
	)
	go consumer.Start(context.Background())
	defer consumer.Close()
	log.Println("✅ Kafka consumer started")

	r := handler.NewRouter(fraudHandler, cfg)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("🚀 %s running on %s", cfg.App.Name, addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start: %v", err)
	}
}