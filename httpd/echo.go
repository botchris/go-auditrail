package httpd

import (
	"github.com/labstack/echo/v4"
)

var _ echo.MiddlewareFunc = EchoMiddleware

// EchoMiddleware is an Echo middleware that injects into request's context
// a [httpd.Details] object holding request details.
func EchoMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		status := "unknown"
		req := c.Request()

		// this may not work when the response is not yet written.
		if req.Response != nil {
			status = req.Response.Status
		}

		d := Details{
			Method:     req.Method,
			StatusCode: status,
			UserAgent:  req.UserAgent(),
			URL: URL{
				Host: req.URL.Host,
				Path: req.URL.Path,
			},
		}

		ctx := AddToContext(req.Context(), d)

		c.SetRequest(req.WithContext(ctx))

		return next(c)
	}
}
