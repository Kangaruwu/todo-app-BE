package db

import (
    "context"
    "log"
    "os"

    "github.com/joho/godotenv"
    "github.com/jackc/pgx/v5"
)

func Connect() (*pgx.Conn, error) {
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found, using system env")
    }

    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        log.Fatal("DATABASE_URL is not set")
    }
    ctx := context.Background()
    db, err := pgx.Connect(ctx, connStr)
	if err != nil {
        panic(err)
    }
    defer db.Close(ctx)

	// Test the connection by querying the database version
	var version string
    err = db.QueryRow(context.Background(), "select version()").Scan(&version)
    if err != nil {
        panic(err)
    }
	log.Printf("Connected to database: %s\n", version)

    return db, nil
}
