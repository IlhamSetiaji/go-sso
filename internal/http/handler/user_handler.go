package handler

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/http/middleware"
	request "app/go-sso/internal/http/request/user"
	authUsecase "app/go-sso/internal/usecase/auth_token"
	usecase "app/go-sso/internal/usecase/user"
	"app/go-sso/utils"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/exp/rand"
	"golang.org/x/oauth2"
)

type UserHandler struct {
	Config             *viper.Viper
	Log                *logrus.Logger
	Validate           *validator.Validate
	OAuthConfig        *config.Authenticator
	GoogleOAuthConfig  *config.GoogleAuthenticator
	ZitadelOAuthConfig *config.ZitadelAuthenticator
}

type UserHandlerInterface interface {
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	LogoutCookie(ctx *gin.Context)
	CheckAuthToken(ctx *gin.Context)
	CheckStoredCookie(ctx *gin.Context)
	Me(ctx *gin.Context)
	LoginOAuth(ctx *gin.Context)
	CallbackOAuth(ctx *gin.Context)
	GoogleLoginOAuth(ctx *gin.Context)
	GoogleCallbackOAuth(ctx *gin.Context)
	ZitadelLoginOAuth(ctx *gin.Context)
	ZitadelCallbackOAuth(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
}

func UserHandlerFactory(log *logrus.Logger, validator *validator.Validate, oAuthConfig *config.Authenticator, googleOAuthConfig *config.GoogleAuthenticator, zitadelOAuthConfig *config.ZitadelAuthenticator) UserHandlerInterface {
	config := viper.New()
	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	return &UserHandler{
		Config:             config,
		Log:                log,
		Validate:           validator,
		OAuthConfig:        oAuthConfig,
		GoogleOAuthConfig:  googleOAuthConfig,
		ZitadelOAuthConfig: zitadelOAuthConfig,
	}
}

var codeVerifierStore = make(map[string]string)

func generateCodeVerifier() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~"
	var seededRand = rand.New(rand.NewSource(rand.Uint64()))

	b := make([]byte, 64)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash[:])
}

func (h *UserHandler) Login(ctx *gin.Context) {
	payload := new(request.LoginRequest)
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, 400, "error", err.Error())
		h.Log.Errorf("Error when binding request: %v", err)
		return
	}
	err := h.Validate.Struct(payload)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "error", err.Error())
		h.Log.Errorf("Error when validating request: %v", err)
		return
	}
	factory := usecase.LoginUseCaseFactory(h.Log)
	response, err := factory.Execute(usecase.ILoginUseCaseRequest{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when login: %v", err)
		return
	}
	token, err := utils.GenerateToken(&response.User)
	if err != nil {
		h.Log.Errorf("Error when generating token: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	// authFactory := authUsecase.StoreTokenUseCaseFactory(h.Log)
	// authToken, err := authFactory.Execute(authUsecase.IStoreTokenUseCaseRequest{
	// 	UserID:    response.User.ID,
	// 	Token:     token,
	// 	ExpiredAt: time.Now().Add((6) * time.Hour),
	// })

	// if err != nil {
	// 	utils.ErrorResponse(ctx, 500, "error", err.Error())
	// 	h.Log.Errorf("Error when storing token: %v", err)
	// 	return
	// }

	// if authToken == nil {
	// 	utils.ErrorResponse(ctx, 500, "error", "Failed to store auth token")
	// 	h.Log.Errorf("Failed to store auth token")
	// 	return
	// }

	// var data = map[string]interface{}{
	// 	"token":      authToken.AuthToken.Token,
	// 	"token_type": "Bearer",
	// 	"user":       response.User,
	// }

	var data = map[string]interface{}{
		"token":      token,
		"token_type": "Bearer",
		"user":       response.User,
	}

	jwtCookie := utils.NewDefaultCookieOptions("jwt_token")
	jwtCookie.Domain = h.Config.GetString("app.domain")
	utils.SetTokenCookie(ctx, token, jwtCookie)

	utils.SuccessResponse(ctx, 200, "success", data)
}

func (h *UserHandler) CheckStoredCookie(ctx *gin.Context) {
	token, err := ctx.Cookie("jwt_token")
	if err != nil {
		utils.ErrorResponse(ctx, 400, "error", "Cookie not found")
		return
	}

	utils.SuccessResponse(ctx, 200, "success", token)
}

func (h *UserHandler) Logout(ctx *gin.Context) {
	user, err := middleware.GetUser(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when getting user: %v", err)
		return
	}
	if user == nil {
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		h.Log.Errorf("User not found")
		return
	}
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header provided"})
		ctx.Abort()
		return
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		utils.ErrorResponse(ctx, 401, "error", "Invalid Authorization header format")
		h.Log.Errorf("Invalid Authorization header format")
		return
	}
	factory := authUsecase.DeleteTokenUseCaseFactory(h.Log)
	message, err := factory.Execute(authUsecase.IDeleteTokenUseCaseRequest{
		UserID: user["userId"].(string),
		Token:  bearerToken[1],
	})
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when deleting token: %v", err)
		return
	}
	utils.SuccessResponse(ctx, 200, "success", message)
}

