package query_test

import (
	"database/sql"
	"notification-service/database/query"
	"notification-service/notification"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestNotificationsStmts_InsertNotification(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	notification := query.Notification{
		Identifier: "test",
		Subject:    "Test Subject",
		Message:    "Test Message",
		Type:       notification.NotificationType(1),
		Attempts:   0,
		CreatedAt:  0,
	}

	mock.ExpectPrepare("INSERT INTO notifications").ExpectExec().
		WithArgs(notification.Identifier, notification.Subject, notification.Message, notification.Type, notification.Attempts, notification.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	notificationsStmts, err := query.NewNotificationsStmts(sqlxDB)
	if err != nil {
		t.Fatalf("Failed to create NotificationsStmts: %v", err)
	}

	err = notificationsStmts.InsertNotification(notification)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestNotificationsStmts_GetPendingNotificationsByType(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectPrepare("INSERT INTO notifications")

	notificationsStmts, err := query.NewNotificationsStmts(sqlxDB)
	if err != nil {
		t.Fatalf("Failed to create NotificationsStmts: %v", err)
	}

	nType := notification.NotificationType(1)
	attempts := int64(3)

	expectedNotifications := []*query.Notification{
		{
			ID:         1,
			Identifier: "test",
			Subject:    "Test Subject",
			Message:    "Test Message",
			Type:       nType,
			Attempts:   2,
			CreatedAt:  0,
			SentAt:     sql.NullInt64{},
		},
	}

	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT (.+) FROM notifications WHERE (.+) FOR UPDATE").
		ExpectQuery().
		WithArgs(nType, attempts).
		WillReturnRows(sqlmock.NewRows([]string{"id", "identifier", "subject", "message", "type", "attempts", "created_at", "sent_at"}).
			AddRow(expectedNotifications[0].ID, expectedNotifications[0].Identifier, expectedNotifications[0].Subject, expectedNotifications[0].Message, expectedNotifications[0].Type, expectedNotifications[0].Attempts, expectedNotifications[0].CreatedAt, expectedNotifications[0].SentAt))

	tx, notifications, err := notificationsStmts.GetPendingNotificationsByType(nType, attempts)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(notifications, expectedNotifications) {
		t.Errorf("Expected notifications to be %v, got %v", expectedNotifications, notifications)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}

	tx.Rollback()
}

func TestNotificationsStmts_UpdateNotifications(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectPrepare("INSERT INTO notifications")

	notificationsStmts, err := query.NewNotificationsStmts(sqlxDB)
	if err != nil {
		t.Fatalf("Failed to create NotificationsStmts: %v", err)
	}

	notifications := []*query.Notification{
		{
			ID:         1,
			Identifier: "test",
			Subject:    "Test Subject",
			Message:    "Test Message",
			Type:       notification.NotificationType(1),
			Attempts:   2,
			CreatedAt:  0,
			SentAt:     sql.NullInt64{},
		},
	}

	mock.ExpectBegin()

	mock.ExpectPrepare("UPDATE notifications SET (.+)").
		ExpectExec().
		WithArgs(notifications[0].Attempts, notifications[0].SentAt, notifications[0].ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	tx, err := sqlxDB.Beginx()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	err = notificationsStmts.UpdateNotifications(tx, notifications)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}

	tx.Rollback()
}
