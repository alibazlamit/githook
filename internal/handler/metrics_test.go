package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetrics(t *testing.T) {
	tests := []struct {
		name            string
		wantStatus      int
		wantContentType string
		wantBodyContain string
	}{
		{
			name:            "returns 200 OK",
			wantStatus:      http.StatusOK,
			wantContentType: "text/plain; version=0.0.4",
			wantBodyContain: "# HELP",
		},
		{
			name:            "body contains HELP line for go_goroutines",
			wantStatus:      http.StatusOK,
			wantContentType: "text/plain; version=0.0.4",
			wantBodyContain: "# HELP go_goroutines",
		},
		{
			name:            "body contains TYPE line",
			wantStatus:      http.StatusOK,
			wantContentType: "text/plain; version=0.0.4",
			wantBodyContain: "# TYPE go_goroutines gauge",
		},
		{
			name:            "body contains metric value line",
			wantStatus:      http.StatusOK,
			wantContentType: "text/plain; version=0.0.4",
			wantBodyContain: "go_goroutines ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
			rr := httptest.NewRecorder()

			Metrics(rr, req)

			if rr.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rr.Code, tc.wantStatus)
			}

			ct := rr.Header().Get("Content-Type")
			if ct != tc.wantContentType {
				t.Errorf("Content-Type = %q, want %q", ct, tc.wantContentType)
			}

			body := rr.Body.String()
			if !strings.Contains(body, tc.wantBodyContain) {
				t.Errorf("body does not contain %q; body = %q", tc.wantBodyContain, body)
			}
		})
	}
}
