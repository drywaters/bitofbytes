package models

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

// LoadEnvConfig load new configuration from environment
func LoadEnvConfig() (Config, error) {
	var cfg Config
	err := godotenv.Load()
	if err != nil {
		return cfg, fmt.Errorf("load config: %w", err)
	}
	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")
	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	isSecure := os.Getenv("CSRF_SECURE")
	cfg.CSRF.Secure, err = strconv.ParseBool(isSecure)
	if err != nil {
		return cfg, fmt.Errorf("load config: CSRF secure: %w", err)
	}
	return cfg, nil
}
