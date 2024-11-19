package utils

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func NewSession(ctx *gin.Context) sessions.Session {
	session := sessions.Default(ctx)
	session.Delete("error")
	session.Delete("success")
	session.Delete("warning")
	session.Delete("status")
	session.Delete("errors")
	session.Save()
	return session
}
