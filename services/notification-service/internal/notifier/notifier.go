package notifier

import (
	"fmt"
	"log"
)

// NotificationPayload adalah payload yang dikirim ke user
type NotificationPayload struct {
	UserID  string
	Email   string
	Type    string
	Title   string
	Message string
}

// Notifier interface — bisa di-extend ke email/push/SMS
type Notifier interface {
	Send(payload *NotificationPayload) error
}

// MockNotifier — simulasi pengiriman notifikasi
type MockNotifier struct{}

func NewMockNotifier() Notifier {
	return &MockNotifier{}
}

func (n *MockNotifier) Send(payload *NotificationPayload) error {
	log.Printf(
		"📧 [MOCK NOTIFICATION] To: %s | Type: %s | Title: %s | Message: %s",
		payload.Email,
		payload.Type,
		payload.Title,
		payload.Message,
	)
	fmt.Printf("✉️  Notification sent to user %s\n", payload.UserID)
	return nil
}