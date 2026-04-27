package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBNAME     string
	DBHost     string
	DBPort     int
	ServerPort int

	MaxWorkersCount int
	MaxAttempts     int

	TimeSleep time.Duration
}

func Load() *Config {
	cfg := &Config{
		DBUser:          "postgres",
		DBPassword:      "",
		DBNAME:          "postgres",
		DBHost:          "localhost",
		DBPort:          5432,
		ServerPort:      8080,
		MaxWorkersCount: 3,
		MaxAttempts:     3,
		TimeSleep:       5 * time.Second,
	}

	if value := os.Getenv("DATABASE_USER"); value != "" {
		cfg.DBUser = value
	}

	if value := os.Getenv("DATABASE_PASSWORD"); value != "" {
		cfg.DBPassword = value
	}

	if value := os.Getenv("DATABASE_DBNAME"); value != "" {
		cfg.DBNAME = value
	}

	if value := os.Getenv("DATABASE_HOST"); value != "" {
		cfg.DBHost = value
	}

	if value := os.Getenv("DATABASE_PORT"); value != "" {
		if port, err := strconv.Atoi(value); err == nil {
			cfg.DBPort = port
		}
	}

	if value := os.Getenv("SERVER_PORT"); value != "" {
		if port, err := strconv.Atoi(value); err == nil {
			cfg.ServerPort = port
		}
	}

	if value := os.Getenv("MAX_WORKERS_COUNT"); value != "" {
		if maxWorkersCount, err := strconv.Atoi(value); err == nil {
			cfg.MaxWorkersCount = maxWorkersCount
		}
	}

	if value := os.Getenv("MAX_ATTEMPTS_COUNT"); value != "" {
		if maxAttemptsCount, err := strconv.Atoi(value); err == nil {
			cfg.MaxAttempts = maxAttemptsCount
		}
	}

	return cfg
}

func (cfg *Config) GetDBDSN() string {
	return "host=" + cfg.DBHost +
		" port=" + strconv.Itoa(cfg.DBPort) +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBNAME +
		" sslmode=disable"
}
