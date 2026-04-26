package service_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ali/githook/internal/service"
)

func TestWebhookServiceValidateSignature(t *testing.T) {
	const secret = "test-secret"
	svc := service.NewWebhookService(secret)
	body := []byte(`{"action":"opened"}`)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	validSig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	tests := []struct {
		name string
		sig  string
		want bool
	}{
		{"valid signature", validSig, true},
		{"missing header", "", false},
		{"wrong prefix", "sha1=abc123", false},
		{"wrong secret", "sha256=" + strings.Repeat("a", 64), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.ValidateSignature(tt.sig, body)
			if got != tt.want {
				t.Errorf("TestWebhookServiceValidateSignature/%s: want %v, got %v", tt.name, tt.want, got)
			}
		})
	}
}
