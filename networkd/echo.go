package networkd

import "github.com/labstack/echo/v4"

var _ echo.MiddlewareFunc = EchoMiddleware

// EchoMiddleware is an Echo middleware that injects into request's context
// a [networkd.Details] object holding the client's IP address.
func EchoMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()
		if ip == "" {
			return next(c)
		}

		d := Details{
			Client: Client{
				IP: ip,
			},
		}

		ctx := c.Request().Context()
		ctx = AddToContext(ctx, d)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
