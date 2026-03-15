package analyzer

import (
	"strings"
	"testing"
)

func TestAnalyze_HTML5Doctype(t *testing.T) {
	input := `<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>`
	result, err := Analyze(strings.NewReader(input), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if result.HTMLVersion != "HTML5" {
		t.Errorf("expected HTML5, got %s", result.HTMLVersion)
	}
}

func TestAnalyze_HTML401Doctype(t *testing.T) {
	input := `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd"><html><head><title>Test</title></head><body></body></html>`
	result, err := Analyze(strings.NewReader(input), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if result.HTMLVersion != "HTML 4.01" {
		t.Errorf("expected HTML 4.01, got %s", result.HTMLVersion)
	}
}

func TestAnalyze_XHTMLDoctype(t *testing.T) {
	input := `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"><html><head><title>Test</title></head><body></body></html>`
	result, err := Analyze(strings.NewReader(input), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if result.HTMLVersion != "XHTML" {
		t.Errorf("expected XHTML, got %s", result.HTMLVersion)
	}
}

func TestAnalyze_UnknownDoctype(t *testing.T) {
	input := `<html><head><title>No doctype</title></head><body></body></html>`
	result, err := Analyze(strings.NewReader(input), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if result.HTMLVersion != "Unknown" {
		t.Errorf("expected Unknown, got %s", result.HTMLVersion)
	}
}

func TestAnalyze_Title(t *testing.T) {
	input := `<!DOCTYPE html><html><head><title>  My Page  </title></head><body></body></html>`
	result, err := Analyze(strings.NewReader(input), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if result.Title != "My Page" {
		t.Errorf("expected 'My Page', got '%s'", result.Title)
	}
}

func TestAnalyze_Headings(t *testing.T) {
	input := `<!DOCTYPE html><html><body>
		<h1>One</h1><h1>Two</h1>
		<h2>Sub</h2>
		<h3>SubSub</h3><h3>SubSub2</h3><h3>SubSub3</h3>
	</body></html>`
	result, err := Analyze(strings.NewReader(input), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]int{"h1": 2, "h2": 1, "h3": 3, "h4": 0, "h5": 0, "h6": 0}
	for tag, want := range expected {
		if got := result.Headings[tag]; got != want {
			t.Errorf("%s: expected %d, got %d", tag, want, got)
		}
	}
}

func TestAnalyze_LoginForm(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected bool
	}{
		{
			name:     "form with password",
			html:     `<html><body><form><input type="text" name="user"><input type="password" name="pass"></form></body></html>`,
			expected: true,
		},
		{
			name:     "form without password",
			html:     `<html><body><form><input type="text" name="search"></form></body></html>`,
			expected: false,
		},
		{
			name:     "password outside form",
			html:     `<html><body><input type="password" name="pass"></body></html>`,
			expected: true,
		},
		{
			name:     "no form",
			html:     `<html><body><p>Hello</p></body></html>`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Analyze(strings.NewReader(tt.html), "https://example.com")
			if err != nil {
				t.Fatal(err)
			}
			if result.HasLoginForm != tt.expected {
				t.Errorf("expected HasLoginForm=%v, got %v", tt.expected, result.HasLoginForm)
			}
		})
	}
}

func TestAnalyze_Links(t *testing.T) {
	input := `<!DOCTYPE html><html><body>
		<a href="/about">About</a>
		<a href="https://example.com/contact">Contact</a>
		<a href="https://external.com/page">External</a>
		<a href="https://other.org">Other</a>
		<a href="#section">Anchor</a>
		<a href="mailto:test@test.com">Email</a>
	</body></html>`

	result, err := Analyze(strings.NewReader(input), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	var internal, external int
	for _, link := range result.Links {
		if link.IsInternal {
			internal++
		} else {
			external++
		}
	}

	if internal != 2 {
		t.Errorf("expected 2 internal links, got %d", internal)
	}
	if external != 2 {
		t.Errorf("expected 2 external links, got %d", external)
	}
}

func TestAnalyze_RelativeLinks(t *testing.T) {
	input := `<html><body>
		<a href="/path">Absolute path</a>
		<a href="relative">Relative</a>
		<a href="../up">Up</a>
	</body></html>`

	result, err := Analyze(strings.NewReader(input), "https://example.com/page/sub")
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Links) != 3 {
		t.Fatalf("expected 3 links, got %d", len(result.Links))
	}

	for _, link := range result.Links {
		if !link.IsInternal {
			t.Errorf("expected all links to be internal, got external: %s", link.URL)
		}
		if !strings.HasPrefix(link.URL, "https://example.com") {
			t.Errorf("expected resolved URL, got: %s", link.URL)
		}
	}
}
