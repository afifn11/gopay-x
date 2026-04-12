package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/afifn11/gopay-x/services/fraud-detection-service/internal/usecase"
)

type PaymentEvent struct {
	EventType   string `json:"event_type"`
	UserID      string `json:"user_id"`
	ReferenceID string `json:"reference_id"`
	Amount      int64  `json:"amount"`
}

type Consumer struct {
	reader  *kafka.Reader
	fraudUC usecase.FraudUsecase
}

func NewConsumer(brokers, topic, groupID string, fraudUC usecase.FraudUsecase) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokers},
		Topic:   topic,
		GroupID: groupID,
	})
	return &Consumer{reader: reader, fraudUC: fraudUC}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Printf("📨 Fraud consumer listening on: %s", c.reader.Config().Topic)

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("Kafka read error: %v", err)
			continue
		}

		var event PaymentEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to unmarshal: %v", err)
			continue
		}

		log.Printf("🔍 Checking fraud: ref=%s amount=%d", event.ReferenceID, event.Amount)

		check, err := c.fraudUC.CheckTransaction(ctx, &usecase.CheckRequest{
			UserID:      event.UserID,
			ReferenceID: event.ReferenceID,
			EventType:   event.EventType,
			Amount:      event.Amount,
		})
		if err != nil {
			log.Printf("Fraud check error: %v", err)
			continue
		}

		log.Printf("✅ Fraud result: ref=%s score=%d level=%s status=%s",
			event.ReferenceID, check.RiskScore, check.RiskLevel, check.Status)
	}
}

func (c *Consumer) Close() {
	c.reader.Close()
}