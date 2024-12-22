package utils

import (
	"encoding/base64"
	"fmt"

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
	encodedValue := base64.StdEncoding.EncodeToString([]byte(token))
	ctx.SetCookie(
		opts.Name,
		encodedValue,
		opts.MaxAge,
		opts.Path,
		opts.Domain,
		opts.Secure,
		opts.HTTPOnly,
	)
	fmt.Println("Cookie set")
}

func GetTokenFromCookie(ctx *gin.Context, cookieName string) (string, error) {
	token, err := ctx.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	decodedValue, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		// Handle error
	}
	return string(decodedValue), nil
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
