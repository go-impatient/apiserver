package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/moocss/apiserver/src/pkg/version"
)

// VersionMiddleware : add version on header.
func VersionMiddleware() gin.HandlerFunc {
	// Set out header value for each response
	return func(c *gin.Context) {
		c.Header("X-APISERSION-VERSION", version.GetVersion())
		c.Next()
	}
}
