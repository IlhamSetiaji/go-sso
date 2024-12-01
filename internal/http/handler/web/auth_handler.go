package web

import (
	"app/go-sso/internal/entity"
	webRequest "app/go-sso/internal/http/request/web/user"
	usecase "app/go-sso/internal/usecase/user"
	"app/go-sso/utils"
	"app/go-sso/views"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AuthHandler struct {
	Config   *viper.Viper
	Log      *logrus.Logger
	Validate *validator.Validate
}

type AuthHandlerInterface interface {
	LoginView(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

func AuthHandlerFactory(log *logrus.Logger, validator *validator.Validate) AuthHandlerInterface {
	config := viper.New()
	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	return &AuthHandler{
		Config:   config,
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
	factory := usecase.LoginUseCaseFactory(h.Log)
	response, err := factory.Execute(usecase.ILoginUseCaseRequest{
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
	var profile = entity.Profile{
		ID:       response.User.ID,
		Name:     response.User.Name,
		Email:    response.User.Email,
		Username: response.User.Username,
	}

	session.Set("profile", profile)
	session.Delete("error")
	if err := session.Save(); err != nil {
		h.Log.Printf("[Auth handler] Session save error: %v", err)
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}
	ctx.Redirect(302, "/")
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	session := utils.NewSession(ctx)
	session.Delete("profile")
	session.Set("success", "You have been logged out")
	session.Save()
	ctx.Redirect(302, "/login")
}
