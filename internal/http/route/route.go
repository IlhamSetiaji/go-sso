package route

import (
	"app/go-sso/internal/http/handler"
	"app/go-sso/internal/http/handler/web"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App                     *gin.Engine
	Viper                   *viper.Viper
	UserHandler             handler.UserHandlerInterface
	UserWebHandler          web.UserHandlerInterface
	RoleWebHandler          web.RoleHandlerInterface
	OrganizationHandler     handler.IOrganizationHandler
	PermissionWebHandler    web.PermissionHandlerInterface
	DashboardHandler        web.DashboardHandlerInterface
	AuthWebHandler          web.AuthHandlerInterface
	WebAuthMiddleware       gin.HandlerFunc
	AuthMiddleware          gin.HandlerFunc
	EmailVerifiedMiddleware gin.HandlerFunc
	JobHandler              handler.IJobHandler
	EmployeeHandler         handler.IEmployeeHandler
	EmployeeWebHandler      web.EmployeeHandlerInterface
	GradeHandler            handler.IGradeHandler
}

func (c *RouteConfig) SetupRoutes() {
	// Setup API, OAuth, and Web routes
	c.SetupApiRoutes()
	c.SetupOAuthRoutes()
	c.SetupWebRoutes()
}

func (c *RouteConfig) SetupApiRoutes() {
	apiRoute := c.App.Group("/api")
	{
		apiRoute.GET("/check-jwt-token", c.UserHandler.CheckStoredCookie)
		apiRoute.POST("/login", c.UserHandler.Login)

		oAuthRoute := apiRoute.Group("/oauth")
		{
			oAuthRoute.GET("/callback", c.UserHandler.CallbackOAuth)
			oAuthRoute.GET("/google/callback", c.UserHandler.GoogleCallbackOAuth)
			oAuthRoute.GET("/zitadel/callback", c.UserHandler.ZitadelCallbackOAuth)
		}

		apiRoute.Use(c.AuthMiddleware)
		{
			// User routes
			apiRoute.GET("/users", c.UserHandler.FindAllPaginated)
			apiRoute.GET("/users/me", c.UserHandler.Me)
			apiRoute.GET("/users/logout/token", c.UserHandler.Logout)
			apiRoute.GET("/users/logout", c.UserHandler.LogoutCookie)
			apiRoute.GET("/users/:id", c.UserHandler.FindById)
			apiRoute.POST("/check-token", c.UserHandler.CheckAuthToken)
			apiRoute.GET("/check-cookie", c.UserHandler.CheckStoredCookie)

			// Organization routes
			apiRoute.GET("/organizations", c.OrganizationHandler.FindAllPaginated)
			apiRoute.GET("/organizations/:id", c.OrganizationHandler.FindById)
			apiRoute.PUT("/organizations/:id/upload-logo", c.OrganizationHandler.UploadLogoOrganization)

			// Organization structure routes
			apiRoute.GET("/organization-structures", c.OrganizationHandler.FindOrganizationStructurePaginated)
			apiRoute.GET("/organization-structures/:id", c.OrganizationHandler.FindOrganizationStructureById)
			apiRoute.GET("/organization-structures/parents/:id", c.OrganizationHandler.FindOrganizationStructureByIdWithParents)

			// Organization location routes
			apiRoute.GET("/organization-locations", c.OrganizationHandler.FindOrganizationLocationsPaginated)
			apiRoute.GET("/organization-locations/organization/:organization_id", c.OrganizationHandler.FindOrganizationLocationByOrganizationId)
			apiRoute.GET("/organization-locations/:id", c.OrganizationHandler.FindOrganizationLocationById)

			// Organization type routes
			apiRoute.GET("/organization-types", c.OrganizationHandler.FindOrganizationTypesPaginated)
			apiRoute.GET("/organization-types/:id", c.OrganizationHandler.FindOrganizationTypeById)

			// Job routes
			apiRoute.GET("/jobs", c.JobHandler.FindAllPaginated)
			apiRoute.GET("/jobs/:id", c.JobHandler.FindById)
			apiRoute.GET("/jobs/job-level/:job_level_id", c.JobHandler.GetJobsByJobLevelId)
			apiRoute.GET("/jobs/organization/:organization_id", c.JobHandler.GetJobsByOrganizationId)

			// Job level routes
			apiRoute.GET("/job-levels", c.JobHandler.FindAllJobLevelsPaginated)
			apiRoute.GET("/job-levels/:id", c.JobHandler.FindJobLevelById)
			apiRoute.GET("/job-levels/organization/:organization_id", c.JobHandler.FindJobLevelsByOrganizationId)

			// Employee routes
			apiRoute.GET("/employees", c.EmployeeHandler.FindAllPaginated)
			apiRoute.GET("/employees/turnover", c.EmployeeHandler.CountEmployeeRetiredEndByDateRange)
			apiRoute.GET("/employees/recruitment-manager", c.EmployeeHandler.FindEmployeeRecruitmentManager)
			apiRoute.GET("/employees/:id", c.EmployeeHandler.FindById)

			// Grade routes
			apiRoute.GET("/grades/job-level/:job_level_id", c.GradeHandler.FindAllByJobLevelID)
		}
	}
}

