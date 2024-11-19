package web

import (
	"app/go-sso/internal/entity"
	apiRequest "app/go-sso/internal/http/request/user"
	webRequest "app/go-sso/internal/http/request/web/user"
	usecase "app/go-sso/internal/usecase/user"
	"app/go-sso/views"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	Log      *log.Logger
	Validate *validator.Validate
}

type AuthHandlerInterface interface {
	LoginView(ctx *gin.Context)
	Login(ctx *gin.Context)
}

func AuthHandlerFactory(log *log.Logger, validator *validator.Validate) AuthHandlerInterface {
	return &AuthHandler{
		Log:      log,
		Validate: validator,
	}
}

func (h *AuthHandler) LoginView(ctx *gin.Context) {
	login := views.NewView("auth_base", "views/auth/login.html")
	data := map[string]interface{}{
		"Title": "Go SSO | Login",
	}
	login.Render(ctx, data)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	session := sessions.Default(ctx)
	payload := new(webRequest.LoginWebRequest)
	if err := ctx.ShouldBind(payload); err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}
	err := h.Validate.Struct(payload)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}
	usecase := usecase.LoginUseCaseFactory(h.Log)
	response, err := usecase.Login(apiRequest.LoginRequest{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}
	var profile = entity.User{
		ID:       response.User.ID,
		Email:    response.User.Email,
		Name:     response.User.Name,
		Username: response.User.Username,
	}
	session.Set("profile", profile)
	session.Delete("error") // Clear any previous error messages
	if err := session.Save(); err != nil {
		h.Log.Printf("Session save error: %v", err)
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}
	ctx.Redirect(302, "/")
}
