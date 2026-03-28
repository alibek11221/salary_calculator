package config

import (
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	defaultReadTimeout  = 15 * time.Second
	defaultWriteTimeout = 15 * time.Second
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	Env            string
	Port           string
	WorkdaysConfig WorkdaysConfig
	Server         ServerConfig
	Database       DatabaseConfig
}

type WorkdaysConfig struct {
	Dir      string
	CacheCap int
}

type ServerConfig struct {
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func GetConfig() *Config {
	once.Do(
		func() {
			instance = newConfig()
		},
	)
	return instance
}

func newConfig() *Config {
	return &Config{
		Env:  getEnvWithDefault("ENV", "development"),
		Port: getEnvWithDefault("PORT", "8080"),
		WorkdaysConfig: WorkdaysConfig{
			Dir:      getEnvWithDefault("WORKDAYS_DIR", "const/workdays"),
			CacheCap: getEnvAsInt("WORKDAYS_CACHE_CAP", 10),
		},
		Server: ServerConfig{
			ReadTimeout:    envToDuration("SERVER_READ_TIMEOUT", defaultReadTimeout),
			WriteTimeout:   envToDuration("SERVER_WRITE_TIMEOUT", defaultWriteTimeout),
			MaxHeaderBytes: 1 << 20, // 1 MB
		},
		Database: DatabaseConfig{
			Host:            getEnvWithDefault("DB_HOST", "localhost"),
			Port:            getEnvWithDefault("DB_PORT", "5432"),
			User:            getEnvWithDefault("DB_USERNAME", "postgres"),
			Password:        getEnvWithDefault("DB_PASSWORD", "postgres"),
			Name:            getEnvWithDefault("DB_DATABASE", "salary_calculator"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: envToDuration("DB_CONN_MAX_LIFETIME", time.Hour),
		},
	}
}

func envToDuration(key string, defaultValue time.Duration) time.Duration {
	if env := os.Getenv(key); env != "" {
		if d, err := time.ParseDuration(env); err == nil {
			return d
		}
	}
	return defaultValue
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
