package model

type AnalyzeRequest struct {
	URL string `json:"url"`
}

type AnalyzeResponse struct {
	URL          string       `json:"url"`
	HTMLVersion  string       `json:"htmlVersion"`
	Title        string       `json:"title"`
	Headings     HeadingCount `json:"headings"`
	Links        LinkSummary  `json:"links"`
	HasLoginForm bool         `json:"hasLoginForm"`
}

type HeadingCount struct {
	H1 int `json:"h1"`
	H2 int `json:"h2"`
	H3 int `json:"h3"`
	H4 int `json:"h4"`
	H5 int `json:"h5"`
	H6 int `json:"h6"`
}

type LinkSummary struct {
	Internal     int `json:"internal"`
	External     int `json:"external"`
	Inaccessible int `json:"inaccessible"`
}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
