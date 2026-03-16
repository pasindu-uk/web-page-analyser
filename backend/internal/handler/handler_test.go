package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pasindu/web-page-analyser/internal/config"
	"github.com/pasindu/web-page-analyser/internal/model"
	"github.com/pasindu/web-page-analyser/internal/service"
)

func setupHandler() (*Handler, *http.ServeMux) {
	cfg := &config.Config{
		Port:                8080,
		RequestTimeout:      5 * time.Second,
		MaxLinkCheckWorkers: 2,
		MaxLinksToCheck:     10,
		LogLevel:            "info",
	}
	svc := service.New(cfg, nil)
	h := New(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return h, mux
}

func TestHandleAnalyze_InvalidJSON(t *testing.T) {
	_, mux := setupHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewBufferString("not json"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var errResp model.ErrorResponse
	json.NewDecoder(w.Body).Decode(&errResp)
	if errResp.Message == "" {
		t.Error("expected error message")
	}
}

func TestHandleAnalyze_EmptyURL(t *testing.T) {
	_, mux := setupHandler()

	body, _ := json.Marshal(model.AnalyzeRequest{URL: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleAnalyze_MethodNotAllowed(t *testing.T) {
	_, mux := setupHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/analyze", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestHandleAnalyze_Success(t *testing.T) {
	// Start a test server serving HTML
	htmlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html><html><head><title>Test</title></head><body><h1>Hello</h1></body></html>`))
	}))
	defer htmlServer.Close()

	_, mux := setupHandler()

	body, _ := json.Marshal(model.AnalyzeRequest{URL: htmlServer.URL})
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp model.AnalyzeResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Title != "Test" {
		t.Errorf("expected title 'Test', got '%s'", resp.Title)
	}
	if resp.Headings.H1 != 1 {
		t.Errorf("expected 1 h1, got %d", resp.Headings.H1)
	}
}

func TestHandleListAnalyses_NoPersistence(t *testing.T) {
	_, mux := setupHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/analyses", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 when persistence not configured, got %d", w.Code)
	}

	var errResp model.ErrorResponse
	json.NewDecoder(w.Body).Decode(&errResp)
	if errResp.Message == "" {
		t.Error("expected error message")
	}
}

func TestCORSMiddleware(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CORSMiddleware(mux)

	// Test preflight
	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204 for OPTIONS, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("missing CORS origin header")
	}

	// Test normal request has CORS headers
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("missing CORS origin header on normal request")
	}
}
