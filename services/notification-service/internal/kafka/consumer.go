package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/afifn11/gopay-x/services/notification-service/internal/notifier"
)

type Event struct {
	EventType   string `json:"event_type"`
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	Amount      int64  `json:"amount"`
	ReferenceID string `json:"reference_id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type Consumer struct {
	readers  []*kafka.Reader
	notifier notifier.Notifier
}

func NewConsumer(brokers string, topics []string, groupID string, n notifier.Notifier) *Consumer {
	var readers []*kafka.Reader
	for _, topic := range topics {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokers},
			Topic:   topic,
			GroupID: groupID,
		})
		readers = append(readers, r)
	}
	return &Consumer{readers: readers, notifier: n}
}

func (c *Consumer) Start(ctx context.Context) {
	for _, reader := range c.readers {
		go c.consume(ctx, reader)
	}
}

func (c *Consumer) consume(ctx context.Context, reader *kafka.Reader) {
	log.Printf("📨 Listening on topic: %s", reader.Config().Topic)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("Kafka read error: %v", err)
			continue
		}

		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			continue
		}

		log.Printf("📩 Event received: %s | user: %s", event.EventType, event.UserID)
		c.handleEvent(event)
	}
}

func (c *Consumer) handleEvent(event Event) {
	var payload *notifier.NotificationPayload

	switch event.EventType {
	case "payment.success":
		payload = &notifier.NotificationPayload{
			UserID:  event.UserID,
			Email:   event.Email,
			Type:    "payment_success",
			Title:   "Pembayaran Berhasil",
			Message: fmt.Sprintf("Transaksi Rp%d berhasil diproses. Ref: %s", event.Amount, event.ReferenceID),
		}
	case "topup.success":
		payload = &notifier.NotificationPayload{
			UserID:  event.UserID,
			Email:   event.Email,
			Type:    "topup_success",
			Title:   "Top Up Berhasil",
			Message: fmt.Sprintf("Saldo kamu berhasil ditambah Rp%d.", event.Amount),
		}
	case "transfer.success":
		payload = &notifier.NotificationPayload{
			UserID:  event.UserID,
			Email:   event.Email,
			Type:    "transfer_success",
			Title:   "Transfer Berhasil",
			Message: fmt.Sprintf("Transfer Rp%d berhasil. Ref: %s", event.Amount, event.ReferenceID),
		}
	case "login.new_device":
		payload = &notifier.NotificationPayload{
			UserID:  event.UserID,
			Email:   event.Email,
			Type:    "security_alert",
			Title:   "Login Perangkat Baru",
			Message: "Terdeteksi login dari perangkat baru. Jika bukan kamu, segera ganti password.",
		}
	default:
		log.Printf("Unknown event type: %s", event.EventType)
		return
	}

	if err := c.notifier.Send(payload); err != nil {
		log.Printf("Failed to send notification: %v", err)
	}
}

func (c *Consumer) Close() {
	for _, r := range c.readers {
		r.Close()
	}
}