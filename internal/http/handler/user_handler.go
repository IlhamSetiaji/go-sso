package handler

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/http/middleware"
	request "app/go-sso/internal/http/request/user"
	authUsecase "app/go-sso/internal/usecase/auth_token"
	usecase "app/go-sso/internal/usecase/user"
	"app/go-sso/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type UserHandler struct {
	Log         *logrus.Logger
	Validate    *validator.Validate
	OAuthConfig *config.Authenticator
}

type UserHandlerInterface interface {
	Login(ctx *gin.Context)
	Me(ctx *gin.Context)
	LoginOAuth(ctx *gin.Context)
	CallbackOAuth(ctx *gin.Context)
}

func UserHandlerFactory(log *logrus.Logger, validator *validator.Validate, oAuthConfig *config.Authenticator) UserHandlerInterface {
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
	factory := usecase.LoginUseCaseFactory(h.Log)
	response, err := factory.Execute(usecase.ILoginUseCaseRequest{
		Email:    payload.Email,
		Password: payload.Password,
	})
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

	authFactory := authUsecase.StoreTokenUseCaseFactory(h.Log)
	authToken, err := authFactory.Execute(authUsecase.IStoreTokenUseCaseRequest{
		UserID:    response.User.ID,
		Token:     token,
		ExpiredAt: time.Now().Add(6 * time.Hour),
	})

	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Panicf("Error when storing token: %v", err)
		return
	}

	if authToken == nil {
		utils.ErrorResponse(ctx, 500, "error", "Failed to store auth token")
		h.Log.Panicf("Failed to store auth token")
		return
	}

	var data = map[string]interface{}{
		"token":      authToken.AuthToken.Token,
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
	state := ctx.Query("state")
	url := h.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *UserHandler) CallbackOAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	appConfig, appExists := config.AppConfigs[state]
	if !appExists {
		utils.ErrorResponse(ctx, 400, "error", "Invalid state")
		h.Log.Panicf("Invalid state")
		return
	}
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
	factory := usecase.FindByEmailUseCaseFactory(h.Log)
	response, err := factory.Execute(usecase.IFindByEmailUseCaseRequest{
		Email: profile["email"].(string),
	})

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
	redirectURL := fmt.Sprintf("%s?token=%s", appConfig.RedirectURI, jwtToken)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
	// utils.SuccessResponse(ctx, 200, "success", jwtToken)
}
