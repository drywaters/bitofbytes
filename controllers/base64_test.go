package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestEncode_JSON(t *testing.T) {
	u := Utils{}
	form := url.Values{"str": {"hi"}}
	r := httptest.NewRequest(http.MethodPost, "/utils/base64/encode", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	u.Encode(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["result"] != "aGk=" {
		t.Fatalf("result = %q", body["result"])
	}
}

func TestDecode_JSON_Invalid(t *testing.T) {
	u := Utils{}
	form := url.Values{"str": {"!!!!"}}
	r := httptest.NewRequest(http.MethodPost, "/utils/base64/decode", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	u.Decode(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["error"] == "" {
		t.Fatalf("expected error in body")
	}
}

func TestEncodeGET_Query_JSON(t *testing.T) {
	u := Utils{}
	r := httptest.NewRequest(http.MethodGet, "/utils/base64/encode?q=hi", nil)
	r.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	u.EncodeGET(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["result"] != "aGk=" {
		t.Fatalf("result = %q", body["result"])
	}
}

func TestDecodeGET_Query_TooLong(t *testing.T) {
	u := Utils{}
	long := strings.Repeat("a", maxBase64InputLen+1)
	r := httptest.NewRequest(http.MethodGet, "/utils/base64/decode?q="+long, nil)
	r.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	u.DecodeGET(w, r)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d", w.Code)
	}
}
