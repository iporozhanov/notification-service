package app_test

import (
	"database/sql"
	"notification-service/app"
	"notification-service/app/mocks"
	"notification-service/config"
	"notification-service/database/query"
	"notification-service/notification"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func TestApp_AddSender(t *testing.T) {
	mockDB := mocks.NewDB(t)
	mockSender := mocks.NewSender(t)
	a := app.NewApp(mockDB, zap.NewNop().Sugar(), config.Config{
		NotificationMaxAttempts:  3,
		NotificationListenTicker: 1 * time.Minute,
	})

	nType := notification.NotificationType(1)

	a.AddSender(nType, mockSender)

	// Verify that the sender was added correctly
	if _, ok := a.Senders[nType]; !ok {
		t.Errorf("expected sender to be added for notification type %s", nType)
	}
	mockDB.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

func TestApp_SendPendingNotifications(t *testing.T) {
	mockDB := mocks.NewDB(t)
	mockSender := mocks.NewSender(t)
	a := app.NewApp(mockDB, zap.NewNop().Sugar(), config.Config{
		NotificationMaxAttempts:  3,
		NotificationListenTicker: 1,
	})

	nType := notification.NotificationType(1)
	a.AddSender(nType, mockSender)

	// Set up mock data
	notifications := []*query.Notification{
		{
			ID:         1,
			Identifier: "test",
			Subject:    "Test Subject",
			Message:    "Test Message",
			Type:       nType,
			Attempts:   0,
			CreatedAt:  0,
			SentAt:     sql.NullInt64{Int64: 2, Valid: true},
		},
	}

	mockDB.EXPECT().GetPendingNotificationsByType(nType, int64(3)).Return(&sqlx.Tx{}, notifications, nil)

	mockSender.EXPECT().SendNotification(notification.Notification{
		Identifier: "test",
		Subject:    "Test Subject",
		Message:    "Test Message",
		Type:       nType,
	}).Return(nil)
	mockDB.EXPECT().UpdateNotifications(&sqlx.Tx{}, notifications).Return(nil)

	a.SendPendingNotifications(nType)

	mockDB.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}

func TestApp_NewNotification(t *testing.T) {
	mockDB := mocks.NewDB(t)
	mockSender := mocks.NewSender(t)
	a := app.NewApp(mockDB, zap.NewNop().Sugar(), config.Config{
		NotificationMaxAttempts:  3,
		NotificationListenTicker: 1,
	})

	nType := notification.NotificationType(1)

	notification := notification.Notification{
		Identifier: "test",
		Subject:    "Test Subject",
		Message:    "Test Message",
		Type:       nType,
	}

	mockDB.EXPECT().InsertNotification(query.Notification{
		Identifier: "test",
		Subject:    "Test Subject",
		Message:    "Test Message",
		Type:       nType,
		CreatedAt:  time.Now().Unix(),
	}).Return(nil)

	err := a.NewNotification(notification)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockDB.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}
