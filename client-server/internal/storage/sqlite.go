package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/paraizofelipe/fullcycle-client-server/internal/ctxlog"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Init(parent context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS exchange_rates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT NOT NULL,
		created_at DATETIME NOT NULL)`,
	)
	if err != nil {
		ctxlog.LogDeadline(ctx, err, "db init")
		return err
	}
	return nil
}

func (s *Store) SaveBid(parent context.Context, bid string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, "INSERT INTO exchange_rates (bid, created_at) VALUES (?, ?)", bid, time.Now())
	if err != nil {
		ctxlog.LogDeadline(ctx, err, "db insert")
		return err
	}
	return nil
}
