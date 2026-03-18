package analyzer

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type LinkChecker struct {
	maxWorkers int
	client     *http.Client
}

func NewLinkChecker(maxWorkers int, timeout time.Duration) *LinkChecker {
	return &LinkChecker{
		maxWorkers: maxWorkers,
		client: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return nil
			},
		},
	}
}

// CheckLinks checks accessibility of the given URLs concurrently and returns
// the count of inaccessible links.
func (lc *LinkChecker) CheckLinks(ctx context.Context, urls []string) int {
	if len(urls) == 0 {
		return 0
	}

	jobs := make(chan string, len(urls))
	for _, u := range urls {
		jobs <- u
	}
	close(jobs)

	var mu sync.Mutex
	inaccessible := 0

	var wg sync.WaitGroup
	workers := lc.maxWorkers
	if workers > len(urls) {
		workers = len(urls)
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for u := range jobs {
				if !lc.isAccessible(ctx, u) {
					mu.Lock()
					inaccessible++
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	return inaccessible
}

func (lc *LinkChecker) isAccessible(ctx context.Context, rawURL string) bool {
	// Try HEAD first (lightweight, no body transfer).
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, rawURL, nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", "WebPageAnalyzer/1.0")

	resp, err := lc.client.Do(req)
	if err == nil {
		resp.Body.Close()
		// Many servers block HEAD or return 405. Fall through to GET in those cases.
		if resp.StatusCode < 400 {
			return true
		}
		if resp.StatusCode != http.StatusMethodNotAllowed && resp.StatusCode != http.StatusForbidden {
			return false
		}
	}

	// Fallback to GET when HEAD fails or is rejected.
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", "WebPageAnalyzer/1.0")

	resp, err = lc.client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()

	return resp.StatusCode < 400
}
