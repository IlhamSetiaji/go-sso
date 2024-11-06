package main

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/http/handler"
	"app/go-sso/internal/http/middleware"
	"app/go-sso/internal/http/route"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	app := config.NewGin(viperConfig)
	validate := config.NewValidator(viperConfig)
	auth, err := config.NewAuth0(viperConfig)
	if err != nil {
		log.Panicf("Failed to initialize the authenticator: %v", err)
	}

	userHandler := handler.UserHandlerFactory(log, validate, auth)

	// handle middleware

	authMiddleware := middleware.NewAuth(viperConfig)

	routeConfig := route.RouteConfig{
		App:            app,
		UserHandler:    userHandler,
		AuthMiddleware: authMiddleware,
	}
	routeConfig.SetupRoutes()

	// webPort := strconv.Itoa(viperConfig.GetInt("web.port"))
	err = app.Run()
	if err != nil {
		// log.Fatalf("Failed to start server: %v", err)
		log.Panicf("Failed to start server: %v", err)
	}
}
