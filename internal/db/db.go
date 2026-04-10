package db

import (
	"fmt"

	"scholaris-v2/internal/config"
	"scholaris-v2/internal/shared/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// connect

func Init(cfg config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.GetDBConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	pingCtx, pingCancel := utils.NewDBContext()
	defer pingCancel()

	if err := pool.Ping(pingCtx); err != nil {
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
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	_, err := pool.Exec(ctx, `
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
