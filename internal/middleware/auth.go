package middleware

import (
	"net/http"

	"github.com/ESP-ODIN/authkit-go"
	"github.com/gin-gonic/gin"
)

func RequireAuth(authenticator *authkit.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticated := false

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authenticated = true
			c.Request = r
			c.Next()
		})

		authenticator.RequireAuth(next).ServeHTTP(c.Writer, c.Request)

		if !authenticated {
			c.Abort()
		}
	}
}
