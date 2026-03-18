package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pasindu-uk/web-page-analyser/internal/analyzer"
	"github.com/pasindu-uk/web-page-analyser/internal/config"
	"github.com/pasindu-uk/web-page-analyser/internal/fetcher"
	"github.com/pasindu-uk/web-page-analyser/internal/handler"
	"github.com/pasindu-uk/web-page-analyser/internal/logger"
	"github.com/pasindu-uk/web-page-analyser/internal/repository"
	"github.com/pasindu-uk/web-page-analyser/internal/service"
)

func main() {
	cfg := config.Load()
	logger.Setup(cfg.LogLevel)

	var repo repository.Repository
	if cfg.MySQLDSN != "" {
		mysqlRepo, err := repository.NewMySQL(cfg.MySQLDSN)
		if err != nil {
			slog.Error("failed to connect to MySQL", "error", err)
			os.Exit(1)
		}
		defer mysqlRepo.Close()

		if err := repository.RunMigrations(mysqlRepo.DB()); err != nil {
			slog.Error("failed to run migrations", "error", err)
			os.Exit(1)
		}
		slog.Info("MySQL persistence enabled")
		repo = repository.NewCached(mysqlRepo)
	}

	f := fetcher.New(cfg.RequestTimeout)
	lc := analyzer.NewLinkChecker(cfg.MaxLinkCheckWorkers, cfg.RequestTimeout)
	svc := service.New(f, lc, repo)
	h := handler.New(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	h.RegisterRoutes(mux)

	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler.CORSMiddleware(mux),
	}

	go func() {
		slog.Info("server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}
}
