package handler

import (
	"context"
	"net/http"
)

type webhookIngestor interface {
	Ingest(ctx context.Context, deliveryID, eventType string, payload []byte) error
}

type WebhookHandler struct {
	svc    webhookIngestor
	secret string
}

func NewWebhookHandler(svc webhookIngestor, secret string) *WebhookHandler {
	return &WebhookHandler{svc: svc, secret: secret}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// implemented in slices 3-5
	http.Error(w, `{"error":"not implemented"}`, http.StatusInternalServerError)
}
