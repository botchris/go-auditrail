package networkd

import (
	"github.com/gin-gonic/gin"
)

var _ gin.HandlerFunc = GinMiddleware

// GinMiddleware is a Gin middleware that injects into request's context
// a [networkd.Details] object holding the client's IP address.
func GinMiddleware(c *gin.Context) {
	ip := c.ClientIP()
	if ip == "" {
		c.Next()

		return
	}

	d := Details{
		Client: Client{
			IP: ip,
		},
	}

	ctx := c.Request.Context()
	ctx = AddToContext(ctx, d)
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}
