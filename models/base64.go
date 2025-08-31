package models

import (
	"bytes"
	"encoding/base64"
)

type Base64Service struct{}

func (s *Base64Service) Encode(b []byte) string {

	return base64.StdEncoding.EncodeToString(b)
}

// EncodeWith encodes using URL-safe and/or no-padding variants when requested.
func (s *Base64Service) EncodeWith(b []byte, urlSafe, noPad bool) string {
	var enc *base64.Encoding
	switch {
	case urlSafe && noPad:
		enc = base64.RawURLEncoding
	case urlSafe && !noPad:
		enc = base64.URLEncoding
	case !urlSafe && noPad:
		enc = base64.RawStdEncoding
	default:
		enc = base64.StdEncoding
	}
	return enc.EncodeToString(b)
}

func (s *Base64Service) Decode(b []byte) (string, error) {
	src := bytes.TrimSpace(b)
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))

	n, err := base64.StdEncoding.Decode(dst, src)

	return string(dst[:n]), err
}

// DecodeWith decodes using URL-safe and/or no-padding variants when requested.
func (s *Base64Service) DecodeWith(b []byte, urlSafe, noPad bool) (string, error) {
	src := bytes.TrimSpace(b)
	var enc *base64.Encoding
	switch {
	case urlSafe && noPad:
		enc = base64.RawURLEncoding
	case urlSafe && !noPad:
		enc = base64.URLEncoding
	case !urlSafe && noPad:
		enc = base64.RawStdEncoding
	default:
		enc = base64.StdEncoding
	}
	dst := make([]byte, enc.DecodedLen(len(src)))
	n, err := enc.Decode(dst, src)
	return string(dst[:n]), err
}
