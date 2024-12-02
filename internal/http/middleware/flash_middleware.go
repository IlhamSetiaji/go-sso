package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func FlashMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		success := session.Get("success")
		error := session.Get("error")
		status := session.Get("status")
		errors := session.Get("errors")
		warning := session.Get("warning")

		if success != nil {
			c.Set("success", success)
			session.Delete("success")
		}
		if error != nil {
			c.Set("error", error)
			session.Delete("error")
		}
		if status != nil {
			c.Set("status", status)
			session.Delete("status")
		}
		if errors != nil {
			c.Set("errors", errors)
			session.Delete("errors")
		}
		if warning != nil {
			c.Set("warning", warning)
			session.Delete("warning")
		}

		session.Save()

		c.Next()
	}
}
