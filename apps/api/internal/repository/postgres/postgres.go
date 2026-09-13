package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nachsyas/arham-porto/apps/api/internal/config"
)

// ErrDatabaseDisabled is returned when attempting DB operations while DB mode is disabled.
var ErrDatabaseDisabled = errors.New("postgres: database mode is disabled")

// Client wraps the PostgreSQL connection pool using pgxpool.
type Client struct {
	dsn   string
	mode  string
	pool  *pgxpool.Pool
	mu    sync.RWMutex
	ready bool
}

// NewClient constructs a PostgreSQL client for the configured database mode.
func NewClient(dsn string, mode string) *Client {
	return &Client{
		dsn:  dsn,
		mode: mode,
	}
}

// Connect attempts to initialize the pgxpool connection pool according to the database mode.
func (c *Client) Connect(ctx context.Context) error {
	if c.mode == config.DatabaseModeDisabled {
		log.Println("[postgres] Database mode is 'disabled'. Skipping connection.")
		return nil
	}

	cfg, err := pgxpool.ParseConfig(c.dsn)
	if err != nil {
		if c.mode == config.DatabaseModeRequired {
			return fmt.Errorf("postgres: invalid DSN: %w", err)
		}
		log.Printf("[postgres] Warning: invalid DSN: %v (continuing in optional mode)", err)
		return nil
	}

	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnLifetime = 1 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	// Short timeout for connection attempt
	connCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connCtx, cfg)
	if err != nil {
		if c.mode == config.DatabaseModeRequired {
			return fmt.Errorf("postgres: connection failed: %w", err)
		}
		log.Printf("[postgres] Warning: connection failed: %v (continuing with degraded status)", err)
		return nil
	}

	// Ping database
	if err := pool.Ping(connCtx); err != nil {
		pool.Close()
		if c.mode == config.DatabaseModeRequired {
			return fmt.Errorf("postgres: ping failed: %w", err)
		}
		log.Printf("[postgres] Warning: ping failed: %v (continuing with degraded status)", err)
		return nil
	}

	c.mu.Lock()
	c.pool = pool
	c.ready = true
	c.mu.Unlock()

	log.Println("[postgres] Connected successfully to PostgreSQL")
	return nil
}

// Ping checks whether PostgreSQL connection pool is alive and responding.
func (c *Client) Ping(ctx context.Context) error {
	if c.mode == config.DatabaseModeDisabled {
		return ErrDatabaseDisabled
	}

	c.mu.RLock()
	pool := c.pool
	ready := c.ready
	c.mu.RUnlock()

	if !ready || pool == nil {
		return errors.New("postgres: connection pool not initialized")
	}

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return pool.Ping(pingCtx)
}

// Status returns the current database operational status: "disabled", "connected", "degraded", or "unavailable".
func (c *Client) Status(ctx context.Context) string {
	if c.mode == config.DatabaseModeDisabled {
		return "disabled"
	}

	if err := c.Ping(ctx); err == nil {
		return "connected"
	}

	if c.mode == config.DatabaseModeRequired {
		return "unavailable"
	}

	return "degraded"
}

// Mode returns the configured database operating mode.
func (c *Client) Mode() string {
	return c.mode
}

// Close closes the pgxpool connection pool.
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.pool != nil {
		c.pool.Close()
		c.pool = nil
		c.ready = false
		log.Println("[postgres] Connection pool closed cleanly")
	}
}

