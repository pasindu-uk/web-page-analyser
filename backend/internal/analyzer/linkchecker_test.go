package analyzer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestCheckLinks_AllAccessible(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	lc := NewLinkChecker(3, 50, 5*time.Second)
	urls := []string{srv.URL + "/a", srv.URL + "/b", srv.URL + "/c"}
	inaccessible := lc.CheckLinks(context.Background(), urls)
	if inaccessible != 0 {
		t.Errorf("expected 0 inaccessible, got %d", inaccessible)
	}
}

func TestCheckLinks_SomeInaccessible(t *testing.T) {
	good := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer good.Close()

	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer bad.Close()

	lc := NewLinkChecker(3, 50, 5*time.Second)
	urls := []string{good.URL, bad.URL, bad.URL}
	inaccessible := lc.CheckLinks(context.Background(), urls)
	if inaccessible != 2 {
		t.Errorf("expected 2 inaccessible, got %d", inaccessible)
	}
}

func TestCheckLinks_EmptyList(t *testing.T) {
	lc := NewLinkChecker(3, 50, 5*time.Second)
	inaccessible := lc.CheckLinks(context.Background(), nil)
	if inaccessible != 0 {
		t.Errorf("expected 0, got %d", inaccessible)
	}
}

func TestCheckLinks_MaxLinksCap(t *testing.T) {
	var callCount atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	lc := NewLinkChecker(2, 3, 5*time.Second)
	urls := []string{srv.URL + "/1", srv.URL + "/2", srv.URL + "/3", srv.URL + "/4", srv.URL + "/5"}
	lc.CheckLinks(context.Background(), urls)

	if callCount.Load() != 3 {
		t.Errorf("expected 3 calls (maxLinks=3), got %d", callCount.Load())
	}
}

func TestCheckLinks_UnreachableHost(t *testing.T) {
	lc := NewLinkChecker(2, 50, 500*time.Millisecond)
	urls := []string{"http://192.0.2.1:12345/unreachable"}
	inaccessible := lc.CheckLinks(context.Background(), urls)
	if inaccessible != 1 {
		t.Errorf("expected 1 inaccessible, got %d", inaccessible)
	}
}
