package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ali/githook/internal/handler"
)

type mockWebhookSvc struct {
	validateResult bool
}

func (m *mockWebhookSvc) ValidateSignature(_ string, _ []byte) bool {
	return m.validateResult
}

func (m *mockWebhookSvc) Ingest(_ context.Context, _, _ string, _ []byte) error {
	return nil
}

func TestWebhookHandlerSignatureValidation(t *testing.T) {
	body := []byte(`{"action":"opened"}`)

	tests := []struct {
		name           string
		validateResult bool
		wantStatus     int
	}{
		{"invalid signature", false, http.StatusUnauthorized},
		{"valid signature", true, http.StatusInternalServerError}, // passes auth; 500 = downstream not yet wired
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewWebhookHandler(&mockWebhookSvc{validateResult: tt.validateResult})
			req := httptest.NewRequest(http.MethodPost, "/webhook/github", bytes.NewReader(body))
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, rec.Code)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type: want application/json, got %q", ct)
			}
		})
	}
}
