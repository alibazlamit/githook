package handler

import (
	"fmt"
	"net/http"
	"runtime"
)

// Wire up with: mux.HandleFunc("GET /metrics", Metrics)

// Metrics writes a minimal Prometheus text exposition (version 0.0.4).
// It reports go_goroutines using runtime.NumGoroutine() from stdlib.
// No external dependencies are required.
func Metrics(w http.ResponseWriter, r *http.Request) {
	goroutines := runtime.NumGoroutine()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "# HELP go_goroutines Number of goroutines.\n")
	fmt.Fprintf(w, "# TYPE go_goroutines gauge\n")
	fmt.Fprintf(w, "go_goroutines %d\n", goroutines)
}
