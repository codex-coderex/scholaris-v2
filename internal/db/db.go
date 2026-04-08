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
			college_code VARCHAR(20)  REFERENCES college(code) ON DELETE SET NULL
		);

		CREATE TABLE IF NOT EXISTS student (
			id           VARCHAR(20)  PRIMARY KEY,
			first_name   VARCHAR(100) NOT NULL,
			last_name    VARCHAR(100) NOT NULL,
			year         INT          NOT NULL CHECK (year BETWEEN 1 AND 4),
			gender       VARCHAR(10)  NOT NULL,
			program_code VARCHAR(20)  REFERENCES program(code) ON DELETE SET NULL
		);

		ALTER TABLE IF EXISTS program
			ALTER COLUMN college_code DROP NOT NULL;

		ALTER TABLE IF EXISTS student
			ALTER COLUMN program_code DROP NOT NULL;

		ALTER TABLE IF EXISTS program
			DROP CONSTRAINT IF EXISTS program_college_code_fkey;

		ALTER TABLE IF EXISTS program
			ADD CONSTRAINT program_college_code_fkey
			FOREIGN KEY (college_code) REFERENCES college(code) ON DELETE SET NULL;

		ALTER TABLE IF EXISTS student
			DROP CONSTRAINT IF EXISTS student_program_code_fkey;

		ALTER TABLE IF EXISTS student
			ADD CONSTRAINT student_program_code_fkey
			FOREIGN KEY (program_code) REFERENCES program(code) ON DELETE SET NULL;
	`)
	return err
}
