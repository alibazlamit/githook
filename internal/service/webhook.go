package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type WebhookService interface {
	ValidateSignature(sig string, body []byte) bool
	Ingest(ctx context.Context, deliveryID, eventType string, payload []byte) error
}

type webhookService struct {
	secret string
}

func NewWebhookService(secret string) WebhookService {
	return &webhookService{secret: secret}
}

func (s *webhookService) ValidateSignature(sig string, body []byte) bool {
	const prefix = "sha256="
	if !strings.HasPrefix(sig, prefix) {
		return false
	}
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(sig[len(prefix):]), []byte(expected))
}

func (s *webhookService) Ingest(_ context.Context, _, _ string, _ []byte) error {
	// implemented in slice 5
	return nil
}
