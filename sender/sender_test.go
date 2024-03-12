package sender_test

import (
	"fmt"
	"testing"

	"notification-service/config"
	"notification-service/sender"

	"github.com/stretchr/testify/assert"
)

func TestNewSender(t *testing.T) {
	cfg := config.Config{
		Mailgun: config.Mailgun{
			From:       "test@example.com",
			Domain:     "example.com",
			PrivateKey: "mailgunPrivateKey",
		},
		Twilio: config.Twilio{
			Username: "twilioUsername",
			Password: "twilioPassword",
			From:     "twilioFrom",
		},
		Sinch: config.Sinch{
			From:          "sinchFrom",
			ServicePlanID: "sinchServicePlanID",
			APIKey:        "sinchAPIKey",
		},
		Slack: config.Slack{
			Token: "slackToken",
		},
	}

	tests := []struct {
		name       string
		senderType sender.SenderType
		sender     sender.Sender
	}{
		{
			name:       "SenderTypeMainguin",
			senderType: sender.SenderTypeMainguin,
			sender:     &sender.MailgunClient{},
		},
		{
			name:       "SenderTypeSinch",
			senderType: sender.SenderTypeSinch,
			sender:     &sender.SinchClient{},
		},
		{
			name:       "SenderTypeSlack",
			senderType: sender.SenderTypeSlack,
			sender:     &sender.SlackClient{},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Test NewSender with type %s", tt.name), func(t *testing.T) {
			sender, err := sender.NewSender(tt.senderType, cfg)
			assert.NoError(t, err)
			assert.NotNil(t, sender)
			assert.IsType(t, tt.sender, sender)
		})
	}

	t.Run("Test NewSender with unsupported sender type", func(t *testing.T) {
		sender, err := sender.NewSender("unsupportedSenderType", cfg)
		assert.Error(t, err)
		assert.Nil(t, sender)
	})
}
