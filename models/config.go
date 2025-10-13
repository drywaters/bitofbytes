package models

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	CSRF struct {
		Key    []byte
		Secure bool
	}
	Server struct {
		Address string
	}
	Logging LoggingConfig
}

type LoggingConfig struct {
	Level  slog.Level
	Format LoggingFormat
}

type LoggingFormat string

const (
	LoggingFormatText LoggingFormat = "text"
	LoggingFormatJSON LoggingFormat = "json"
)

// LoadEnvConfig load new configuration from environment
func LoadEnvConfig() (Config, error) {
	var cfg Config
	if _, statErr := os.Stat(".env"); statErr == nil {
		if err := godotenv.Load(); err != nil {
			return cfg, fmt.Errorf("load config: %w", err)
		}
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return cfg, fmt.Errorf("load config: %w", statErr)
	}
	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		return cfg, errors.New("load config: SERVER_ADDRESS environment variable is required")
	}
	cfg.Server.Address = serverAddr

	secureEnv := os.Getenv("CSRF_SECURE")
	if secureEnv == "" {
		return cfg, errors.New("load config: CSRF_SECURE environment variable is required")
	}
	secure, err := strconv.ParseBool(secureEnv)
	if err != nil {
		return cfg, fmt.Errorf("load config: CSRF secure: %w", err)
	}
	cfg.CSRF.Secure = secure

	csrfKey := strings.TrimSpace(os.Getenv("CSRF_KEY"))
	if csrfKey == "" {
		keyPath := os.Getenv("CSRF_KEY_FILE")
		if keyPath == "" {
			keyPath = "/run/secrets/csrf_key"
		}

		keyData, err := os.ReadFile(keyPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return cfg, fmt.Errorf("load config: CSRF key file %q: %w", keyPath, err)
			}
			return cfg, fmt.Errorf("load config: read CSRF key file %q: %w", keyPath, err)
		}

		csrfKey = strings.TrimSpace(string(keyData))
	}

	if csrfKey == "" {
		return cfg, errors.New("load config: CSRF key is empty")
	}

	decodedKey, err := base64.StdEncoding.DecodeString(csrfKey)
	if err != nil {
		return cfg, fmt.Errorf("load config: decode CSRF key: %w", err)
	}

	if keyLen := len(decodedKey); keyLen != 32 {
		return cfg, fmt.Errorf("load config: CSRF key must decode to 32 bytes, got %d", keyLen)
	}

	cfg.CSRF.Key = decodedKey

	if err := loadLoggingConfig(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func loadLoggingConfig(cfg *Config) error {
	level := strings.TrimSpace(os.Getenv("LOG_LEVEL"))

	if level == "" {
		level = "info"
	}

	switch strings.ToLower(level) {
	case "debug":
		cfg.Logging.Level = slog.LevelDebug
	case "info":
		cfg.Logging.Level = slog.LevelInfo
	case "warn", "warning":
		cfg.Logging.Level = slog.LevelWarn
	case "error":
		cfg.Logging.Level = slog.LevelError
	default:
		return fmt.Errorf("load config: invalid LOG_LEVEL %q", level)
	}

	format := strings.TrimSpace(os.Getenv("LOG_FORMAT"))
	if format == "" {
		cfg.Logging.Format = LoggingFormatText
		return nil
	}

	switch strings.ToLower(format) {
	case string(LoggingFormatText):
		cfg.Logging.Format = LoggingFormatText
	case string(LoggingFormatJSON):
		cfg.Logging.Format = LoggingFormatJSON
	default:
		return fmt.Errorf("load config: invalid LOG_FORMAT %q", format)
	}

	return nil
}
