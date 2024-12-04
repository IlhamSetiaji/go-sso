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
	gob.Register(entity.Profile{})
	// gob.Register([]entity.Role{})
	// gob.Register([]entity.Permission{})
}

func main() {
	// setup config
	viperConfig := config.NewViper()
	log := config.NewLogrus(viperConfig)
	validate := config.NewValidator(viperConfig)
	auth, err := config.NewAuth0(viperConfig)
	if err != nil {
		log.Printf("Failed to initialize the auth0 authenticator: %v", err)
	}
	googleAuth, err := config.NewGoogleAuthenticator(viperConfig)
	if err != nil {
		log.Printf("Failed to initialize the google authenticator: %v", err)
	}
	zitadelAuth, err := config.NewZitadelAuthenticator(viperConfig)
	if err != nil {
		log.Printf("Failed to initialize the zitadel authenticator: %v", err)
	}

	// err = rabbitmq.InitializeConnection(viperConfig.GetString("rabbitmq.url"))
	// if err != nil {
	// 	log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	// }
	// defer rabbitmq.CloseConnection()

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

	// app.Use(middleware.FlashMiddleware())

	//handle handler
	userHandler := handler.UserHandlerFactory(log, validate, auth, googleAuth, zitadelAuth)
	organizationHandler := handler.OrganizationHandlerFactory(log, validate)
	jobHandler := handler.JobHandlerFactory(log, validate)
	employeeHandler := handler.EmployeeHandlerFactory(log, validate)

	// handle web handler
	dashboardHandler := web.DashboardHandlerFactory(log, validate)
	authWebHandler := web.AuthHandlerFactory(log, validate)
	userWebHandler := web.UserHandlerFactory(log, validate)
	roleWebHandler := web.RoleHandlerFactory(log, validate)
	permissionWebHandler := web.PermissionHandlerFactory(log, validate)

	// handle middleware
	authMiddleware := middleware.NewAuth(viperConfig)
	authWebMiddleware := middleware.WebAuthMiddleware()

	// setup route config
	routeConfig := route.RouteConfig{
		App:                  app,
		UserHandler:          userHandler,
		DashboardHandler:     dashboardHandler,
		AuthWebHandler:       authWebHandler,
		OrganizationHandler:  organizationHandler,
		JobHandler:           jobHandler,
		UserWebHandler:       userWebHandler,
		RoleWebHandler:       roleWebHandler,
		PermissionWebHandler: permissionWebHandler,
		AuthMiddleware:       authMiddleware,
		WebAuthMiddleware:    authWebMiddleware,
		EmployeeHandler:      employeeHandler,
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
