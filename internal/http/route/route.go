package route

import (
	"app/go-sso/internal/http/handler"
	"app/go-sso/views"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App            *gin.Engine
	UserHandler    handler.UserHandlerInterface
	AuthMiddleware gin.HandlerFunc
}

func (c *RouteConfig) SetupRoutes() {
	c.App.GET("/", func(ctx *gin.Context) {
		index := views.NewView("base", "views/index.html")
		data := map[string]interface{}{
			"Title": "Go SSO",
		}
		index.Render(ctx, data)
	})
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
