package main

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/http/handler"
	"app/go-sso/internal/http/handler/web"
	"app/go-sso/internal/http/middleware"
	"app/go-sso/internal/http/route"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// setup config
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	validate := config.NewValidator(viperConfig)
	auth, err := config.NewAuth0(viperConfig)
	if err != nil {
		log.Panicf("Failed to initialize the authenticator: %v", err)
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

	//handle handler
	userHandler := handler.UserHandlerFactory(log, validate, auth)
	dashboardHandler := web.DashboardHandlerFactory(log, validate)

	// handle middleware
	authMiddleware := middleware.NewAuth(viperConfig)

	// setup route config
	routeConfig := route.RouteConfig{
		App:              app,
		UserHandler:      userHandler,
		DashboardHandler: dashboardHandler,
		AuthMiddleware:   authMiddleware,
	}
	routeConfig.SetupRoutes()

	// run server
	webPort := strconv.Itoa(viperConfig.GetInt("web.port"))
	err = app.Run(":" + webPort)
	if err != nil {
		log.Panicf("Failed to start server: %v", err)
	}
}
