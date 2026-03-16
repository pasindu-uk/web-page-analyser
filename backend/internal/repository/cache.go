package repository

import (
	"context"
	"sync"

	"github.com/pasindu/web-page-analyser/internal/model"
)

// CachedRepository wraps a Repository and caches List results in memory.
// The cache is invalidated whenever Save inserts a new analysis.
type CachedRepository struct {
	inner  Repository
	mu     sync.RWMutex
	cached []model.AnalyzeResponse
	valid  bool
}

// NewCached returns a CachedRepository that wraps inner.
func NewCached(inner Repository) *CachedRepository {
	return &CachedRepository{inner: inner}
}

func (c *CachedRepository) List(ctx context.Context) ([]model.AnalyzeResponse, error) {
	c.mu.RLock()
	if c.valid {
		result := c.cached
		c.mu.RUnlock()
		return result, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock.
	if c.valid {
		return c.cached, nil
	}

	results, err := c.inner.List(ctx)
	if err != nil {
		return nil, err
	}
	c.cached = results
	c.valid = true
	return results, nil
}

// Invalidate marks the cache as stale so the next List call re-queries the inner repository.
func (c *CachedRepository) Invalidate() {
	c.mu.Lock()
	c.valid = false
	c.mu.Unlock()
}

func (c *CachedRepository) Save(ctx context.Context, resp *model.AnalyzeResponse) error {
	if err := c.inner.Save(ctx, resp); err != nil {
		return err
	}
	c.mu.Lock()
	c.valid = false
	c.mu.Unlock()
	return nil
}
