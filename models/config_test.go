package models

import (
	"encoding/base64"
	"os"
	"strings"
	"testing"
)

func writeEnvFile(t *testing.T, content string) {
	t.Helper()
	if err := os.WriteFile(".env", []byte(content), 0600); err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(".env")
	})
}

func unsetConfigEnvVars(t *testing.T) {
	t.Helper()
	keys := []string{"CSRF_KEY", "CSRF_SECURE", "SERVER_ADDRESS"}
	for _, key := range keys {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("failed to unset env %s: %v", key, err)
		}
	}
}

func TestLoadEnvConfigValidCSRFKey(t *testing.T) {
	unsetConfigEnvVars(t)
	validKey := "0123456789abcdef0123456789abcdef"
	encodedKey := base64.StdEncoding.EncodeToString([]byte(validKey))
	env := strings.Join([]string{
		"CSRF_KEY=" + encodedKey,
		"CSRF_SECURE=false",
		"SERVER_ADDRESS=:8080",
		"",
	}, "\n")
	writeEnvFile(t, env)

	cfg, err := LoadEnvConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(cfg.CSRF.Key) != validKey {
		t.Fatalf("unexpected csrf key: %x", cfg.CSRF.Key)
	}

	if cfg.CSRF.Secure {
		t.Fatalf("expected csrf secure to be false")
	}

	if cfg.Server.Address != ":8080" {
		t.Fatalf("unexpected server address: %s", cfg.Server.Address)
	}
}

func TestLoadEnvConfigInvalidCSRFKey(t *testing.T) {
	unsetConfigEnvVars(t)
	env := strings.Join([]string{
		"CSRF_KEY=" + base64.StdEncoding.EncodeToString([]byte("short-key")),
		"CSRF_SECURE=false",
		"SERVER_ADDRESS=:8080",
		"",
	}, "\n")
	writeEnvFile(t, env)

	_, err := LoadEnvConfig()
	if err == nil {
		t.Fatalf("expected error due to invalid csrf key length")
	}
}

func TestLoadEnvConfigInvalidBase64CSRFKey(t *testing.T) {
	unsetConfigEnvVars(t)
	env := strings.Join([]string{
		"CSRF_KEY=invalid-base64***",
		"CSRF_SECURE=false",
		"SERVER_ADDRESS=:8080",
		"",
	}, "\n")
	writeEnvFile(t, env)

	_, err := LoadEnvConfig()
	if err == nil {
		t.Fatalf("expected error due to invalid base64 encoding")
	}
}
