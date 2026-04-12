package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/afifn11/gopay-x/services/audit-service/internal/domain"
	"github.com/afifn11/gopay-x/services/audit-service/internal/usecase"
)

type AuditEvent struct {
	EventType    string `json:"event_type"`
	ServiceName  string `json:"service_name"`
	ActorID      string `json:"actor_id"`
	ActorType    string `json:"actor_type"`
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	Action       string `json:"action"`
	Status       string `json:"status"`
	IPAddress    string `json:"ip_address"`
	Payload      string `json:"payload"`
	ErrorMsg     string `json:"error_msg"`
}

type Consumer struct {
	readers []*kafka.Reader
	auditUC usecase.AuditUsecase
}

func NewConsumer(brokers string, topics []string, groupID string, auditUC usecase.AuditUsecase) *Consumer {
	var readers []*kafka.Reader
	for _, topic := range topics {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokers},
			Topic:   topic,
			GroupID: groupID,
		})
		readers = append(readers, r)
	}
	return &Consumer{readers: readers, auditUC: auditUC}
}

func (c *Consumer) Start(ctx context.Context) {
	for _, reader := range c.readers {
		go c.consume(ctx, reader)
	}
}

func (c *Consumer) consume(ctx context.Context, reader *kafka.Reader) {
	log.Printf("📋 Audit consumer listening: %s", reader.Config().Topic)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("Kafka read error: %v", err)
			continue
		}

		var event AuditEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to unmarshal audit event: %v", err)
			continue
		}

		auditLog := &domain.AuditLog{
			ServiceName:  event.ServiceName,
			EventType:    event.EventType,
			ActorID:      event.ActorID,
			ActorType:    event.ActorType,
			ResourceID:   event.ResourceID,
			ResourceType: event.ResourceType,
			Action:       event.Action,
			Status:       event.Status,
			IPAddress:    event.IPAddress,
			Payload:      event.Payload,
			ErrorMsg:     event.ErrorMsg,
			CreatedAt:    time.Now(),
		}

		if err := c.auditUC.RecordLog(ctx, auditLog); err != nil {
			log.Printf("Failed to record audit log: %v", err)
			continue
		}

		log.Printf("📋 Audit recorded: %s | %s | actor=%s", event.ServiceName, event.EventType, event.ActorID)
	}
}

func (c *Consumer) Close() {
	for _, r := range c.readers {
		r.Close()
	}
}