package models

import (
	"fmt"
	"log/slog"
	"os"
)

// NewLogger constructs a slog.Logger using the configured level and format.
func NewLogger(cfg LoggingConfig) (*slog.Logger, error) {
	handlerOpts := &slog.HandlerOptions{Level: cfg.Level}

	switch cfg.Format {
	case LoggingFormatText:
		return slog.New(slog.NewTextHandler(os.Stdout, handlerOpts)), nil
	case LoggingFormatJSON:
		return slog.New(slog.NewJSONHandler(os.Stdout, handlerOpts)), nil
	default:
		return nil, fmt.Errorf("new logger: unsupported logging format %q", cfg.Format)
	}
}
