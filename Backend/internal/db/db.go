package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func Connect(dbURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	var pool *pgxpool.Pool

	// retry connection
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				cancel()
				return pool, nil
			}
		}

		cancel()
		log.Println("⏳ waiting for database...")
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after retries")
}

// func Connect(dbURL string) (*Database, error) {
// 	config, err := pgxpool.ParseConfig(dbURL)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to parse database URL: %w", err)
// 	}

// 	config.MaxConns = 25
// 	config.MinConns = 5
// 	config.MaxConnLifetime = 30 * time.Minute
// 	config.MaxConnIdleTime = 5 * time.Minute

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	pool, err := pgxpool.NewWithConfig(ctx, config)

// 	if err != nil {
// 		return nil, fmt.Errorf("could not create connection pool: %w", err)
// 	}

// 	if err := pool.Ping(ctx); err != nil {
// 		return nil, fmt.Errorf("database unreachable: %w", err)
// 	}

// 	return &Database{Pool: pool}, nil
// }

func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
