package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pasindu/web-page-analyser/internal/analyzer"
	"github.com/pasindu/web-page-analyser/internal/config"
	"github.com/pasindu/web-page-analyser/internal/fetcher"
	"github.com/pasindu/web-page-analyser/internal/model"
)

type AnalyzeService struct {
	fetcher     *fetcher.Fetcher
	linkChecker *analyzer.LinkChecker
}

func New(cfg *config.Config) *AnalyzeService {
	return &AnalyzeService{
		fetcher:     fetcher.New(cfg.RequestTimeout),
		linkChecker: analyzer.NewLinkChecker(cfg.MaxLinkCheckWorkers, cfg.MaxLinksToCheck, cfg.RequestTimeout),
	}
}

func (s *AnalyzeService) Analyze(ctx context.Context, rawURL string) (*model.AnalyzeResponse, error) {
	if err := validateURL(rawURL); err != nil {
		return nil, &ValidationError{Message: err.Error()}
	}

	result, err := s.fetcher.Fetch(ctx, rawURL)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	analysis, err := analyzer.Analyze(result.Body, result.FinalURL)
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	var allURLs []string
	for _, link := range analysis.Links {
		allURLs = append(allURLs, link.URL)
	}

	inaccessible := s.linkChecker.CheckLinks(ctx, allURLs)

	var internal, external int
	for _, link := range analysis.Links {
		if link.IsInternal {
			internal++
		} else {
			external++
		}
	}

	return &model.AnalyzeResponse{
		URL:         rawURL,
		HTMLVersion: analysis.HTMLVersion,
		Title:       analysis.Title,
		Headings: model.HeadingCount{
			H1: analysis.Headings["h1"],
			H2: analysis.Headings["h2"],
			H3: analysis.Headings["h3"],
			H4: analysis.Headings["h4"],
			H5: analysis.Headings["h5"],
			H6: analysis.Headings["h6"],
		},
		Links: model.LinkSummary{
			Internal:     internal,
			External:     external,
			Inaccessible: inaccessible,
		},
		HasLoginForm: analysis.HasLoginForm,
	}, nil
}

func validateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL is required")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme")
	}

	if u.Host == "" {
		return fmt.Errorf("URL must include a host")
	}

	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
