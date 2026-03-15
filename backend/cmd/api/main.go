package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/pasindu/web-page-analyser/internal/config"
	"github.com/pasindu/web-page-analyser/internal/logger"
)

func main() {
	cfg := config.Load()
	logger.Setup(cfg.LogLevel)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	slog.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server failed", "error", err)
	}
}
