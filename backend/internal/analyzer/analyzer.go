package analyzer

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Result struct {
	HTMLVersion  string
	Title        string
	Headings     map[string]int
	Links        []Link
	HasLoginForm bool
}

type Link struct {
	URL        string
	IsInternal bool
}

func Analyze(body io.Reader, pageURL string) (*Result, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	pageU, _ := url.Parse(pageURL)

	result := &Result{
		Headings: map[string]int{
			"h1": 0, "h2": 0, "h3": 0,
			"h4": 0, "h5": 0, "h6": 0,
		},
	}

	var inTitle bool
	var titleBuilder strings.Builder
	var inForm bool

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		switch n.Type {
		case html.DoctypeNode:
			result.HTMLVersion = detectHTMLVersion(n)
		case html.ElementNode:
			tag := strings.ToLower(n.Data)

			if tag == "title" {
				inTitle = true
			}

			if _, ok := result.Headings[tag]; ok {
				result.Headings[tag]++
			}

			if tag == "a" {
				if href := getAttr(n, "href"); href != "" {
					link := resolveLink(href, pageU)
					if link != nil {
						result.Links = append(result.Links, *link)
					}
				}
			}

			if tag == "form" {
				inForm = true
				defer func() { inForm = false }()
			}

			if inForm && tag == "input" && getAttr(n, "type") == "password" {
				result.HasLoginForm = true
			}
		case html.TextNode:
			if inTitle {
				titleBuilder.WriteString(n.Data)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}

		if n.Type == html.ElementNode && strings.ToLower(n.Data) == "title" {
			inTitle = false
			result.Title = strings.TrimSpace(titleBuilder.String())
		}
	}

	walk(doc)

	if result.HTMLVersion == "" {
		result.HTMLVersion = "Unknown"
	}

	return result, nil
}

func detectHTMLVersion(n *html.Node) string {
	if n.Type != html.DoctypeNode {
		return "Unknown"
	}

	// HTML5: <!DOCTYPE html> with no public/system identifiers
	if strings.EqualFold(n.Data, "html") && n.Attr == nil {
		return "HTML5"
	}

	for _, attr := range n.Attr {
		val := strings.ToLower(attr.Val)
		if strings.Contains(val, "xhtml") {
			return "XHTML"
		}
		if strings.Contains(val, "html 4.01") {
			return "HTML 4.01"
		}
		if strings.Contains(val, "html 4.0") {
			return "HTML 4.0"
		}
	}

	return "Unknown"
}

func resolveLink(href string, pageURL *url.URL) *Link {
	href = strings.TrimSpace(href)
	if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "javascript:") || strings.HasPrefix(href, "mailto:") {
		return nil
	}

	parsed, err := url.Parse(href)
	if err != nil {
		return nil
	}

	resolved := pageURL.ResolveReference(parsed)

	isInternal := resolved.Host == pageURL.Host
	return &Link{
		URL:        resolved.String(),
		IsInternal: isInternal,
	}
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, key) {
			return a.Val
		}
	}
	return ""
}
