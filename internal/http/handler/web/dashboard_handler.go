package web

import (
	"app/go-sso/internal/http/middleware"
	usecase "app/go-sso/internal/usecase/application"
	"app/go-sso/utils"
	"app/go-sso/views"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DashboardHandler struct {
	Config   *viper.Viper
	Log      *logrus.Logger
	validate *validator.Validate
}

type DashboardHandlerInterface interface {
	Index(ctx *gin.Context)
	Portal(ctx *gin.Context)
}

func DashboardHandlerFactory(log *logrus.Logger, validate *validator.Validate) DashboardHandlerInterface {
	config := viper.New()
	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	return &DashboardHandler{
		Config:   config,
		Log:      log,
		validate: validate,
	}
}

func (h *DashboardHandler) Index(ctx *gin.Context) {
	middleware.RoleMiddleware("superadmin")(ctx)
	if ctx.IsAborted() {
		ctx.Redirect(302, "/logout")
		return
	}
	index := views.NewView("base", "views/index.html")
	data := map[string]interface{}{
		"Title": "Go SSO | Dashboard",
	}
	index.Render(ctx, data)
}

func (h *DashboardHandler) Portal(ctx *gin.Context) {
	middleware.ExceptRoleMiddleware("Applicant")(ctx)
	if ctx.IsAborted() {
		ctx.Redirect(302, "/logout")
		return
	}
	session := sessions.Default(ctx)
	// token := ctx.Query("token")
	token, err := utils.GetTokenFromCookie(ctx, "jwt_token")
	if err != nil {
		ctx.JSON(200, gin.H{
			"message": "Cookie not found",
		})
		return
	}

	if token == "" {
		session.Set("error", "There are no token")
		session.Save()
		ctx.Redirect(302, "/logout")
		return
	}

	isJWT, err := h.checkIfTokenIsJWT(token)
	if err != nil {
		h.Log.Error(err)
		session.Set("error", err.Error())
		session.Save()
		ctx.Redirect(302, "/login")
		return
	}

	if !isJWT {
		session.Set("error", "Invalid token")
		session.Save()
		ctx.Redirect(302, "/logout")
		return
	}

	factory := usecase.GetAllApplicationsUseCaseFactory(h.Log)

	resp, err := factory.Execute()
	if err != nil {
		h.Log.Error(err)
	}

	index := views.NewView("base", "views/portal.html")
	data := map[string]interface{}{
		"Title":        "Go SSO | Portal",
		"Applications": resp.Applications,
		"Token":        token,
	}
	index.Render(ctx, data)
}

func (h *DashboardHandler) checkIfTokenIsJWT(stateToken string) (bool, error) {
	jwtToken, err := jwt.Parse(stateToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.Config.GetString("jwt.secret")), nil
	})

	if err != nil {
		return false, err
	}

	if _, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		return true, nil
	} else {
		return false, nil
	}

}
