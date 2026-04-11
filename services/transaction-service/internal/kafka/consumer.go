package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/afifn11/gopay-x/services/transaction-service/internal/usecase"
)

type TransactionEvent struct {
	EventType     string `json:"event_type"`
	UserID        string `json:"user_id"`
	CounterpartID string `json:"counterpart_id,omitempty"`
	TxType        string `json:"tx_type"`
	Amount        int64  `json:"amount"`
	Fee           int64  `json:"fee"`
	BalanceBefore int64  `json:"balance_before"`
	BalanceAfter  int64  `json:"balance_after"`
	ReferenceID   string `json:"reference_id"`
	ServiceSource string `json:"service_source"`
	Description   string `json:"description"`
	Status        string `json:"status"`
}

type Consumer struct {
	reader *kafka.Reader
	txUC   usecase.TransactionUsecase
}

func NewConsumer(brokers, topic, groupID string, txUC usecase.TransactionUsecase) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokers},
		Topic:   topic,
		GroupID: groupID,
	})

	return &Consumer{reader: reader, txUC: txUC}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Printf("📨 Kafka consumer started, listening on topic: %s", c.reader.Config().Topic)

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Kafka consumer stopped")
				return
			}
			log.Printf("Kafka read error: %v", err)
			continue
		}

		var event TransactionEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			continue
		}

		log.Printf("📩 Received event: %s | ref: %s", event.EventType, event.ReferenceID)

		req := &usecase.RecordTransactionRequest{
			UserID:        event.UserID,
			CounterpartID: event.CounterpartID,
			Type:          event.TxType,
			Amount:        event.Amount,
			Fee:           event.Fee,
			BalanceBefore: event.BalanceBefore,
			BalanceAfter:  event.BalanceAfter,
			ReferenceID:   event.ReferenceID,
			ServiceSource: event.ServiceSource,
			Description:   event.Description,
			Status:        event.Status,
		}

		if _, err := c.txUC.RecordTransaction(ctx, req); err != nil {
			log.Printf("Failed to record transaction: %v", err)
		}
	}
}

func (c *Consumer) Close() {
	c.reader.Close()
}