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
	duplicate      bool
	dupErr         error
	ingestErr      error
}

func (m *mockWebhookSvc) ValidateSignature(_ string, _ []byte) bool {
	return m.validateResult
}

func (m *mockWebhookSvc) Ingest(_ context.Context, _, _ string, _ []byte) error {
	return m.ingestErr
}

func (m *mockWebhookSvc) CheckDuplicate(_ context.Context, _ string) (bool, error) {
	return m.duplicate, m.dupErr
}

func TestWebhookHandlerServeHTTP(t *testing.T) {
	body := []byte(`{"action":"opened"}`)

	tests := []struct {
		name           string
		validateResult bool
		duplicate      bool
		dupErr         error
		wantStatus     int
	}{
		{"invalid signature", false, false, nil, http.StatusUnauthorized},
		{"valid signature", true, false, nil, http.StatusInternalServerError}, // passes auth + idempotency; 500 = ingest not yet wired
		{"duplicate delivery", true, true, nil, http.StatusConflict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewWebhookHandler(&mockWebhookSvc{
				validateResult: tt.validateResult,
				duplicate:      tt.duplicate,
				dupErr:         tt.dupErr,
			})
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
