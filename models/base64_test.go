package models

import (
	"strings"
	"testing"
)

func TestBase64Service_EncodeDecode(t *testing.T) {
	s := Base64Service{}
	cases := []struct {
		name    string
		input   string
		urlSafe bool
		noPad   bool
	}{
		{"std", "hello", false, false},
		{"noPad", "hello", false, true},
		{"urlSafe", "sample input", true, false},
		{"urlSafeNoPad", "sample input", true, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			enc := s.EncodeWith([]byte(tc.input), tc.urlSafe, tc.noPad)
			if tc.urlSafe {
				if strings.Contains(enc, "+") || strings.Contains(enc, "/") {
					t.Fatalf("expected URL-safe output, got %q", enc)
				}
			}
			dec, err := s.DecodeWith([]byte(enc), tc.urlSafe, tc.noPad)
			if err != nil {
				t.Fatalf("decode error: %v", err)
			}
			if dec != tc.input {
				t.Fatalf("roundtrip mismatch: got %q want %q", dec, tc.input)
			}
		})
	}
}

func TestBase64Service_Decode_Whitespace(t *testing.T) {
	s := Base64Service{}
	enc := s.Encode([]byte("hello"))
	withWS := "  \n\t" + enc + "  \n"
	dec, err := s.Decode([]byte(withWS))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dec != "hello" {
		t.Fatalf("got %q want %q", dec, "hello")
	}
}

func TestBase64Service_Decode_Invalid(t *testing.T) {
	s := Base64Service{}
	if _, err := s.Decode([]byte("!!!!")); err == nil {
		t.Fatalf("expected error for invalid input")
	}
}
