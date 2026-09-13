package postgres

import (
	"context"
	"errors"
)

// ErrNotImplemented indicates that database operations are scheduled for Phase 4.
var ErrNotImplemented = errors.New("postgres repository: implementation scheduled for Phase 4")

// Client wraps the PostgreSQL connection pool (to be initialized with pgx in Phase 4).
type Client struct {
	dsn string
}

// NewClient constructs a PostgreSQL client stub for Phase 0.
func NewClient(dsn string) *Client {
	return &Client{dsn: dsn}
}

// Ping provides a connectivity check stub.
func (c *Client) Ping(ctx context.Context) error {
	// Phase 0: Return nil; full pgx pool ping configured in Phase 4.
	return nil
}
