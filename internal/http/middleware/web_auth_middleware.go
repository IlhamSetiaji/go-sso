package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// func WebAuthMiddleware(ctx *gin.Context) {
// 	if sessions.Default(ctx).Get("profile") == nil {
// 		ctx.Redirect(http.StatusSeeOther, "/")
// 	} else {
// 		ctx.Next()
// 	}
// }

func WebAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if sessions.Default(c).Get("profile") == nil {
			c.Redirect(http.StatusSeeOther, "/login")
		} else {
			c.Next()
		}
	}
}
