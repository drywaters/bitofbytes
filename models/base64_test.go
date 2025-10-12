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

	if got, want := service.Encode(input), base64.StdEncoding.EncodeToString(input); got != want {
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

			decoded, err := service.Decode([]byte(tt.input))

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
