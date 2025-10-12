package models

import (
	"bytes"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func withTempEnv(env map[string]string, fn func(t *testing.T)) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		original := make(map[string]string, len(env))
		for k := range env {
			original[k] = os.Getenv(k)
		}
		for k, v := range env {
			if v == "" {
				os.Unsetenv(k)
				continue
			}
			os.Setenv(k, v)
		}
		t.Cleanup(func() {
			for k, v := range original {
				if v == "" {
					os.Unsetenv(k)
					continue
				}
				os.Setenv(k, v)
			}
		})
		fn(t)
	}
}

func TestLoadEnvConfigRequiresEnv(t *testing.T) {
	validKey := base64.StdEncoding.EncodeToString([]byte("01234567890123456789012345678901"))

	t.Run("errors when server address missing", withTempEnv(map[string]string{
		"SERVER_ADDRESS": "",
		"CSRF_SECURE":    "true",
		"CSRF_KEY":       validKey,
	}, func(t *testing.T) {
		if _, err := LoadEnvConfig(); err == nil {
			t.Fatalf("expected error when SERVER_ADDRESS missing")
		}
	}))

	t.Run("errors when csrf secure missing", withTempEnv(map[string]string{
		"SERVER_ADDRESS": "127.0.0.1:8080",
		"CSRF_SECURE":    "",
		"CSRF_KEY":       validKey,
	}, func(t *testing.T) {
		if _, err := LoadEnvConfig(); err == nil {
			t.Fatalf("expected error when CSRF_SECURE missing")
		}
	}))

	t.Run("loads config from env", withTempEnv(map[string]string{
		"SERVER_ADDRESS": "127.0.0.1:8080",
		"CSRF_SECURE":    "false",
		"CSRF_KEY":       validKey,
	}, func(t *testing.T) {
		cfg, err := LoadEnvConfig()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if cfg.Server.Address != "127.0.0.1:8080" {
			t.Fatalf("expected server address from env, got %q", cfg.Server.Address)
		}
		if cfg.CSRF.Secure {
			t.Fatalf("expected CSRF secure false from env")
		}
		wantKey := []byte("01234567890123456789012345678901")
		if !bytes.Equal(cfg.CSRF.Key, wantKey) {
			t.Fatalf("unexpected CSRF key %q", cfg.CSRF.Key)
		}
	}))
}

func TestLoadEnvConfigReadsSecretFile(t *testing.T) {
	dir := t.TempDir()
	secretPath := filepath.Join(dir, "csrf_key")
	fileSecret := []byte("file-secret-32-bytes--0123456789")
	fileKey := base64.StdEncoding.EncodeToString(fileSecret)
	if err := os.WriteFile(secretPath, []byte(fileKey+"\n"), 0600); err != nil {
		t.Fatalf("write secret: %v", err)
	}

	t.Run("loads from file when env empty", withTempEnv(map[string]string{
		"SERVER_ADDRESS": "127.0.0.1:8080",
		"CSRF_SECURE":    "true",
		"CSRF_KEY":       "",
		"CSRF_KEY_FILE":  secretPath,
	}, func(t *testing.T) {
		cfg, err := LoadEnvConfig()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		wantKey := fileSecret
		if !bytes.Equal(cfg.CSRF.Key, wantKey) {
			t.Fatalf("expected CSRF key trimmed from file, got %q", cfg.CSRF.Key)
		}
	}))

	t.Run("errors when file missing", withTempEnv(map[string]string{
		"SERVER_ADDRESS": "127.0.0.1:8080",
		"CSRF_SECURE":    "true",
		"CSRF_KEY":       "",
		"CSRF_KEY_FILE":  filepath.Join(dir, "missing"),
	}, func(t *testing.T) {
		if _, err := LoadEnvConfig(); err == nil {
			t.Fatalf("expected error when CSRF key file missing")
		}
	}))
}

func TestLoadEnvConfigValidatesCSRFKey(t *testing.T) {
	t.Run("errors when key not base64", withTempEnv(map[string]string{
		"SERVER_ADDRESS": "127.0.0.1:8080",
		"CSRF_SECURE":    "true",
		"CSRF_KEY":       "not-base64!!",
	}, func(t *testing.T) {
		if _, err := LoadEnvConfig(); err == nil {
			t.Fatalf("expected error when CSRF key is not base64 encoded")
		}
	}))

	shortKey := base64.StdEncoding.EncodeToString([]byte("short key"))
	t.Run("errors when key wrong length", withTempEnv(map[string]string{
		"SERVER_ADDRESS": "127.0.0.1:8080",
		"CSRF_SECURE":    "true",
		"CSRF_KEY":       shortKey,
	}, func(t *testing.T) {
		if _, err := LoadEnvConfig(); err == nil {
			t.Fatalf("expected error when CSRF key decodes to wrong length")
		}
	}))

	longKey := base64.StdEncoding.EncodeToString([]byte("0123456789012345678901234567890101234567890123456789012345678901"))
	t.Run("errors when key decodes to 64 bytes", withTempEnv(map[string]string{
		"SERVER_ADDRESS": "127.0.0.1:8080",
		"CSRF_SECURE":    "true",
		"CSRF_KEY":       longKey,
	}, func(t *testing.T) {
		if _, err := LoadEnvConfig(); err == nil {
			t.Fatalf("expected error when CSRF key decodes to 64 bytes")
		}
	}))
}
