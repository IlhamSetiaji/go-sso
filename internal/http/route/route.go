package route

import (
	"app/go-sso/internal/http/handler"
	"app/go-sso/internal/http/handler/web"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App               *gin.Engine
	UserHandler       handler.UserHandlerInterface
	DashboardHandler  web.DashboardHandlerInterface
	AuthWebHandler    web.AuthHandlerInterface
	WebAuthMiddleware gin.HandlerFunc
	AuthMiddleware    gin.HandlerFunc
}

func (c *RouteConfig) SetupRoutes() {
	c.SetupApiRoutes()
	c.SetupWebRoutes()
	c.SetupOAuthRoutes()
}

func (c *RouteConfig) SetupApiRoutes() {
	apiRoute := c.App.Group("/api")
	{
		userRoute := apiRoute.Group("/user")
		{
			userRoute.POST("/login", c.UserHandler.Login)
			userRoute.Use(c.AuthMiddleware)
			{
				userRoute.GET("/me", c.UserHandler.Me)
			}
		}
		oAuthRoute := apiRoute.Group("/oauth")
		{
			oAuthRoute.GET("/callback", c.UserHandler.CallbackOAuth)
		}
	}
}

func (c *RouteConfig) SetupWebRoutes() {

	c.App.GET("/login", c.AuthWebHandler.LoginView)
	c.App.POST("/login", c.AuthWebHandler.Login)
	c.App.Use(c.WebAuthMiddleware)
	{
		c.App.GET("/", c.DashboardHandler.Index)
		c.App.GET("/logout", c.AuthWebHandler.Logout)
	}
}

func (c *RouteConfig) SetupOAuthRoutes() {
	oAuthRoute := c.App.Group("/oauth")
	{
		oAuthRoute.GET("/login", c.UserHandler.LoginOAuth)
	}
}
