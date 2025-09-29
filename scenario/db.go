package scenario

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDBPool() (*pgxpool.Pool, error) {
	conn := os.Getenv("DATABASE_URL")
	if conn == "" {
		// return nil, fmt.Errorf("Invalid DATABASE_URL")
		conn = "postgres://admin:admin@localhost:5432/postgres"
	}
	config, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse pgxpool config: %s", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("Failed to create pool")
	}

	return pool, nil
}
