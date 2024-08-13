package models

import (
	"encoding/base64"
)

type Base64Service struct{}

func (s *Base64Service) Encode(b []byte) string {

	return base64.StdEncoding.EncodeToString(b)
}

func (s *Base64Service) Decode(b []byte) (string, error) {
	dst := make([]byte, len(b))

	n, err := base64.StdEncoding.Decode(dst, b)

	return string(dst[:n]), err
}
