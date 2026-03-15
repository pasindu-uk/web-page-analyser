package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                int
	RequestTimeout      time.Duration
	MaxLinkCheckWorkers int
	MaxLinksToCheck     int
	LogLevel            string
}

func Load() *Config {
	return &Config{
		Port:                getEnvInt("PORT", 8080),
		RequestTimeout:      getEnvDuration("REQUEST_TIMEOUT", 10*time.Second),
		MaxLinkCheckWorkers: getEnvInt("MAX_LINK_CHECK_WORKERS", 5),
		MaxLinksToCheck:     getEnvInt("MAX_LINKS_TO_CHECK", 50),
		LogLevel:            getEnvStr("LOG_LEVEL", "info"),
	}
}

func getEnvStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
