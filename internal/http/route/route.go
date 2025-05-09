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
	webRoute := c.App.Group("/")
	webRoute.GET("/login", c.AuthWebHandler.LoginView)
	webRoute.GET("/choose-roles", c.AuthWebHandler.ChooseRoles)
	webRoute.POST("/continue-login", c.AuthWebHandler.ContinueLogin)
	webRoute.POST("/login", c.AuthWebHandler.Login)
	webRoute.GET("/register", c.AuthWebHandler.RegisterView)
	webRoute.POST("/register", c.AuthWebHandler.Register)
	webRoute.Use(c.WebAuthMiddleware)
	{
		webRoute.GET("/", c.DashboardHandler.Index)
		webRoute.GET("/test", c.AuthWebHandler.CheckCookieTest)
		webRoute.GET("/logout", c.AuthWebHandler.Logout)
		webRoute.GET("/otp", c.AuthWebHandler.OtpView)
		webRoute.POST("/verify-email", c.AuthWebHandler.VerifyEmail)
		webRoute.GET("/resend-verify-email/:email", c.AuthWebHandler.ResendVerifyEmail)
		webRoute.Use(c.EmailVerifiedMiddleware)
		{
			webRoute.GET("/portal", c.DashboardHandler.Portal)
			userRoutes := webRoute.Group("/users")
			{
				userRoutes.GET("/", c.UserWebHandler.Index)
				userRoutes.POST("/", c.UserWebHandler.StoreUser)
				userRoutes.POST("/update", c.UserWebHandler.UpdateUser)
				userRoutes.POST("/delete", c.UserWebHandler.DeleteUser)
			}
			roleRoutes := webRoute.Group("/roles")
			{
				roleRoutes.GET("/", c.RoleWebHandler.Index)
				roleRoutes.POST("/", c.RoleWebHandler.StoreRole)
				roleRoutes.POST("/assign-permissions", c.RoleWebHandler.AssignRoleToPermissions)
				roleRoutes.POST("/resign-permissions", c.RoleWebHandler.ResignRoleFromPermission)
				roleRoutes.POST("/update", c.RoleWebHandler.UpdateRole)
				roleRoutes.POST("/delete", c.RoleWebHandler.DeleteRole)
			}
			permissionRoutes := webRoute.Group("/permissions")
			{
				permissionRoutes.GET("/", c.PermissionWebHandler.Index)
				permissionRoutes.GET("/role/:role_id", c.PermissionWebHandler.GetPermissionsByRoleID)
				permissionRoutes.POST("/", c.PermissionWebHandler.StorePermission)
				permissionRoutes.POST("/update", c.PermissionWebHandler.UpdatePermission)
				permissionRoutes.POST("/delete", c.PermissionWebHandler.DeletePermission)
			}
			employeeRoutes := webRoute.Group("/employees")
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
