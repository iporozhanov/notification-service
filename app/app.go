package app

import (
	"database/sql"
	"notification-service/config"
	"notification-service/database/query"
	"notification-service/notification"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Sender defines the interface for sending notifications.
type Sender interface {
	SendNotification(notification notification.Notification) error
}

// Senders is a map of notification types to their corresponding senders.
type Senders map[notification.NotificationType]Sender

// DB defines the interface for interacting with the database.
type DB interface {
	InsertNotification(notification query.Notification) (err error)
	GetPendingNotificationsByType(nType notification.NotificationType, attempts int64) (*sqlx.Tx, []*query.Notification, error)
	UpdateNotifications(tx *sqlx.Tx, notifications []*query.Notification) error
	Close() error
}

// App represents the main application.
type App struct {
	db                       DB
	Senders                  Senders
	Log                      *zap.SugaredLogger
	NotificationMaxAttempts  int64
	NotificationListenTicker time.Duration
	quitChann                chan struct{}
}

// NewApp creates a new instance of the App.
func NewApp(db DB, log *zap.SugaredLogger, cfg config.Config) *App {
	return &App{
		db:                       db,
		Log:                      log,
		quitChann:                make(chan struct{}),
		NotificationMaxAttempts:  cfg.NotificationMaxAttempts,
		NotificationListenTicker: cfg.NotificationListenTicker,
		Senders:                  make(Senders),
	}
}

func (a *App) AddSender(nType notification.NotificationType, s Sender) {
	a.Senders[nType] = s
	go a.startNotificationPoller(nType)
}

// Shutdown gracefully shuts down the application.
func (a *App) Shutdown() {
	a.Log.Info("shutting down app")
	go func() {
		a.quitChann <- struct{}{}
		close(a.quitChann)
	}()

	err := a.db.Close()
	if err != nil {
		a.Log.Errorf("error closing db: %v", err)
	}
}

// startNotificationPoller starts a poller for the specified notification type.
func (a *App) startNotificationPoller(nType notification.NotificationType) {
	ticker := time.NewTicker(a.NotificationListenTicker)
	for {
		select {
		case <-ticker.C:
			a.SendPendingNotifications(nType)
		case <-a.quitChann:
			return
		}
	}
}

// SendPendingNotifications sends pending notifications of the specified type using the corresponding sender.
func (a *App) SendPendingNotifications(nType notification.NotificationType) {
	tx, notifications, err := a.db.GetPendingNotificationsByType(nType, a.NotificationMaxAttempts)
	if err != nil {
		tx.Rollback()
		a.Log.Errorf("error getting pending notifications for type %s: %s", nType, err)
		return
	}
	if len(notifications) == 0 {
		tx.Rollback()
		return
	}

	for key, n := range notifications {
		if err := a.Senders[n.Type].SendNotification(notification.Notification{
			Identifier: n.Identifier,
			Subject:    n.Subject,
			Message:    n.Message,
			Type:       n.Type,
		}); err != nil {
			a.Log.Errorf("error sending notification %d type: %s: %s", n.ID, nType, err)
			notifications[key].Attempts++
			continue
		}

		notifications[key].SentAt = sql.NullInt64{Int64: time.Now().Unix(), Valid: true}
	}
	if err := a.db.UpdateNotifications(tx, notifications); err != nil {
		tx.Rollback()
		a.Log.Errorf("error updating notifications for type %s: %s", nType, err)
	}
}

// NewNotification creates a new notification and inserts it into the database.
func (a *App) NewNotification(n notification.Notification) error {
	return a.db.InsertNotification(query.Notification{
		Identifier: n.Identifier,
		Subject:    n.Subject,
		Message:    n.Message,
		Type:       n.Type,
		CreatedAt:  time.Now().Unix(),
	})
}
