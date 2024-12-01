package utils

import (
	"github.com/gin-gonic/gin"
)

const (
	DefaultCookieMaxAge = 3600 * 24 // 24 hours in seconds
)

type CookieOptions struct {
	Name     string
	Domain   string
	MaxAge   int
	Path     string
	Secure   bool
	HTTPOnly bool
}

func NewDefaultCookieOptions(name string) CookieOptions {
	return CookieOptions{
		Name:     name,
		MaxAge:   DefaultCookieMaxAge,
		Path:     "/",
		Secure:   false, // set to true if in production
		HTTPOnly: true,
	}
}

func SetTokenCookie(ctx *gin.Context, token string, opts CookieOptions) {
	ctx.SetCookie(
		opts.Name,
		token,
		opts.MaxAge,
		opts.Path,
		opts.Domain,
		opts.Secure,
		opts.HTTPOnly,
	)
}

func GetTokenFromCookie(ctx *gin.Context, cookieName string) (string, error) {
	token, err := ctx.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ClearTokenCookie(ctx *gin.Context, cookieName, domain string) {
	ctx.SetCookie(
		cookieName,
		"",
		-1,
		"/",
		domain,
		true,
		true,
	)
}
