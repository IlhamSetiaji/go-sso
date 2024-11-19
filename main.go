package main

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/handler"
	"app/go-sso/internal/http/handler/web"
	"app/go-sso/internal/http/middleware"
	"app/go-sso/internal/http/route"
	"encoding/gob"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func init() {
	gob.Register(entity.User{})
	// gob.Register([]entity.Role{})
	// gob.Register([]entity.Permission{})
}

func main() {
	// setup config
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	validate := config.NewValidator(viperConfig)
	auth, err := config.NewAuth0(viperConfig)
	if err != nil {
		log.Printf("Failed to initialize the authenticator: %v", err)
	}

	// setup gin engine
	app := gin.Default()
	app.Static("/assets", "./public")
	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("App-Name", viperConfig.GetString("app.name"))
	})

	// setup session and cookie
	store := cookie.NewStore([]byte(viperConfig.GetString("web.cookie.secret")))
	app.Use(sessions.Sessions(viperConfig.GetString("web.session.name"), store))

	// setup custom csrf middleware
	app.Use(func(c *gin.Context) {
		if !shouldExcludeFromCSRF(c.Request.URL.Path) {
			csrf.Middleware(csrf.Options{
				Secret: viperConfig.GetString("web.csrf_secret"),
				ErrorFunc: func(c *gin.Context) {
					c.String(http.StatusForbidden, "CSRF token mismatch")
					c.Abort()
				},
			})(c)
		}
		c.Next()
	})

	//handle handler
	userHandler := handler.UserHandlerFactory(log, validate, auth)
	dashboardHandler := web.DashboardHandlerFactory(log, validate)
	authWebHandler := web.AuthHandlerFactory(log, validate)
	userWebHandler := web.UserHandlerFactory(log, validate)

	// handle middleware
	authMiddleware := middleware.NewAuth(viperConfig)
	authWebMiddleware := middleware.WebAuthMiddleware()

	// setup route config
	routeConfig := route.RouteConfig{
		App:               app,
		UserHandler:       userHandler,
		DashboardHandler:  dashboardHandler,
		AuthWebHandler:    authWebHandler,
		UserWebHandler:    userWebHandler,
		AuthMiddleware:    authMiddleware,
		WebAuthMiddleware: authWebMiddleware,
	}
	routeConfig.SetupRoutes()

	// run server
	webPort := strconv.Itoa(viperConfig.GetInt("web.port"))
	err = app.Run(":" + webPort)
	if err != nil {
		log.Panicf("Failed to start server: %v", err)
	}
}

func shouldExcludeFromCSRF(path string) bool {
	return len(path) >= 4 && path[:4] == "/api"
}
