package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pasindu/web-page-analyser/internal/config"
)

func testConfig() *config.Config {
	return &config.Config{
		Port:                8080,
		RequestTimeout:      5 * time.Second,
		MaxLinkCheckWorkers: 2,
		MaxLinksToCheck:     10,
		LogLevel:            "info",
	}
}

func TestAnalyze_ValidationErrors(t *testing.T) {
	svc := New(testConfig(), nil)

	tests := []struct {
		name string
		url  string
	}{
		{"empty URL", ""},
		{"missing scheme", "example.com"},
		{"ftp scheme", "ftp://example.com"},
		{"no host", "http://"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Analyze(context.Background(), tt.url)
			if err == nil {
				t.Fatal("expected validation error")
			}
			if _, ok := err.(*ValidationError); !ok {
				t.Errorf("expected ValidationError, got %T: %v", err, err)
			}
		})
	}
}

func TestAnalyze_FullFlow(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>Test Page</title></head>
<body>
	<h1>Hello</h1>
	<h2>World</h2>
	<h2>Again</h2>
	<a href="/internal">Internal</a>
	<a href="https://external.com">External</a>
	<form><input type="text"><input type="password"></form>
</body>
</html>`))
	}))
	defer srv.Close()

	svc := New(testConfig(), nil)
	resp, err := svc.Analyze(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.HTMLVersion != "HTML5" {
		t.Errorf("expected HTML5, got %s", resp.HTMLVersion)
	}
	if resp.Title != "Test Page" {
		t.Errorf("expected 'Test Page', got '%s'", resp.Title)
	}
	if resp.Headings.H1 != 1 {
		t.Errorf("expected 1 h1, got %d", resp.Headings.H1)
	}
	if resp.Headings.H2 != 2 {
		t.Errorf("expected 2 h2, got %d", resp.Headings.H2)
	}
	if resp.Links.Internal != 1 {
		t.Errorf("expected 1 internal link, got %d", resp.Links.Internal)
	}
	if resp.Links.External != 1 {
		t.Errorf("expected 1 external link, got %d", resp.Links.External)
	}
	if !resp.HasLoginForm {
		t.Error("expected HasLoginForm=true")
	}
}
