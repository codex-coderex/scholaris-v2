package db

import (
	"context"
	"fmt"

	"scholaris-v2/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// connect

func Init(cfg config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.GetDBConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return pool, nil
}

// close

func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}

// create tables

func CreateTables(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS college (
			code VARCHAR(20)  PRIMARY KEY,
			name VARCHAR(255) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS program (
			code         VARCHAR(20)  PRIMARY KEY,
			name         VARCHAR(255) NOT NULL,
			college_code VARCHAR(20)  NOT NULL REFERENCES college(code)
		);

		CREATE TABLE IF NOT EXISTS student (
			id           VARCHAR(20)  PRIMARY KEY,
			first_name   VARCHAR(100) NOT NULL,
			last_name    VARCHAR(100) NOT NULL,
			year         INT          NOT NULL CHECK (year BETWEEN 1 AND 4),
			gender       VARCHAR(10)  NOT NULL,
			program_code VARCHAR(20)  NOT NULL REFERENCES program(code)
		);
	`)
	return err
}
