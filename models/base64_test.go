package models

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"
)

func TestBase64Encode(t *testing.T) {
	t.Parallel()

	service := Base64Service{}
	input := []byte("foobar")

	got, err := service.Encode(input, "standard")
	if err != nil {
		t.Fatalf("Encode(%q) unexpected error: %v", input, err)
	}

	if want := base64.StdEncoding.EncodeToString(input); got != want {
		t.Fatalf("Encode(%q) = %q, want %q", input, got, want)
	}
}

func TestBase64Decode(t *testing.T) {
	t.Parallel()

	service := Base64Service{}
	encoded := base64.StdEncoding.EncodeToString([]byte("foobar"))

	tests := []struct {
		name     string
		input    string
		want     string
		wantErr  error
		errCheck func(error) bool
	}{
		{
			name:  "trims whitespace before decoding",
			input: "  \n" + encoded + "\t  ",
			want:  "foobar",
		},
		{
			name:  "rejects malformed payloads",
			input: "@@@",
			errCheck: func(err error) bool {
				_, ok := err.(base64.CorruptInputError)
				return ok
			},
		},
		{
			name:    "enforces maximum accepted length",
			input:   strings.Repeat("A", maxBase64DecodeLen+1),
			wantErr: ErrBase64InputTooLarge,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decoded, err := service.Decode([]byte(tt.input), "standard")

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Decode(%q) error = %v, want %v", tt.input, err, tt.wantErr)
				}

				return
			}

			if tt.errCheck != nil {
				if tt.errCheck(err) {
					return
				}
				t.Fatalf("Decode(%q) error = %v, expected matching error", tt.input, err)
			}

			if err != nil {
				t.Fatalf("Decode(%q) unexpected error: %v", tt.input, err)
			}

			if decoded != tt.want {
				t.Fatalf("Decode(%q) = %q, want %q", tt.input, decoded, tt.want)
			}
		})
	}
}

func TestBase64SupportsAdditionalEncodings(t *testing.T) {
	t.Parallel()

	service := Base64Service{}

	tests := []struct {
		name          string
		encodingName  string
		encodedSample string
	}{
		{
			name:          "URL safe variant",
			encodingName:  "url",
			encodedSample: base64.URLEncoding.EncodeToString([]byte("hello")),
		},
		{
			name:          "raw standard variant",
			encodingName:  "raw-standard",
			encodedSample: base64.RawStdEncoding.EncodeToString([]byte("hello")),
		},
		{
			name:          "raw URL safe variant",
			encodingName:  "raw-url",
			encodedSample: base64.RawURLEncoding.EncodeToString([]byte("hello")),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			encoded, err := service.Encode([]byte("hello"), tt.encodingName)
			if err != nil {
				t.Fatalf("Encode unexpected error: %v", err)
			}
			if encoded != tt.encodedSample {
				t.Fatalf("Encode = %q, want %q", encoded, tt.encodedSample)
			}

			decoded, err := service.Decode([]byte(encoded), tt.encodingName)
			if err != nil {
				t.Fatalf("Decode unexpected error: %v", err)
			}
			if decoded != "hello" {
				t.Fatalf("Decode = %q, want %q", decoded, "hello")
			}
		})
	}
}

func TestBase64RejectsUnknownEncoding(t *testing.T) {
	t.Parallel()

	service := Base64Service{}

	if _, err := service.Encode([]byte("hello"), "does-not-exist"); !errors.Is(err, ErrBase64UnknownEncoding) {
		t.Fatalf("Encode error = %v, want %v", err, ErrBase64UnknownEncoding)
	}

	if _, err := service.Decode([]byte("hello"), "does-not-exist"); !errors.Is(err, ErrBase64UnknownEncoding) {
		t.Fatalf("Decode error = %v, want %v", err, ErrBase64UnknownEncoding)
	}
}
