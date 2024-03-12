package sender

import (
	"context"
	"fmt"
	"notification-service/notification"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailgunClient struct {
	from   string
	client *mailgun.MailgunImpl
}

func NewMailgunClient(from, domain, privateAPIKey string) (*MailgunClient, error) {
	if from == "" || domain == "" || privateAPIKey == "" {
		return nil, fmt.Errorf("missing required mailgun configuration")
	}

	mg := mailgun.NewMailgun(domain, privateAPIKey)

	return &MailgunClient{from, mg}, nil
}

func (m *MailgunClient) SendNotification(notification notification.Notification) error {
	message := m.client.NewMessage(
		m.from,
		notification.Subject,
		notification.Message,
		notification.Identifier,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := m.client.Send(ctx, message)

	return err
}
