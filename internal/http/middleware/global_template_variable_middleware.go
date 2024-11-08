package middleware

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GlobalTemplateVariablesMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		c.Set("Title", "Go SSO")
		c.Set("AssetBase", "/assets")
		c.Set("Errors", session.Flashes("errors"))
		c.Set("Success", session.Get("success"))
		c.Set("Error", session.Get("error"))
		c.Set("Info", session.Get("info"))
		c.Set("Warning", session.Get("warning"))
		session.Save()
		fmt.Printf("session: %v\n", session)
		c.Next()
	}
}
