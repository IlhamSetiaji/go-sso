package middleware

import (
	"app/go-sso/internal/entity"
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

func EmailVerifiedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		profile := session.Get("profile")

		if profile == nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		userProfile, ok := profile.(entity.Profile)
		if !ok {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		if userProfile.EmailVerifiedAt.IsZero() {
			c.Redirect(http.StatusSeeOther, "/otp")
			c.Abort()
			return
		}

		c.Next()
	}
}
