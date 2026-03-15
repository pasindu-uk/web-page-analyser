package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port                int
	RequestTimeout      time.Duration
	MaxLinkCheckWorkers int
	MaxLinksToCheck     int
	LogLevel            string
	MySQLDSN            string
}

func Load() *Config {
	// Load .env.local first (local overrides), then .env as fallback.
	// Existing env vars take precedence — file values won't overwrite them.
	loadEnvFile(".env.local")
	loadEnvFile(".env")

	return &Config{
		Port:                getEnvInt("PORT", 8080),
		RequestTimeout:      getEnvDuration("REQUEST_TIMEOUT", 10*time.Second),
		MaxLinkCheckWorkers: getEnvInt("MAX_LINK_CHECK_WORKERS", 5),
		MaxLinksToCheck:     getEnvInt("MAX_LINKS_TO_CHECK", 50),
		LogLevel:            getEnvStr("LOG_LEVEL", "info"),
		MySQLDSN:            getEnvStr("MYSQL_DSN", ""),
	}
}

// loadEnvFile reads a .env file and sets any variables not already in the environment.
func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return // file doesn't exist, skip silently
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		// Don't overwrite existing env vars
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
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
