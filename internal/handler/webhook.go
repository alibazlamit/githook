package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type webhookIngestor interface {
	ValidateSignature(sig string, body []byte) bool
	Ingest(ctx context.Context, deliveryID, eventType string, payload []byte) error
	CheckDuplicate(ctx context.Context, deliveryID string) (bool, error)
}

type WebhookHandler struct {
	svc webhookIngestor
}

func NewWebhookHandler(svc webhookIngestor) *WebhookHandler {
	return &WebhookHandler{svc: svc}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read request body")
		return
	}

	if !h.svc.ValidateSignature(r.Header.Get("X-Hub-Signature-256"), body) {
		writeError(w, http.StatusUnauthorized, "invalid or missing signature")
		return
	}

	deliveryID := r.Header.Get("X-GitHub-Delivery")
	dup, err := h.svc.CheckDuplicate(r.Context(), deliveryID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "idempotency check failed")
		return
	}
	if dup {
		writeError(w, http.StatusConflict, "duplicate delivery")
		return
	}

	// slice 5: ingest
	writeError(w, http.StatusInternalServerError, "not implemented")
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg}) // error unrecoverable after WriteHeader
}
