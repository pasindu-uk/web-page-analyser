package fetcher

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type Result struct {
	Body     io.ReadCloser
	FinalURL string
}

// Option configures a Fetcher.
type Option func(*Fetcher)

// WithAllowPrivateIPs disables the SSRF protection that blocks requests to
// private/loopback IP addresses. Use this only in tests with httptest servers.
func WithAllowPrivateIPs() Option {
	return func(f *Fetcher) {
		f.client.Transport = nil // use default transport
	}
}

type Fetcher struct {
	client *http.Client
}

func New(timeout time.Duration, opts ...Option) *Fetcher {
	f := &Fetcher{
		client: &http.Client{
			Timeout:   timeout,
			Transport: &http.Transport{DialContext: safeDialContext},
		},
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// isPrivateIP returns true if the IP is loopback, private, link-local, or unspecified.
func isPrivateIP(ip net.IP) bool {
	return ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsUnspecified()
}

// safeDialContext resolves the address and rejects connections to private IPs.
func safeDialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid address %q: %w", addr, err)
	}

	ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("DNS lookup failed for %q: %w", host, err)
	}

	for _, ipAddr := range ips {
		if isPrivateIP(ipAddr.IP) {
			return nil, fmt.Errorf("request to private IP %s is blocked (SSRF protection)", ipAddr.IP)
		}
	}

	var d net.Dialer
	return d.DialContext(ctx, network, net.JoinHostPort(host, port))
}

func (f *Fetcher) Fetch(ctx context.Context, url string) (*Result, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	req.Header.Set("User-Agent", "WebPageAnalyzer/1.0")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, &HTTPError{StatusCode: resp.StatusCode}
	}

	contentType := resp.Header.Get("Content-Type")
	if !isHTML(contentType) {
		resp.Body.Close()
		return nil, fmt.Errorf("response is not HTML (Content-Type: %s)", contentType)
	}

	const maxBodySize = 10 << 20 // 10 MB
	limitedBody := struct {
		io.Reader
		io.Closer
	}{
		Reader: io.LimitReader(resp.Body, maxBodySize),
		Closer: resp.Body,
	}

	return &Result{
		Body:     limitedBody,
		FinalURL: resp.Request.URL.String(),
	}, nil
}

func isHTML(contentType string) bool {
	ct := strings.ToLower(contentType)
	return strings.Contains(ct, "text/html") || strings.Contains(ct, "application/xhtml+xml")
}

type HTTPError struct {
	StatusCode int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, http.StatusText(e.StatusCode))
}
