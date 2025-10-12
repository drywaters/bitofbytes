package models

import (
	"bytes"
	"encoding/base64"
	"errors"
)

const maxBase64DecodeLen = 4 * 1024

var ErrBase64InputTooLarge = errors.New("base64 input exceeds maximum length")

type Base64Service struct{}

func (s *Base64Service) Encode(b []byte) string {

	return base64.StdEncoding.EncodeToString(b)
}

func (s *Base64Service) Decode(b []byte) (string, error) {
	trimmed := bytes.TrimSpace(b)
	if len(trimmed) > maxBase64DecodeLen {
		return "", ErrBase64InputTooLarge
	}

	dst := make([]byte, base64.StdEncoding.DecodedLen(len(trimmed)))

	n, err := base64.StdEncoding.Decode(dst, trimmed)
	if err != nil {
		return "", err
	}

	return string(dst[:n]), nil
}
