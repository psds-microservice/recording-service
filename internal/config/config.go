package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config holds recording-service configuration.
type Config struct {
	AppEnv   string
	AppHost  string
	GRPCPort string

	// StorageDir is the directory where recording files are written (e.g. ./recordings).
	StorageDir string
	// BaseURL is the base URL for recording links (e.g. https://cdn.example.com/recordings).
	// The final URL will be BaseURL + "/" + session_id + ".webm".
	BaseURL string
}

// Load loads config from environment.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		AppHost:    getEnv("APP_HOST", "0.0.0.0"),
		GRPCPort:   getEnv("GRPC_PORT", "8096"),
		StorageDir: getEnv("STORAGE_DIR", "./recordings"),
		BaseURL:    getEnv("RECORDING_BASE_URL", "http://localhost:8096/recordings"),
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate checks required fields.
func (c *Config) Validate() error {
	if c.StorageDir == "" {
		return errors.New("config: STORAGE_DIR is required")
	}
	return nil
}

// GRPCAddr returns listen address for gRPC server.
func (c *Config) GRPCAddr() string {
	return c.AppHost + ":" + c.GRPCPort
}

// RecordingPath returns the full path for a session recording file.
func (c *Config) RecordingPath(sessionID string) string {
	return filepath.Join(c.StorageDir, sessionID+".webm")
}

// RecordingURL returns the public URL for a session recording.
func (c *Config) RecordingURL(sessionID string) string {
	return c.BaseURL + "/" + sessionID + ".webm"
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
