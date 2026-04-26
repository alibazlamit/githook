package handler_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ali/githook/internal/handler"
)

func TestWebhookHandlerSignatureValidation(t *testing.T) {
	const secret = "test-secret"
	body := []byte(`{"action":"opened"}`)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	validSig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	tests := []struct {
		name       string
		sigHeader  string
		wantStatus int
	}{
		{"missing header", "", http.StatusUnauthorized},
		{"wrong prefix", "sha1=abc123", http.StatusUnauthorized},
		{"wrong secret", "sha256=" + strings.Repeat("a", 64), http.StatusUnauthorized},
		{"valid signature", validSig, http.StatusInternalServerError}, // passes auth; 500 = not yet wired
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewWebhookHandler(nil, secret)
			req := httptest.NewRequest(http.MethodPost, "/webhook/github", bytes.NewReader(body))
			if tt.sigHeader != "" {
				req.Header.Set("X-Hub-Signature-256", tt.sigHeader)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}
