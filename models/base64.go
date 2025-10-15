package models

import (
	"bytes"
	"encoding/base64"
	"errors"
	"strings"
)

const maxBase64DecodeLen = 4 * 1024

var ErrBase64InputTooLarge = errors.New("base64 input exceeds maximum length")

var ErrBase64UnknownEncoding = errors.New("base64 encoding not supported")

type Base64Service struct{}

func (s *Base64Service) Encode(b []byte, encodingName string) (string, error) {
	encoding, err := s.lookupEncoding(encodingName)
	if err != nil {
		return "", err
	}

	return encoding.EncodeToString(b), nil
}

func (s *Base64Service) Decode(b []byte, encodingName string) (string, error) {
	encoding, err := s.lookupEncoding(encodingName)
	if err != nil {
		return "", err
	}

	trimmed := bytes.TrimSpace(b)
	if len(trimmed) > maxBase64DecodeLen {
		return "", ErrBase64InputTooLarge
	}

	dst := make([]byte, encoding.DecodedLen(len(trimmed)))

	n, err := encoding.Decode(dst, trimmed)
	if err != nil {
		return "", err
	}

	return string(dst[:n]), nil
}

func (s *Base64Service) lookupEncoding(name string) (*base64.Encoding, error) {
	normalized := strings.ToLower(strings.TrimSpace(name))
	if normalized == "" {
		normalized = "standard"
	}

	switch normalized {
	case "standard":
		return base64.StdEncoding, nil
	case "url", "url-safe":
		return base64.URLEncoding, nil
	case "raw-standard":
		return base64.RawStdEncoding, nil
	case "raw-url", "raw-url-safe":
		return base64.RawURLEncoding, nil
	default:
		return nil, ErrBase64UnknownEncoding
	}
}
