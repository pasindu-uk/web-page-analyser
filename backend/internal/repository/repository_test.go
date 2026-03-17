package repository

import (
	"context"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/pasindu-uk/web-page-analyser/internal/model"
)

func setupTestDB(t *testing.T) *MySQLRepository {
	t.Helper()
	dsn := os.Getenv("MYSQL_DSN_TEST")
	if dsn == "" {
		t.Skip("MYSQL_DSN_TEST not set, skipping integration test")
	}

	repo, err := NewMySQL(dsn)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	if err := RunMigrations(repo.DB()); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Clean table before each test
	if _, err := repo.DB().Exec("DELETE FROM analyses"); err != nil {
		t.Fatalf("failed to clean table: %v", err)
	}

	t.Cleanup(func() { repo.Close() })
	return repo
}

func TestMySQLRepository_SaveAndList(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	sample := &model.AnalyzeResponse{
		URL:         "https://example.com",
		HTMLVersion: "HTML5",
		Title:       "Example",
		Headings:    model.HeadingCount{H1: 1, H2: 3},
		Links:       model.LinkSummary{Internal: 5, External: 10, Inaccessible: 2},
		HasLoginForm: false,
	}

	if err := repo.Save(ctx, sample); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	results, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	got := results[0]
	if got.URL != sample.URL {
		t.Errorf("URL: expected %q, got %q", sample.URL, got.URL)
	}
	if got.HTMLVersion != sample.HTMLVersion {
		t.Errorf("HTMLVersion: expected %q, got %q", sample.HTMLVersion, got.HTMLVersion)
	}
	if got.Title != sample.Title {
		t.Errorf("Title: expected %q, got %q", sample.Title, got.Title)
	}
	if got.Headings != sample.Headings {
		t.Errorf("Headings: expected %+v, got %+v", sample.Headings, got.Headings)
	}
	if got.Links != sample.Links {
		t.Errorf("Links: expected %+v, got %+v", sample.Links, got.Links)
	}
	if got.HasLoginForm != sample.HasLoginForm {
		t.Errorf("HasLoginForm: expected %v, got %v", sample.HasLoginForm, got.HasLoginForm)
	}
}

func TestMySQLRepository_ListEmpty(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	results, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestMySQLRepository_ListOrder(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	first := &model.AnalyzeResponse{URL: "https://first.com", HTMLVersion: "HTML5", Title: "First"}
	second := &model.AnalyzeResponse{URL: "https://second.com", HTMLVersion: "HTML5", Title: "Second"}

	if err := repo.Save(ctx, first); err != nil {
		t.Fatalf("Save first failed: %v", err)
	}
	if err := repo.Save(ctx, second); err != nil {
		t.Fatalf("Save second failed: %v", err)
	}

	results, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// Most recent first (DESC order)
	if results[0].URL != "https://second.com" {
		t.Errorf("expected second.com first, got %s", results[0].URL)
	}
	if results[1].URL != "https://first.com" {
		t.Errorf("expected first.com second, got %s", results[1].URL)
	}
}
