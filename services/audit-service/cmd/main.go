package main

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/afifn11/gopay-x/services/audit-service/config"
	"github.com/afifn11/gopay-x/services/audit-service/internal/domain"
	"github.com/afifn11/gopay-x/services/audit-service/internal/handler"
	kafkaconsumer "github.com/afifn11/gopay-x/services/audit-service/internal/kafka"
	"github.com/afifn11/gopay-x/services/audit-service/internal/repository"
	"github.com/afifn11/gopay-x/services/audit-service/internal/usecase"
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

	if err := db.AutoMigrate(&domain.AuditLog{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	log.Println("✅ Database migrated")

	auditRepo := repository.NewAuditLogRepository(db)
	auditUC := usecase.NewAuditUsecase(auditRepo)
	auditHandler := handler.NewAuditHandler(auditUC)

	// Subscribe ke semua topic untuk full audit trail
	topics := []string{
		"payment.created",
		"payment.success",
		"topup.success",
		"transfer.success",
		"user.registered",
		"login.new_device",
		"fraud.flagged",
	}

	consumer := kafkaconsumer.NewConsumer(
		cfg.Kafka.Brokers,
		topics,
		"audit-service-group",
		auditUC,
	)
	go consumer.Start(context.Background())
	defer consumer.Close()
	log.Println("✅ Kafka consumers started")

	r := handler.NewRouter(auditHandler, cfg)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("🚀 %s running on %s", cfg.App.Name, addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start: %v", err)
	}
}