package sender

import (
	"fmt"
	"notification-service/config"
	"notification-service/notification"
)

// Sender defines the interface for sending notifications.
type Sender interface {
	SendNotification(notification notification.Notification) error
}

type Senders map[notification.NotificationType]Sender

// NewSender creates a new instance of a sender based on the provided type.
func NewSender(senderType SenderType, cfg config.Config) (Sender, error) {
	switch senderType {
	case SenderTypeMainguin:
		return NewMailgunClient(cfg.Mailgun.From, cfg.Mailgun.Domain, cfg.Mailgun.PrivateKey)
	case SenderTypeSinch:
		return NewSinchClient(cfg.Sinch.From, cfg.Sinch.ServicePlanID, cfg.Sinch.APIKey)
	case SenderTypeSlack:
		return NewSlackClient(cfg.Slack.Token)
	}
	return nil, fmt.Errorf("sender type %s not supported", senderType)

}

type SenderType string

const (
	SenderTypeMainguin SenderType = "mailgun"
	SenderTypeSlack    SenderType = "slack"
	SenderTypeSinch    SenderType = "sinch"
)

var AvailableSubTypes = []SenderType{
	SenderTypeMainguin,
	SenderTypeSlack,
	SenderTypeSinch,
}

func (n SenderType) IsValid() bool {
	for _, v := range AvailableSubTypes {
		if n == v {
			return true
		}
	}
	return false
}
