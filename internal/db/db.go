package db

import (
	"context"
	"log"

	"go-backend-todo/internal/config"

	"github.com/jackc/pgx/v5"
)

// Connect establishes connection to PostgreSQL database using configuration
func Connect() (*pgx.Conn, error) {
	cfg := config.Load()
	connStr := config.GetDatabaseURL(cfg)

	ctx := context.Background()
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

    dbHealthCheck(db)
	return db, nil
}

// ConnectWithConfig connects using specific configuration
func ConnectWithConfig(cfg *config.Config) (*pgx.Conn, error) {
	connStr := config.GetDatabaseURL(cfg)

	ctx := context.Background()
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

    dbHealthCheck(db)
	return db, nil
}


func dbHealthCheck(db *pgx.Conn) error {
    ctx := context.Background()
    var version string
    err := db.QueryRow(ctx, "select version()").Scan(&version)
    if err != nil {
        return err
    }
    log.Printf("Database connection is healthy: %s\n", version)
    return nil
}