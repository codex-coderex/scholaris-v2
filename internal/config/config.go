package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// config

type Config struct {
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
}

func (c Config) GetDBConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

// resolve

func resolveConfigPath() string {
	// check working directory first (wails dev)
	if wd, err := os.Getwd(); err == nil {
		candidate := filepath.Join(wd, "config.json")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	// check next to the executable (production)
	if exe, err := os.Executable(); err == nil {
		candidate := filepath.Join(filepath.Dir(exe), "config.json")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	// default path
	return "config.json"
}

// load

func Load() (Config, error) {
	defaults := Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "scholarisdb",
	}

	data, err := os.ReadFile(resolveConfigPath())
	if err != nil {
		return defaults, nil
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaults, fmt.Errorf("invalid config.json: %w", err)
	}

	return cfg, nil
}