func (c *RouteConfig) SetupWebRoutes() {
	c.App.GET("/login", c.AuthWebHandler.LoginView)
	c.App.GET("/choose-roles", c.AuthWebHandler.ChooseRoles)
	c.App.POST("/continue-login", c.AuthWebHandler.ContinueLogin)
	c.App.POST("/login", c.AuthWebHandler.Login)
	c.App.GET("/register", c.AuthWebHandler.RegisterView)
	c.App.POST("/register", c.AuthWebHandler.Register)
	c.App.Use(c.WebAuthMiddleware)
	{
		c.App.GET("/", c.DashboardHandler.Index)
		c.App.GET("/test", c.AuthWebHandler.CheckCookieTest)
		c.App.GET("/logout", c.AuthWebHandler.Logout)
		c.App.GET("/otp", c.AuthWebHandler.OtpView)
		c.App.POST("/verify-email", c.AuthWebHandler.VerifyEmail)
		c.App.GET("/resend-verify-email/:email", c.AuthWebHandler.ResendVerifyEmail)
		c.App.Use(c.EmailVerifiedMiddleware)
		{
			c.App.GET("/portal", c.DashboardHandler.Portal)
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
				roleRoutes.POST("/assign-permissions", c.RoleWebHandler.AssignRoleToPermissions)
				roleRoutes.POST("/resign-permissions", c.RoleWebHandler.ResignRoleFromPermission)
				roleRoutes.POST("/update", c.RoleWebHandler.UpdateRole)
				roleRoutes.POST("/delete", c.RoleWebHandler.DeleteRole)
			}
			permissionRoutes := c.App.Group("/permissions")
			{
				permissionRoutes.GET("/", c.PermissionWebHandler.Index)
				permissionRoutes.GET("/role/:role_id", c.PermissionWebHandler.GetPermissionsByRoleID)
				permissionRoutes.POST("/", c.PermissionWebHandler.StorePermission)
				permissionRoutes.POST("/update", c.PermissionWebHandler.UpdatePermission)
				permissionRoutes.POST("/delete", c.PermissionWebHandler.DeletePermission)
			}
			employeeRoutes := c.App.Group("/employees")
			{
				employeeRoutes.GET("/", c.EmployeeWebHandler.Index)
				employeeRoutes.POST("/", c.EmployeeWebHandler.Store)
				employeeRoutes.POST("/update", c.EmployeeWebHandler.Update)
				employeeRoutes.POST("/delete", c.EmployeeWebHandler.Delete)
				employeeRoutes.GET("/:id/job", c.EmployeeWebHandler.EmployeeJobs)
				employeeRoutes.POST("/store-job", c.EmployeeWebHandler.StoreEmployeeJob)
				employeeRoutes.POST("/update-job", c.EmployeeWebHandler.UpdateEmployeeJob)
			}
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
