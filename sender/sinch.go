package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"notification-service/notification"
)

type SinchClient struct {
	from     string
	planID   string
	apiToken string
}

func NewSinchClient(from, planID, apiToken string) (*SinchClient, error) {
	if from == "" || planID == "" || apiToken == "" {
		return nil, fmt.Errorf("missing required sinch configuration")
	}

	return &SinchClient{from, planID, apiToken}, nil
}

type SMSRequest struct {
	From string   `json:"from"`
	To   []string `json:"to"`
	Body string   `json:"body"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Text    string `json:"text"`
}

func (s *SinchClient) SendNotification(notification notification.Notification) error {
	apiUrl := fmt.Sprintf("https://us.sms.api.sinch.com/xms/v1/%s/batches", s.planID)
	payload := SMSRequest{
		From: s.from,
		To:   []string{notification.Identifier},
		Body: notification.Message,
	}

	bPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %w", err)
	}

	request, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(bPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiToken))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error sending SMS message: %w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	jsonErr := ErrorResponse{}
	if resp.StatusCode == http.StatusUnauthorized { // 401 has no body
		jsonErr.Message = "Unauthorized"
		jsonErr.Text = "Invalid API token"
	} else {
		responseBody, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}

		err = json.Unmarshal(responseBody, &jsonErr)
		if err != nil {
			return fmt.Errorf("error unmarshalling response body: %w", err)
		}
	}

	return fmt.Errorf("error sending SMS message: %v", jsonErr)

}
