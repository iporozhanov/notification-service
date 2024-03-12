package database

import (
	"notification-service/database/query"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// DB is a struct that holds the database connection and prepared statements for notifications.
type DB struct {
	db *sqlx.DB
	*query.NotificationsStmts
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.db.Close()
}

// NewDB creates a new instance of the database with the given DSN (Data Source Name).
func NewDB(dsn string) (*DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver,
	)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	NotificationsStmts, err := query.NewNotificationsStmts(db)
	if err != nil {
		return nil, err
	}

	return &DB{
		db:                 db,
		NotificationsStmts: NotificationsStmts,
	}, nil
}
