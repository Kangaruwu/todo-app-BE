package db

import (
	"context"
	"log"

	"go-backend-todo/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Singleton instances
var DbInstance *pgx.Conn
var DbPool *pgxpool.Pool

// ConnectPool creates connection pool (RECOMMENDED for production)
func ConnectPool() (*pgxpool.Pool, error) {
	cfg := config.Load()
	connStr := config.GetDatabaseURL(cfg)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	poolHealthCheck(pool)
	return pool, nil
}

// ConnectPoolWithConfig creates connection pool with specific config
func ConnectPoolWithConfig(cfg *config.Config) (*pgxpool.Pool, error) {
	connStr := config.GetDatabaseURL(cfg)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	poolHealthCheck(pool)
	return pool, nil
}

// GetPool returns singleton connection pool
func GetPool() *pgxpool.Pool {
	if DbPool == nil {
		var err error
		DbPool, err = ConnectPool()
		if err != nil {
			log.Fatalf("Failed to create connection pool: %v", err)
		}
	}
	return DbPool
}

// Connect establishes single connection (use only for testing or simple cases)
func Connect() (*pgx.Conn, error) {
	cfg := config.Load()
	connStr := config.GetDatabaseURL(cfg)

	ctx := context.Background()
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	connectHealthCheck(db)

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

	connectHealthCheck(db)
	return db, nil
}

// GetDB returns a singleton instance of the database connection
func GetDB() *pgx.Conn {
	if DbInstance == nil {
		var err error
		DbInstance, err = Connect()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	}
	return DbInstance
}

// For connection pool health check
func poolHealthCheck(pool *pgxpool.Pool) {
	ctx := context.Background()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection pool is healthy")

	DbPool = pool
}

// For single connection 
func connectHealthCheck(db *pgx.Conn) {
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection is healthy")

	DbInstance = db
}

