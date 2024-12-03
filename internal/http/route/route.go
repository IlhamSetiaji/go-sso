package route

import (
	"app/go-sso/internal/http/handler"
	"app/go-sso/internal/http/handler/web"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App                  *gin.Engine
	UserHandler          handler.UserHandlerInterface
	UserWebHandler       web.UserHandlerInterface
	RoleWebHandler       web.RoleHandlerInterface
	OrganizationHandler  handler.IOrganizationHandler
	PermissionWebHandler web.PermissionHandlerInterface
	DashboardHandler     web.DashboardHandlerInterface
	AuthWebHandler       web.AuthHandlerInterface
	WebAuthMiddleware    gin.HandlerFunc
	AuthMiddleware       gin.HandlerFunc
	JobHandler           handler.IJobHandler
	EmployeeHandler      handler.IEmployeeHandler
}

func (c *RouteConfig) SetupRoutes() {
	c.SetupApiRoutes()
	c.SetupOAuthRoutes()
	c.SetupWebRoutes()
}

func (c *RouteConfig) SetupApiRoutes() {
	apiRoute := c.App.Group("/api")
	{
		apiRoute.POST("/login", c.UserHandler.Login)

		oAuthRoute := apiRoute.Group("/oauth")
		{
			oAuthRoute.GET("/callback", c.UserHandler.CallbackOAuth)
			oAuthRoute.GET("/google/callback", c.UserHandler.GoogleCallbackOAuth)
			oAuthRoute.GET("/zitadel/callback", c.UserHandler.ZitadelCallbackOAuth)
		}

		apiRoute.Use(c.AuthMiddleware)
		{
			userRoute := apiRoute.Group("/users")
			{
				userRoute.GET("/me", c.UserHandler.Me)
				userRoute.GET("/logout/token", c.UserHandler.Logout)
				userRoute.GET("/logout", c.UserHandler.LogoutCookie)
				userRoute.POST("/check-token", c.UserHandler.CheckAuthToken)
				userRoute.GET("/", c.UserHandler.FindAllPaginated)
				userRoute.GET("/:id", c.UserHandler.FindById)
				userRoute.GET("/check-cookie", c.UserHandler.CheckStoredCookie)
			}
			organizationRoute := apiRoute.Group("/organizations")
			{
				organizationRoute.GET("/", c.OrganizationHandler.FindAllPaginated)
				organizationRoute.GET("/:id", c.OrganizationHandler.FindById)
			}
			organizationStructureRoute := apiRoute.Group("/organization-structures")
			{
				organizationStructureRoute.GET("/", c.OrganizationHandler.FindOrganizationStructurePaginated)
				organizationStructureRoute.GET("/:id", c.OrganizationHandler.FindOrganizationStructureById)
			}
			organizationLocationRoute := apiRoute.Group("/organization-locations")
			{
				organizationLocationRoute.GET("/", c.OrganizationHandler.FindOrganizationLocationsPaginated)
				organizationLocationRoute.GET("/:id", c.OrganizationHandler.FindOrganizationLocationById)
			}
			jobRoute := apiRoute.Group("/jobs")
			{
				jobRoute.GET("/", c.JobHandler.FindAllPaginated)
				jobRoute.GET("/:id", c.JobHandler.FindById)
			}
			jobLevelRoute := apiRoute.Group("/job-levels")
			{
				jobLevelRoute.GET("/", c.JobHandler.FindAllJobLevelsPaginated)
				jobLevelRoute.GET("/:id", c.JobHandler.FindJobLevelById)
			}
			employeeRoute := apiRoute.Group("/employees")
			{
				employeeRoute.GET("/", c.EmployeeHandler.FindAllPaginated)
				employeeRoute.GET("/:id", c.EmployeeHandler.FindById)
			}
		}
	}
}

func (c *RouteConfig) SetupWebRoutes() {

	c.App.GET("/login", c.AuthWebHandler.LoginView)
	c.App.POST("/login", c.AuthWebHandler.Login)
	c.App.Use(c.WebAuthMiddleware)
	{
		c.App.GET("/", c.DashboardHandler.Index)
		c.App.GET("/test", c.AuthWebHandler.CheckCookieTest)
		c.App.GET("/portal", c.DashboardHandler.Portal)
		c.App.GET("/logout", c.AuthWebHandler.Logout)
		userRoutes := c.App.Group("/users")
		{
			userRoutes.GET("/", c.UserWebHandler.Index)
			userRoutes.POST("/", c.UserWebHandler.StoreUser)
			userRoutes.POST("/update", c.UserWebHandler.UpdateUser)
			userRoutes.POST("/delete", c.UserWebHandler.DeleteUser)
		}
		roleRoutes := c.App.Group("/roles")
		{
			roleRoutes.GET("/", c.RoleWebHandler.Index)
			roleRoutes.POST("/", c.RoleWebHandler.StoreRole)
			roleRoutes.POST("/update", c.RoleWebHandler.UpdateRole)
			roleRoutes.POST("/delete", c.RoleWebHandler.DeleteRole)
		}
		permissionRoutes := c.App.Group("/permissions")
		{
			permissionRoutes.GET("/", c.PermissionWebHandler.Index)
			permissionRoutes.POST("/", c.PermissionWebHandler.StorePermission)
			permissionRoutes.POST("/update", c.PermissionWebHandler.UpdatePermission)
			permissionRoutes.POST("/delete", c.PermissionWebHandler.DeletePermission)
		}
	}
}

func (c *RouteConfig) SetupOAuthRoutes() {
	oAuthRoute := c.App.Group("/oauth")
	{
		oAuthRoute.GET("/login", c.UserHandler.LoginOAuth)
		oAuthRoute.GET("/google/login", c.UserHandler.GoogleLoginOAuth)
		oAuthRoute.GET("/zitadel/login", c.UserHandler.ZitadelLoginOAuth)
	}
}
