package handler

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/http/middleware"
	request "app/go-sso/internal/http/request/user"
	usecase "app/go-sso/internal/usecase/user"
	"app/go-sso/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/oauth2"
)

type UserHandler struct {
	Log         *log.Logger
	Validate    *validator.Validate
	OAuthConfig *config.Authenticator
}

func UserHandlerFactory(log *log.Logger, validator *validator.Validate, oAuthConfig *config.Authenticator) *UserHandler {
	return &UserHandler{
		Log:         log,
		Validate:    validator,
		OAuthConfig: oAuthConfig,
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	payload := new(request.LoginRequest)
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, 400, "error", err.Error())
		h.Log.Panicf("Error when binding request: %v", err)
		return
	}
	err := h.Validate.Struct(payload)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "error", err.Error())
		h.Log.Panicf("Error when validating request: %v", err)
		return
	}
	usecase := usecase.LoginUseCaseFactory(h.Log)
	response, err := usecase.Login(*payload)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Panicf("Error when login: %v", err)
		return
	}
	token, err := utils.GenerateToken(&response.User)
	if err != nil {
		h.Log.Panicf("Error when generating token: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}
	if err != nil {
		h.Log.Panicf("Error when login: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}
	var data = map[string]interface{}{
		"token":      token,
		"token_type": "Bearer",
		"user":       response.User,
	}
	utils.SuccessResponse(ctx, 200, "success", data)
}

func (h *UserHandler) Me(ctx *gin.Context) {
	user, err := middleware.GetUser(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Panicf("Error when getting user: %v", err)
		return
	}
	if user == nil {
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		h.Log.Panicf("User not found")
		return
	}

	utils.SuccessResponse(ctx, 200, "success", user)
}

func (h *UserHandler) LoginOAuth(ctx *gin.Context) {
	url := h.OAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *UserHandler) CallbackOAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	token, err := h.OAuthConfig.Exchange(ctx, code)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Panicf("Error when exchanging token: %v", err)
		return
	}
	idToken, err := h.OAuthConfig.VerifyIDToken(ctx, token)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Panicf("Error when verifying id token: %v", err)
		return
	}
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Panicf("Error when getting profile: %v", err)
		return
	}
	usecase := usecase.FindByEmailUseCaseFactory(h.Log)
	response, err := usecase.FindByEmail(profile["email"].(string))
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Panicf("Error when finding user by email: %v", err)
		return
	}
	jwtToken, err := utils.GenerateToken(response.User)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Panicf("Error when generating token: %v", err)
		return
	}
	utils.SuccessResponse(ctx, 200, "success", jwtToken)
}
