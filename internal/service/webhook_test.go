package service_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ali/githook/internal/repository"
	"github.com/ali/githook/internal/service"
)

type mockRepo struct {
	exists bool
	err    error
}

func (m *mockRepo) Insert(_ context.Context, _ *repository.WebhookEvent) error { return nil }
func (m *mockRepo) UpdateStatus(_ context.Context, _, _ string, _ *time.Time) error {
	return nil
}
func (m *mockRepo) FindStuck(_ context.Context, _ time.Duration) ([]*repository.WebhookEvent, error) {
	return nil, nil
}
func (m *mockRepo) DeliveryExists(_ context.Context, _ string) (bool, error) {
	return m.exists, m.err
}

func TestWebhookServiceValidateSignature(t *testing.T) {
	const secret = "test-secret"
	svc := service.NewWebhookService(secret, &mockRepo{})
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

func TestWebhookServiceCheckDuplicate(t *testing.T) {
	repoErr := errors.New("db down")

	tests := []struct {
		name    string
		exists  bool
		repoErr error
		want    bool
		wantErr bool
	}{
		{"not duplicate", false, nil, false, false},
		{"is duplicate", true, nil, true, false},
		{"repo error", false, repoErr, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewWebhookService("secret", &mockRepo{exists: tt.exists, err: tt.repoErr})
			got, err := svc.CheckDuplicate(context.Background(), "delivery-id")
			if (err != nil) != tt.wantErr {
				t.Errorf("err: want err=%v, got %v", tt.wantErr, err)
			}
			if got != tt.want {
				t.Errorf("dup: want %v, got %v", tt.want, got)
			}
		})
	}
}
