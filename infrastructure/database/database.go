package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"runtime"
	"sync"
	"time"
)

type DB struct {
	Conn *pgxpool.Pool
	mu   sync.Mutex
}

// Close gracefully shuts down the database connection.
func (d *DB) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.Conn != nil {
		d.Conn.Close()
	}
	return nil
}

// Init initializes the database connection and applies migrations.
func Init(DBURI string) (*DB, error) {
	// Apply migrations
	if err := applyMigrations(DBURI); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Initialize database connection
	var db DB
	config, err := pgxpool.ParseConfig(DBURI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	config.MaxConns = int32(runtime.NumCPU())

	db.Conn, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database connection: %w", err)
	}

	// Test the connection
	if err = testConnection(db.Conn); err != nil {
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}

	// Start a connection worker for automatic reconnection
	go connectionWorker(&db, DBURI)

	return &db, nil
}

// testConnection pings the database to ensure the connection is alive.
func testConnection(db *pgxpool.Pool) error {
	if err := db.Ping(context.Background()); err != nil {
		return fmt.Errorf("can't ping database: %w", err)
	}
	return nil
}

// connectionWorker periodically checks the database connection and reconnects if necessary.
func connectionWorker(db *DB, dbURI string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := testConnection(db.Conn); err != nil {
			log.Printf("lost connection to database: %v", err)
			db.mu.Lock()
			newConn, err := pgxpool.New(context.Background(), dbURI)
			if err != nil {
				log.Fatalf("failed to reconnect to PostgreSQL: %v", err)
			}
			db.Conn.Close()
			db.Conn = newConn
			db.mu.Unlock()
			log.Println("successfully reconnected to PostgreSQL.")
		}
	}
}

// applyMigrations runs all pending database migrations.
func applyMigrations(DBURI string) error {
	m, err := migrate.New("file://migrations/", DBURI)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if sourceErr, dbErr := m.Close(); sourceErr != nil || dbErr != nil {
		return fmt.Errorf("failed to close migrate instance: sourceErr: %v, dbErr: %v", sourceErr, dbErr)
	}

	return nil
}
