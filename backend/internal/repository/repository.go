package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/pasindu-uk/web-page-analyser/internal/model"
)

// Repository defines persistence operations for analysis results.
type Repository interface {
	Save(ctx context.Context, resp *model.AnalyzeResponse) error
	List(ctx context.Context) ([]model.AnalyzeResponse, error)
}

// MySQLRepository implements Repository using a MySQL database.
type MySQLRepository struct {
	db *sql.DB
}

// NewMySQL wraps an existing *sql.DB connection.
// The caller owns the connection and is responsible for closing it.
func NewMySQL(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Save(ctx context.Context, resp *model.AnalyzeResponse) error {
	query := `INSERT INTO analyses
		(url, html_version, title, h1_count, h2_count, h3_count, h4_count, h5_count, h6_count,
		 internal_links, external_links, inaccessible_links, has_login_form)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		resp.URL, resp.HTMLVersion, resp.Title,
		resp.Headings.H1, resp.Headings.H2, resp.Headings.H3,
		resp.Headings.H4, resp.Headings.H5, resp.Headings.H6,
		resp.Links.Internal, resp.Links.External, resp.Links.Inaccessible,
		resp.HasLoginForm,
	)
	if err != nil {
		return fmt.Errorf("inserting analysis: %w", err)
	}
	return nil
}

func (r *MySQLRepository) List(ctx context.Context) ([]model.AnalyzeResponse, error) {
	query := `SELECT url, html_version, title,
		h1_count, h2_count, h3_count, h4_count, h5_count, h6_count,
		internal_links, external_links, inaccessible_links, has_login_form
		FROM analyses ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying analyses: %w", err)
	}
	defer rows.Close()

	var results []model.AnalyzeResponse
	for rows.Next() {
		var resp model.AnalyzeResponse
		if err := rows.Scan(
			&resp.URL, &resp.HTMLVersion, &resp.Title,
			&resp.Headings.H1, &resp.Headings.H2, &resp.Headings.H3,
			&resp.Headings.H4, &resp.Headings.H5, &resp.Headings.H6,
			&resp.Links.Internal, &resp.Links.External, &resp.Links.Inaccessible,
			&resp.HasLoginForm,
		); err != nil {
			return nil, fmt.Errorf("scanning analysis row: %w", err)
		}
		results = append(results, resp)
	}
	return results, rows.Err()
}
