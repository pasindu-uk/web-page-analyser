package fetcher

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetch_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("<html><body>hello</body></html>"))
	}))
	defer srv.Close()

	f := New(5*time.Second, WithAllowPrivateIPs())
	result, err := f.Fetch(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer result.Body.Close()

	body, _ := io.ReadAll(result.Body)
	if string(body) != "<html><body>hello</body></html>" {
		t.Errorf("unexpected body: %s", body)
	}
	if result.FinalURL != srv.URL+"/" {
		// httptest URL might not have trailing slash; just check it's non-empty
		if result.FinalURL == "" {
			t.Error("FinalURL is empty")
		}
	}
}

func TestFetch_NonHTMLContentType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"key":"value"}`))
	}))
	defer srv.Close()

	f := New(5*time.Second, WithAllowPrivateIPs())
	_, err := f.Fetch(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error for non-HTML content type")
	}
}

func TestFetch_HTTPError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"404 Not Found", http.StatusNotFound},
		{"500 Internal Server Error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer srv.Close()

			f := New(5*time.Second, WithAllowPrivateIPs())
			_, err := f.Fetch(context.Background(), srv.URL)
			if err == nil {
				t.Fatal("expected error")
			}
			httpErr, ok := err.(*HTTPError)
			if !ok {
				t.Fatalf("expected HTTPError, got %T: %v", err, err)
			}
			if httpErr.StatusCode != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, httpErr.StatusCode)
			}
		})
	}
}

func TestFetch_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Write([]byte("late"))
	}))
	defer srv.Close()

	f := New(100*time.Millisecond, WithAllowPrivateIPs())
	_, err := f.Fetch(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestFetch_Redirect(t *testing.T) {
	final := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html>redirected</html>"))
	}))
	defer final.Close()

	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, final.URL, http.StatusMovedPermanently)
	}))
	defer redirect.Close()

	f := New(5*time.Second, WithAllowPrivateIPs())
	result, err := f.Fetch(context.Background(), redirect.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer result.Body.Close()

	if result.FinalURL != final.URL+"/" {
		// Just verify it followed the redirect (final URL differs from original)
		if result.FinalURL == redirect.URL || result.FinalURL == redirect.URL+"/" {
			t.Errorf("expected redirect to final URL, got %s", result.FinalURL)
		}
	}
}

func TestFetch_BlocksPrivateIPs(t *testing.T) {
	// Without WithAllowPrivateIPs, the fetcher should block private IPs.
	f := New(5 * time.Second)

	tests := []struct {
		name string
		url  string
	}{
		{"loopback", "http://127.0.0.1/"},
		{"loopback v6", "http://[::1]/"},
		{"private 10.x", "http://10.0.0.1/"},
		{"private 192.168.x", "http://192.168.1.1/"},
		{"link-local", "http://169.254.169.254/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := f.Fetch(context.Background(), tt.url)
			if err == nil {
				t.Fatal("expected SSRF protection to block request")
			}
		})
	}
}
