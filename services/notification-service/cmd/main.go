package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/afifn11/gopay-x/services/notification-service/config"
	kafkaconsumer "github.com/afifn11/gopay-x/services/notification-service/internal/kafka"
	"github.com/afifn11/gopay-x/services/notification-service/internal/notifier"
)

func main() {
	cfg := config.Load()

	// Setup mock notifier
	n := notifier.NewMockNotifier()

	// Subscribe ke multiple Kafka topics
	topics := []string{
		"payment.success",
		"topup.success",
		"transfer.success",
		"login.new_device",
	}

	consumer := kafkaconsumer.NewConsumer(
		cfg.Kafka.Brokers,
		topics,
		"notification-service-group",
		n,
	)

	go consumer.Start(context.Background())
	defer consumer.Close()

	log.Println("✅ Kafka consumers started")

	// Simple health check server
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "notification-service",
			"topics":  topics,
		})
	})

	log.Printf("🚀 %s running on :%s", cfg.App.Name, cfg.App.Port)
	r.Run(":" + cfg.App.Port)
}