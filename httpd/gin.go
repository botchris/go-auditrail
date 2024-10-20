package httpd

import (
	"github.com/gin-gonic/gin"
)

var _ gin.HandlerFunc = GinMiddleware

// GinMiddleware is a Gin middleware that injects into request's context
// a [httpd.Details] object holding request details.
func GinMiddleware(c *gin.Context) {
	status := "200"

	// this may not work when the response is not yet written.
	if c.Request.Response != nil {
		status = c.Request.Response.Status
	}

	d := Details{
		Method:     c.Request.Method,
		StatusCode: status,
		UserAgent:  c.Request.UserAgent(),
		URL: URL{
			Host: c.Request.URL.Host,
			Path: c.Request.URL.Path,
		},
	}

	ctx := c.Request.Context()
	ctx = AddToContext(ctx, d)
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}
