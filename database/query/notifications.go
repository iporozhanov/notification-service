package query

import (
	"database/sql"
	"notification-service/notification"

	"github.com/jmoiron/sqlx"
)

type Notification struct {
	ID         int64                         `db:"id"`
	Identifier string                        `db:"identifier"`
	Subject    string                        `db:"subject"`
	Message    string                        `db:"message"`
	Type       notification.NotificationType `db:"type"`
	Attempts   int64                         `db:"attempts"`
	CreatedAt  int64                         `db:"created_at"`
	SentAt     sql.NullInt64                 `db:"sent_at"`
}

type NotificationsStmts struct {
	db                 *sqlx.DB
	insertNotification *sqlx.NamedStmt
}

func NewNotificationsStmts(db *sqlx.DB) (*NotificationsStmts, error) {
	ns := &NotificationsStmts{
		db: db,
	}
	if err := ns.prepare(); err != nil {
		return nil, err
	}

	return ns, nil
}

// prepare prepares the SQL statements used by NotificationsStmts.
func (n *NotificationsStmts) prepare() (err error) {
	if n.insertNotification, err = n.db.PrepareNamed(`
		INSERT INTO notifications (identifier, subject, message, type, attempts, created_at)
		VALUES (:identifier, :subject, :message, :type, :attempts, :created_at)
	`); err != nil {
		return
	}

	return
}

// InsertNotification inserts a new notification into the database.
func (n *NotificationsStmts) InsertNotification(notification Notification) (err error) {
	_, err = n.insertNotification.Exec(notification)
	return
}

// GetPendingNotificationsByType retrieves pending notifications of a specific type from the database.
func (n *NotificationsStmts) GetPendingNotificationsByType(nType notification.NotificationType, attempts int64) (*sqlx.Tx, []*Notification, error) {
	var notifications []*Notification

	tx, err := n.db.Beginx()
	if err != nil {
		return nil, nil, err
	}

	stmt, err := tx.PrepareNamed(`
		SELECT 
			id, identifier, subject, message, type, attempts, created_at, sent_at 
		FROM notifications
		WHERE type = :type 
			AND attempts < :attempts 
			AND sent_at IS NULL
		FOR UPDATE
	`)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	err = stmt.Select(&notifications, Notification{
		Type:     nType,
		Attempts: attempts,
	})

	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return nil, nil, err
	}

	return tx, notifications, nil
}

// UpdateNotifications updates the attempts and sent_at fields of multiple notifications in the database.
func (n *NotificationsStmts) UpdateNotifications(tx *sqlx.Tx, notifications []*Notification) error {
	stmt, err := tx.PrepareNamed(`
		UPDATE notifications
		SET attempts = :attempts, sent_at = :sent_at
		WHERE id = :id
	`)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, notification := range notifications {
		_, err = stmt.Exec(notification)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
