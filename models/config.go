package models

import (
	"encoding/base64"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

const csrfKeyLength = 32

type Config struct {
	CSRF struct {
		Key    []byte
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
	rawKey := os.Getenv("CSRF_KEY")
	decodedKey, err := base64.StdEncoding.DecodeString(rawKey)
	if err != nil {
		return cfg, fmt.Errorf("load config: decode CSRF key: %w", err)
	}
	if len(decodedKey) != csrfKeyLength {
		return cfg, fmt.Errorf("load config: CSRF key must be %d bytes", csrfKeyLength)
	}
	cfg.CSRF.Key = decodedKey
	isSecure := os.Getenv("CSRF_SECURE")
	cfg.CSRF.Secure, err = strconv.ParseBool(isSecure)
	if err != nil {
		return cfg, fmt.Errorf("load config: CSRF secure: %w", err)
	}
	return cfg, nil
}
