package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/pasindu-uk/web-page-analyser/internal/fetcher"
	"github.com/pasindu-uk/web-page-analyser/internal/model"
	"github.com/pasindu-uk/web-page-analyser/internal/service"
)

type Handler struct {
	service *service.AnalyzeService
}

func New(svc *service.AnalyzeService) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/analyze", h.handleAnalyze)
	mux.HandleFunc("GET /api/analyses", h.handleListAnalyses)
	mux.HandleFunc("DELETE /api/cache", h.handleClearCache)
}

const maxRequestBodySize = 1 << 20 // 1 MB

func (h *Handler) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

	var req model.AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	resp, err := h.service.Analyze(r.Context(), req.URL)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleListAnalyses(w http.ResponseWriter, r *http.Request) {
	results, err := h.service.ListAnalyses(r.Context())
	if err != nil {
		slog.Error("failed to list analyses", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to retrieve analysis history")
		return
	}

	if results == nil {
		writeError(w, http.StatusServiceUnavailable, "persistence is not configured")
		return
	}

	writeJSON(w, http.StatusOK, results)
}

func (h *Handler) handleClearCache(w http.ResponseWriter, r *http.Request) {
	if h.service.ClearCache() {
		writeJSON(w, http.StatusOK, map[string]string{"status": "cache cleared"})
	} else {
		writeJSON(w, http.StatusOK, map[string]string{"status": "no cache configured"})
	}
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	var validationErr *service.ValidationError
	if errors.As(err, &validationErr) {
		writeError(w, http.StatusBadRequest, validationErr.Message)
		return
	}

	var httpErr *fetcher.HTTPError
	if errors.As(err, &httpErr) {
		writeError(w, httpErr.StatusCode, httpErr.Error())
		return
	}

	slog.Error("analysis failed", "error", err)
	writeError(w, http.StatusBadGateway, err.Error())
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, model.ErrorResponse{
		StatusCode: status,
		Message:    message,
	})
}

// CORSMiddleware adds CORS headers for frontend integration.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
