package middleware

import (
	"app/go-sso/internal/entity"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func WebAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("WebAuthMiddleware checking:", c.Request.URL.Path)

		if strings.HasPrefix(c.Request.URL.Path, "/api") || strings.HasPrefix(c.Request.URL.Path, "/oauth") {
			// Don't apply web auth middleware to API or OAuth paths
			c.Next()
			return
		}

		session := sessions.Default(c)
		profile := session.Get("profile")

		if profile == nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		_, ok := profile.(entity.Profile)
		if !ok {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
