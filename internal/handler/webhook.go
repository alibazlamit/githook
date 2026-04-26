package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read request body")
		return
	}

	if !h.validSignature(r.Header.Get("X-Hub-Signature-256"), body) {
		writeError(w, http.StatusUnauthorized, "invalid or missing signature")
		return
	}

	// slices 4-5: idempotency check and ingest
	writeError(w, http.StatusInternalServerError, "not implemented")
}

// validSignature reports whether sig is a valid HMAC-SHA256 signature of body
// using h.secret. The sig must be prefixed with "sha256=".
func (h *WebhookHandler) validSignature(sig string, body []byte) bool {
	const prefix = "sha256="
	if !strings.HasPrefix(sig, prefix) {
		return false
	}
	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(sig[len(prefix):]), []byte(expected))
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
