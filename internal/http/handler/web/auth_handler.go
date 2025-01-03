package web

import (
	"app/go-sso/internal/entity"
	webRequest "app/go-sso/internal/http/request/web/user"
	appUsecase "app/go-sso/internal/usecase/application"
	usecase "app/go-sso/internal/usecase/user"
	"app/go-sso/utils"
	"app/go-sso/views"
	"fmt"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	ChooseRoles(ctx *gin.Context)
	ContinueLogin(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	CheckCookieTest(ctx *gin.Context)
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
	state := ctx.Query("state")
	login := views.NewView("auth_base", "views/auth/login.html")
	data := map[string]interface{}{
		"Title": "Go SSO | Login",
	}

	if state != "" {
		data["State"] = state
	}

	login.Render(ctx, data)
}

func (h *AuthHandler) ChooseRoles(ctx *gin.Context) {
	state := ctx.Query("state")
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	if profile == nil {
		session.Set("error", "Profile not found")
		session.Save()
		h.Log.Printf("Profile not found")
		ctx.Redirect(302, "/logout")
		return
	}

	userProfile, ok := profile.(entity.Profile)
	if !ok {
		session.Set("error", "Profile not found")
		session.Save()
		h.Log.Printf("Profile not found")
		ctx.Redirect(302, "/logout")
		return
	}

	factory := usecase.FindByIdUseCaseFactory(h.Log)
	response, err := factory.Execute(&usecase.IFindByIdUseCaseRequest{
		ID: userProfile.ID,
	})
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, "/logout")
		return
	}

	if response.User == nil {
		session.Set("error", "User not found")
		session.Save()
		h.Log.Printf("User not found")
	}

	viewRoles := views.NewView("auth_base", "views/auth/choose_roles.html")
	data := map[string]interface{}{
		"Title": "Go SSO | Choose Roles",
		"Roles": response.User.Roles,
	}

	if state != "" {
		data["State"] = state
	}

	viewRoles.Render(ctx, data)
}

func (h *AuthHandler) ContinueLogin(ctx *gin.Context) {
	testCookie := utils.NewDefaultCookieOptions("test_haha")
	testCookie.Domain = h.Config.GetString("app.domain")
	utils.SetTokenCookie(ctx, "test_cookie", testCookie)

	session := sessions.Default(ctx)
	payload := new(webRequest.ChooseRolesWebRequest)
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

	factory := usecase.FindByIdOnlyUseCaseFactory(h.Log)
	response, err := factory.Execute(&usecase.IFindByIdOnlyUseCaseRequest{
		ID: uuid.MustParse(payload.ID),
	})

	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Printf(err.Error())
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	filteredRoles := []entity.Role{}
	for _, role := range response.User.Roles {
		if role.ID.String() == payload.RoleID {
			filteredRoles = append(filteredRoles, role)
			break
		}
	}
	response.User.Roles = filteredRoles

	token, err := utils.GenerateToken(response.User)
	if err != nil {
		h.Log.Errorf("Error when generating token: %v", err)
		session.Set("error", err.Error())
		session.Save()
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	jwtCookie := utils.NewDefaultCookieOptions("jwt_token")
	jwtCookie.Domain = h.Config.GetString("app.domain")
	utils.SetTokenCookie(ctx, token, jwtCookie)

	if payload.State != "" {
		data, err := h.loginAsApplication(token, payload.State, response.User)
		if err != nil {
			session.Set("error", err.Error())
			session.Save()
			h.Log.Printf(err.Error())
			ctx.Redirect(302, ctx.Request.Referer())
			return
		}

		application := data["application"].(*entity.Application)

		redirectURL := fmt.Sprintf("%s?token=%s", application.RedirectURI, data["token"])
		h.Log.Printf("Redirecting to URL: %s", redirectURL)

		if !strings.HasPrefix(redirectURL, "http") {
			redirectURL = "http://" + redirectURL
		}

		ctx.Redirect(302, redirectURL)
		return
	}

	ctx.Redirect(302, "/portal")

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

	ctx.Redirect(302, "/choose-roles")
	return

	// token, err := utils.GenerateToken(&response.User)
	// if err != nil {
	// 	h.Log.Errorf("Error when generating token: %v", err)
	// 	session.Set("error", err.Error())
	// 	session.Save()
	// 	ctx.Redirect(302, ctx.Request.Referer())
	// 	return
	// }

	// if payload.State != "" {
	// 	data, err := h.loginAsApplication(token, payload.State, &response.User)
	// 	if err != nil {
	// 		session.Set("error", err.Error())
	// 		session.Save()
	// 		h.Log.Printf(err.Error())
	// 		ctx.Redirect(302, ctx.Request.Referer())
	// 		return
	// 	}

	// 	application := data["application"].(*entity.Application)

	// 	redirectURL := fmt.Sprintf("%s?token=%s", application.RedirectURI, data["token"])
	// 	h.Log.Printf("Redirecting to URL: %s", redirectURL)

	// 	if !strings.HasPrefix(redirectURL, "http") {
	// 		redirectURL = "http://" + redirectURL
	// 	}

	// 	ctx.Redirect(302, redirectURL)
	// 	return
	// }

	// if !h.checkUserRole(&response.User, "superadmin") {
	// 	session.Set("error", "You are not allowed to access this page")
	// 	session.Save()
	// 	ctx.Redirect(302, ctx.Request.Referer())
	// 	return
	// }

	// jwtCookie := utils.NewDefaultCookieOptions("jwt_token")
	// jwtCookie.Domain = h.Config.GetString("app.domain")
	// utils.SetTokenCookie(ctx, token, jwtCookie)
	// utils.SuccessResponse(ctx, 200, "success", response.User)
	// return
	// ctx.Redirect(302, "/portal")
}

func (h *AuthHandler) CheckCookieTest(ctx *gin.Context) {
	cookie, err := utils.GetTokenFromCookie(ctx, "test_haha")
	if err != nil {
		ctx.JSON(200, gin.H{
			"message": "Cookie not found",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": cookie,
	})
}

func (h *AuthHandler) checkUserRole(user *entity.User, role string) bool {
	for _, r := range user.Roles {
		if r.Name == role {
			return true
		}
	}
	return false
}

func (h *AuthHandler) loginAsApplication(token string, state string, user *entity.User) (map[string]interface{}, error) {
	factory := appUsecase.FindApplicationByNameUsecaseFactory(h.Log)
	application, err := factory.Execute(&appUsecase.IFindApplicationByNameUsecaseRequest{
		Name: state,
	})

	if err != nil {
		h.Log.Errorf("Error when finding application: %v", err)
		return nil, err
	}

	data := map[string]interface{}{
		"token":       token,
		"application": application.Application,
		"user":        user,
	}

	return data, nil
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	session := utils.NewSession(ctx)
	session.Delete("profile")
	session.Set("success", "You have been logged out")
	session.Save()
	utils.ClearTokenCookie(ctx, "jwt_token", h.Config.GetString("app.domain"))
	ctx.Redirect(302, "/login")
}
