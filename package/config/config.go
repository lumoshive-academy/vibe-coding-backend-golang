package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

// Config aggregates application configuration loaded from environment variables.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name  string
	Port  int
	Debug bool
}

type DatabaseConfig struct {
	Name          string
	Username      string
	Password      string
	Host          string
	Port          int
	Timezone      string
	MaxIdleConns  int
	MaxOpenConns  int
	MaxIdleTime   int
	MaxLifeTime   int
	SSLMode       string
	LogLevel      int
	MigrationsDir string
}

type JWTConfig struct {
	Secret string
	Issuer string
	TTL    int
}

var (
	config     Config
	configOnce sync.Once
)

// Load reads environment variables and returns a fully populated Config instance.
func Load() (Config, error) {
	var err error
	configOnce.Do(func() {
		_ = godotenv.Load()
		appConfig, e := loadAppConfig()
		if e != nil {
			err = fmt.Errorf("load app config: %w", e)
			return
		}
		dbConfig, e := loadDatabaseConfig()
		if e != nil {
			err = fmt.Errorf("load database config: %w", e)
			return
		}
		jwtConfig, e := loadJWTConfig()
		if e != nil {
			err = fmt.Errorf("load jwt config: %w", e)
			return
		}
		config = Config{
			App:      appConfig,
			Database: dbConfig,
			JWT:      jwtConfig,
		}
	})
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func loadAppConfig() (AppConfig, error) {
	port, err := intFromEnv("PORT", 8080)
	if err != nil {
		return AppConfig{}, err
	}
	debug, err := boolFromEnv("DEBUG", false)
	if err != nil {
		return AppConfig{}, err
	}
	return AppConfig{
		Name:  stringFromEnv("APP_NAME", "todolist"),
		Port:  port,
		Debug: debug,
	}, nil
}

func loadDatabaseConfig() (DatabaseConfig, error) {
	port, err := intFromEnv("DB_PORT", 5432)
	if err != nil {
		return DatabaseConfig{}, err
	}
	maxIdle, err := intFromEnv("DB_MAX_IDLE_CONNS", 5)
	if err != nil {
		return DatabaseConfig{}, err
	}
	maxOpen, err := intFromEnv("DB_MAX_OPEN_CONNS", 10)
	if err != nil {
		return DatabaseConfig{}, err
	}
	maxIdleTime, err := intFromEnv("DB_MAX_IDLE_TIME", 30)
	if err != nil {
		return DatabaseConfig{}, err
	}
	maxLifeTime, err := intFromEnv("DB_MAX_LIFE_TIME", 300)
	if err != nil {
		return DatabaseConfig{}, err
	}
	logLevel, err := intFromEnv("DB_LOG_LEVEL", 4)
	if err != nil {
		return DatabaseConfig{}, err
	}
	return DatabaseConfig{
		Name:          stringFromEnv("DB_NAME", "todolist"),
		Username:      stringFromEnv("DB_USERNAME", "postgres"),
		Password:      stringFromEnv("DB_PASSWORD", "postgres"),
		Host:          stringFromEnv("DB_HOST", "localhost"),
		Port:          port,
		Timezone:      stringFromEnv("DB_TIMEZONE", "UTC"),
		MaxIdleConns:  maxIdle,
		MaxOpenConns:  maxOpen,
		MaxIdleTime:   maxIdleTime,
		MaxLifeTime:   maxLifeTime,
		SSLMode:       stringFromEnv("DB_SSL_MODE", "disable"),
		LogLevel:      logLevel,
		MigrationsDir: stringFromEnv("DB_MIGRATIONS_DIR", "migrations"),
	}, nil
}

func loadJWTConfig() (JWTConfig, error) {
	ttl, err := intFromEnv("JWT_TTL", 3600)
	if err != nil {
		return JWTConfig{}, err
	}
	return JWTConfig{
		Secret: stringFromEnv("JWT_SECRET", "please-change-me"),
		Issuer: stringFromEnv("JWT_ISSUER", "todolist"),
		TTL:    ttl,
	}, nil
}

func stringFromEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func intFromEnv(key string, fallback int) (int, error) {
	if value, ok := os.LookupEnv(key); ok {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("parse %s: %w", key, err)
		}
		return parsed, nil
	}
	return fallback, nil
}

func boolFromEnv(key string, fallback bool) (bool, error) {
	if value, ok := os.LookupEnv(key); ok {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return false, fmt.Errorf("parse %s: %w", key, err)
		}
		return parsed, nil
	}
	return fallback, nil
}
