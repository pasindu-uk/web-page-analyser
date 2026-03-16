package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/pasindu/web-page-analyser/internal/model"
)

// mockRepo is a minimal Repository implementation for testing the cache.
type mockRepo struct {
	listCalls atomic.Int32
	data      []model.AnalyzeResponse
}

func (m *mockRepo) List(_ context.Context) ([]model.AnalyzeResponse, error) {
	m.listCalls.Add(1)
	return m.data, nil
}

func (m *mockRepo) Save(_ context.Context, resp *model.AnalyzeResponse) error {
	m.data = append(m.data, *resp)
	return nil
}

func TestCachedList_ReturnsCachedOnSecondCall(t *testing.T) {
	mock := &mockRepo{
		data: []model.AnalyzeResponse{{URL: "https://example.com", Title: "Example"}},
	}
	cached := NewCached(mock)
	ctx := context.Background()

	first, err := cached.List(ctx)
	if err != nil {
		t.Fatalf("first List: %v", err)
	}
	second, err := cached.List(ctx)
	if err != nil {
		t.Fatalf("second List: %v", err)
	}

	if mock.listCalls.Load() != 1 {
		t.Errorf("expected inner List called once, got %d", mock.listCalls.Load())
	}
	if len(first) != 1 || len(second) != 1 {
		t.Errorf("expected 1 result each, got %d and %d", len(first), len(second))
	}
}

func TestCachedSave_InvalidatesCache(t *testing.T) {
	mock := &mockRepo{
		data: []model.AnalyzeResponse{{URL: "https://example.com"}},
	}
	cached := NewCached(mock)
	ctx := context.Background()

	// Populate cache.
	if _, err := cached.List(ctx); err != nil {
		t.Fatalf("List: %v", err)
	}
	if mock.listCalls.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", mock.listCalls.Load())
	}

	// Save should invalidate.
	if err := cached.Save(ctx, &model.AnalyzeResponse{URL: "https://new.com"}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Next List should hit inner again.
	results, err := cached.List(ctx)
	if err != nil {
		t.Fatalf("List after Save: %v", err)
	}
	if mock.listCalls.Load() != 2 {
		t.Errorf("expected 2 inner List calls, got %d", mock.listCalls.Load())
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestCachedInvalidate_ForcesReQuery(t *testing.T) {
	mock := &mockRepo{
		data: []model.AnalyzeResponse{{URL: "https://example.com"}},
	}
	cached := NewCached(mock)
	ctx := context.Background()

	// Populate cache.
	if _, err := cached.List(ctx); err != nil {
		t.Fatalf("List: %v", err)
	}
	if mock.listCalls.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", mock.listCalls.Load())
	}

	// Invalidate should mark cache stale.
	cached.Invalidate()

	// Next List should hit inner again.
	if _, err := cached.List(ctx); err != nil {
		t.Fatalf("List after Invalidate: %v", err)
	}
	if mock.listCalls.Load() != 2 {
		t.Errorf("expected 2 inner List calls after Invalidate, got %d", mock.listCalls.Load())
	}
}

func TestCachedList_ConcurrentReads(t *testing.T) {
	mock := &mockRepo{
		data: []model.AnalyzeResponse{{URL: "https://example.com"}},
	}
	cached := NewCached(mock)
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := cached.List(ctx); err != nil {
				t.Errorf("concurrent List: %v", err)
			}
		}()
	}
	wg.Wait()
}
