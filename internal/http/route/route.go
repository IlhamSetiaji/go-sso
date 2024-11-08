package route

import (
	"app/go-sso/internal/http/handler"
	"app/go-sso/internal/http/handler/web"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App              *gin.Engine
	UserHandler      handler.UserHandlerInterface
	DashboardHandler web.DashboardHandlerInterface
	AuthMiddleware   gin.HandlerFunc
}

func (c *RouteConfig) SetupRoutes() {
	c.App.GET("/", c.DashboardHandler.Index)
	c.SetupApiRoutes()
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

	webRoute := c.App.Group("/")
	{
		webOAuthRoute := webRoute.Group("/oauth")
		{
			webOAuthRoute.GET("/login", c.UserHandler.LoginOAuth)
		}
	}
}