func (h *UserHandler) LogoutCookie(ctx *gin.Context) {
	utils.ClearTokenCookie(ctx, "access_token", h.Config.GetString("app.domain"))
	utils.ClearTokenCookie(ctx, "jwt_token", h.Config.GetString("app.domain"))
	utils.SuccessResponse(ctx, 200, "success", "Logged out successfully")
}

func (h *UserHandler) CheckAuthToken(ctx *gin.Context) {
	payload := new(request.CheckAuthTokenRequest)
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, 400, "error", err.Error())
		h.Log.Errorf("Error when binding request: %v", err)
		return
	}
	err := h.Validate.Struct(payload)
	if err != nil {
		utils.ErrorResponse(ctx, 400, "error", err.Error())
		h.Log.Errorf("Error when validating request: %v", err)
		return
	}
	factory := authUsecase.FindTokenUseCaseFactory(h.Log)
	response, err := factory.Execute(authUsecase.IFindTokenUseCaseRequest{
		UserID: uuid.MustParse(payload.UserID),
		Token:  payload.Token,
	})
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when finding token: %v", err)
		return
	}

	if response == nil {
		utils.ErrorResponse(ctx, 404, "error", "Token not found")
		h.Log.Errorf("Token not found")
		return
	}

	utils.SuccessResponse(ctx, 200, "success", response)
}

