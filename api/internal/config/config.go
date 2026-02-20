package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	Port        string `json:"port"`
	DatabaseURL string `json:"database_url"`
	MLEngineURL string `json:"ml_engine_url"`
	LogLevel    string `json:"log_level"`

	// Auth configuration
	JWTSecret    string   `json:"jwt_secret"`
	JWTIssuer    string   `json:"jwt_issuer"`
	APIKeys      []string `json:"api_keys"` // format: "key:subject:role"
	CORSOrigins  []string `json:"cors_origins"`
	MaxBodyBytes int64    `json:"max_body_bytes"`
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Port:         getEnv("QRAP_PORT", "8083"),
		DatabaseURL:  getEnv("QRAP_DATABASE_URL", ""),
		MLEngineURL:  getEnv("QRAP_ML_ENGINE_URL", "http://127.0.0.1:8084"),
		LogLevel:     getEnv("QRAP_LOG_LEVEL", "info"),
		JWTSecret:    getEnv("QUANTUN_JWT_SECRET", ""),
		JWTIssuer:    getEnv("QUANTUN_JWT_ISSUER", "quantun"),
		MaxBodyBytes: 1 << 20, // 1 MB
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("QRAP_DATABASE_URL is required")
	}

	// Parse API keys (comma-separated, format: key:subject:role)
	if apiKeysStr := getEnv("QUANTUN_API_KEYS", ""); apiKeysStr != "" {
		cfg.APIKeys = strings.Split(apiKeysStr, ",")
	}

	// Parse CORS origins (comma-separated)
	if corsStr := getEnv("QUANTUN_CORS_ORIGINS", ""); corsStr != "" {
		cfg.CORSOrigins = strings.Split(corsStr, ",")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
