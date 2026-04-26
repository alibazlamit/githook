package main

import (
	"log"
	"net/http"

	"github.com/ali/githook/internal/config"
	"github.com/ali/githook/internal/handler"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)

	log.Printf("starting on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}
