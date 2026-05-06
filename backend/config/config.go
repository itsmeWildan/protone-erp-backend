package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
	MaxConns int
	MinConns int
}

type JWTConfig struct {
	Secret           string
	AccessTTLMinutes int
	RefreshTTLDays   int
}

func Load() (*Config, error) {
	// Load .env file (ignore error in production where env vars are set directly)
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "ProtoERP"),
			Env:  getEnv("APP_ENV", "production"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "protone_erp"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			MaxConns: getEnvInt("DB_MAX_CONNS", 25),
			MinConns: getEnvInt("DB_MIN_CONNS", 5),
		},
		JWT: JWTConfig{
			Secret:           getEnv("JWT_SECRET", ""),
			AccessTTLMinutes: getEnvInt("JWT_ACCESS_TTL_MINUTES", 60),
			RefreshTTLDays:   getEnvInt("JWT_REFRESH_TTL_DAYS", 7),
		},
	}

	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s pool_max_conns=%d pool_min_conns=%d",
		d.Host, d.Port, d.Name, d.User, d.Password, d.SSLMode, d.MaxConns, d.MinConns,
	)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