func (h *UserHandler) Me(ctx *gin.Context) {
	user, err := middleware.GetUser(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when getting user: %v", err)
		return
	}
	if user == nil {
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		h.Log.Errorf("User not found")
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
		h.Log.Errorf("Invalid state")
		return
	}
	token, err := h.OAuthConfig.Exchange(ctx, code)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when exchanging token: %v", err)
		return
	}
	idToken, err := h.OAuthConfig.VerifyIDToken(ctx, token)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when verifying id token: %v", err)
		return
	}
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when getting profile: %v", err)
		return
	}
	factory := usecase.FindByEmailUseCaseFactory(h.Log)
	response, err := factory.Execute(usecase.IFindByEmailUseCaseRequest{
		Email: profile["email"].(string),
	})

	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when finding user by email: %v", err)
		return
	}
	jwtToken, err := utils.GenerateToken(response.User)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when generating token: %v", err)
		return
	}
	redirectURL := fmt.Sprintf("%s?token=%s", appConfig.RedirectURI, jwtToken)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *UserHandler) GoogleLoginOAuth(ctx *gin.Context) {
	state := ctx.Query("state")
	url := h.GoogleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *UserHandler) GoogleCallbackOAuth(ctx *gin.Context) {
	// Get the code and state from the query
	code := ctx.Query("code")
	state := ctx.Query("state")

	if code == "" {
		utils.ErrorResponse(ctx, 400, "error", "code is required")
		return
	}

	appConfig, appExists := config.GoogleAppConfigs[state]
	if !appExists {
		utils.ErrorResponse(ctx, 400, "error", "Invalid state")
		h.Log.Errorf("Invalid state")
		return
	}

	token, err := h.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", "failed to exchange token")
		return
	}

	client := h.GoogleOAuthConfig.Client(context.Background(), token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", "failed to get user info")
		return
	}
	defer userInfoResp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		utils.ErrorResponse(ctx, 500, "error", "failed to decode user info")
		return
	}

	email := userInfo.Email
	if email == "" {
		utils.ErrorResponse(ctx, 400, "error", "email is required")
		return
	}

	factory := usecase.FindByEmailUseCaseFactory(h.Log)
	response, err := factory.Execute(usecase.IFindByEmailUseCaseRequest{
		Email: email,
	})

	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when finding user by email: %v", err)
		return
	}
	jwtToken, err := utils.GenerateToken(response.User)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when generating token: %v", err)
		return
	}

	redirectURL := fmt.Sprintf("%s?token=%s", appConfig.RedirectURI, jwtToken)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *UserHandler) ZitadelLoginOAuth(ctx *gin.Context) {
	state := ctx.Query("state")
	codeVerifier := generateCodeVerifier()
	codeChallenge := generateCodeChallenge(codeVerifier)

	codeVerifierStore[state] = codeVerifier

	authURL := h.ZitadelOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline) +
		"&code_challenge=" + codeChallenge +
		"&code_challenge_method=S256"

	ctx.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *UserHandler) ZitadelCallbackOAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	if code == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "code is required")
		return
	}

	appConfig, appExists := config.ZitadelAppConfigs[state]
	if !appExists {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "Invalid state")
		return
	}

	codeVerifier, ok := codeVerifierStore[state]
	if !ok {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "Invalid or expired state")
		return
	}
	delete(codeVerifierStore, state)

	token, err := h.ZitadelOAuthConfig.Exchange(
		context.Background(),
		code,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier),
		oauth2.SetAuthURLParam("grant_type", "authorization_code"),
	)
	if err != nil {
		h.Log.Errorf("Token exchange failed: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", "Failed to exchange token")
		return
	}

	client := h.ZitadelOAuthConfig.Client(context.Background(), token)
	userInfoResp, err := client.Get("https://signals99-kjlnde.us1.zitadel.cloud/oidc/v1/userinfo")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", "Failed to get user info: "+err.Error())
		return
	}
	defer userInfoResp.Body.Close()

	accessToken := token.AccessToken

	accessTokenCookie := utils.NewDefaultCookieOptions("access_token")
	accessTokenCookie.Domain = h.Config.GetString("app.domain")
	utils.SetTokenCookie(ctx, accessToken, accessTokenCookie)

	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		h.Log.Errorf("Failed to decode user info: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", "Failed to decode user info: "+err.Error())
		return
	}

	if userInfo.Email == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "Email is required")
		return
	}

	factory := usecase.FindByEmailUseCaseFactory(h.Log)
	response, err := factory.Execute(usecase.IFindByEmailUseCaseRequest{
		Email: userInfo.Email,
	})
	if err != nil {
		h.Log.Errorf("Error finding user by email: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	jwtToken, err := utils.GenerateToken(response.User)
	if err != nil {
		h.Log.Errorf("Error generating token: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", "Failed to generate token")
		return
	}

	jwtCookie := utils.NewDefaultCookieOptions("jwt_token")
	jwtCookie.Domain = h.Config.GetString("app.domain")
	utils.SetTokenCookie(ctx, jwtToken, jwtCookie)

	redirectURL := fmt.Sprintf("%s?token=%s", appConfig.RedirectURI, jwtToken)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *UserHandler) FindById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, 400, "error", "ID is required")
		h.Log.Errorf("ID is required")
		return
	}
	factory := usecase.FindByIdUseCaseFactory(h.Log)
	response, err := factory.Execute(&usecase.IFindByIdUseCaseRequest{
		ID: uuid.MustParse(id),
	})
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when finding user by ID: %v", err)
		return
	}
	utils.SuccessResponse(ctx, 200, "success", response.User)
}

func (h *UserHandler) FindAllPaginated(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := ctx.Query("search")
	if err != nil {
		search = ""
	}

	factory := usecase.FindAllPaginatedUseCaseFactory(h.Log)
	response, err := factory.Execute(&usecase.IFindAllPaginatedRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	})

	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when finding all users: %v", err)
		return
	}

	utils.SuccessResponse(ctx, 200, "success", response)
}
