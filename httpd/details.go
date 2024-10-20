package httpd

// Details represents the details of an HTTP request.
type Details struct {
	Method     string `json:"method"`
	StatusCode string `json:"status_code"`
	UserAgent  string `json:"user_agent"`
	URL        URL    `json:"url"`
}

// IsEmpty returns true if the details are empty.
func (d Details) IsEmpty() bool {
	return d == Details{}
}

// URL HTTP request url details.
type URL struct {
	Host string `json:"host"`
	Path string `json:"path"`
}
