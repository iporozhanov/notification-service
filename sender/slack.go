package sender

import (
	"fmt"
	"notification-service/notification"

	"github.com/slack-go/slack"
)

type SlackClient struct {
	client *slack.Client
}

func NewSlackClient(apiToken string) (*SlackClient, error) {
	if apiToken == "" {
		return nil, fmt.Errorf("missing required slack configuration")
	}

	client := slack.New(apiToken)
	return &SlackClient{
		client: client,
	}, nil
}

func (s *SlackClient) SendNotification(notification notification.Notification) error {
	_, _, err := s.client.PostMessage(notification.Identifier, slack.MsgOptionText(notification.Message, false), slack.MsgOptionAsUser(true))
	if err != nil {
		return fmt.Errorf("error sending slack message: %w", err)
	}

	return nil
}
