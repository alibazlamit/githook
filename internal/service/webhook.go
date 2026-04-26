package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ali/githook/internal/repository"
)

type WebhookService interface {
	ValidateSignature(sig string, body []byte) bool
	Ingest(ctx context.Context, deliveryID, eventType string, payload []byte) error
	CheckDuplicate(ctx context.Context, deliveryID string) (bool, error)
}

type webhookService struct {
	secret string
	repo   repository.WebhookRepository
}

func NewWebhookService(secret string, repo repository.WebhookRepository) WebhookService {
	return &webhookService{secret: secret, repo: repo}
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

func (s *webhookService) Ingest(ctx context.Context, deliveryID, eventType string, payload []byte) error {
	// implemented in slice 5
	return nil
}

func (s *webhookService) CheckDuplicate(ctx context.Context, deliveryID string) (bool, error) {
	exists, err := s.repo.DeliveryExists(ctx, deliveryID)
	if err != nil {
		return false, fmt.Errorf("CheckDuplicate: %w", err)
	}
	return exists, nil
}
