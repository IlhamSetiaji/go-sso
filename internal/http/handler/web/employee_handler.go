package web

import (
	"app/go-sso/internal/http/middleware"
	usecase "app/go-sso/internal/usecase/employee"
	jobUsecase "app/go-sso/internal/usecase/job"
	orgUsecase "app/go-sso/internal/usecase/organization"
	orgLocUsecase "app/go-sso/internal/usecase/organization_location"
	userUsecase "app/go-sso/internal/usecase/user"
	"app/go-sso/views"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type EmployeeHandler struct {
	Log      *logrus.Logger
	Validate *validator.Validate
}

type EmployeeHandlerInterface interface {
	Index(ctx *gin.Context)
	Store(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	EmployeeJobs(ctx *gin.Context)
}

func EmployeeHandlerFactory(log *logrus.Logger, validator *validator.Validate) EmployeeHandlerInterface {
	return &EmployeeHandler{
		Log:      log,
		Validate: validator,
	}
}

func (h *EmployeeHandler) Index(ctx *gin.Context) {
	middleware.PermissionMiddleware("read-employee")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	factory := usecase.FindAllEmployeesUsecaseFactory(h.Log)

	resp, err := factory.Execute()
	if err != nil {
		h.Log.Error(err)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	userFactory := userUsecase.GetAllUsersDoesntHaveEmployeeUsecaseFactory(h.Log)

	userResp, err := userFactory.Execute()
	if err != nil {
		h.Log.Error(err)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	orgFactory := orgUsecase.FindAllOrganizationsUsecaseFactory(h.Log)

	orgResp, err := orgFactory.Execute()
	if err != nil {
		h.Log.Error(err)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	index := views.NewView("base", "views/employees/index.html")
	data := map[string]interface{}{
		"Title":         "Go SSO | Employee",
		"Employees":     resp.Employees,
		"Users":         userResp.Users,
		"Organizations": orgResp.Organizations,
	}

	index.Render(ctx, data)
}

func (h *EmployeeHandler) Store(ctx *gin.Context) {
	middleware.PermissionMiddleware("create-employee")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	session := sessions.Default(ctx)

	var request usecase.IStoreEmployeeUsecaseRequest
	if err := ctx.ShouldBind(&request); err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err)
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	if err := h.Validate.Struct(request); err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err)
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	factory := usecase.StoreEmployeeUsecaseFactory(h.Log)

	_, err := factory.Execute(&request)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err)
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	session.Set("success", "Employee created")
	session.Save()
	ctx.Redirect(302, ctx.Request.Referer())
}

func (h *EmployeeHandler) Update(ctx *gin.Context) {
	middleware.PermissionMiddleware("update-employee")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	session := sessions.Default(ctx)

	var request usecase.IUpdateEmployeeUsecaseRequest
	if err := ctx.ShouldBind(&request); err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err)
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	if err := h.Validate.Struct(request); err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err)
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	factory := usecase.UpdateEmployeeUsecaseFactory(h.Log)

	_, err := factory.Execute(&request)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err)
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	session.Set("success", "Employee updated")
	session.Save()
	ctx.Redirect(302, ctx.Request.Referer())
}

func (h *EmployeeHandler) Delete(ctx *gin.Context) {
	middleware.PermissionMiddleware("delete-employee")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	session := sessions.Default(ctx)

	var request usecase.IDeleteEmployeeUsecaseRequest
	if err := ctx.ShouldBind(&request); err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err)
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	if err := h.Validate.Struct(request); err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err)
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	factory := usecase.DeleteEmployeeUsecaseFactory(h.Log)

	_, err := factory.Execute(&request)
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err)
		ctx.Redirect(302, ctx.Request.Referer())
		return
	}

	session.Set("success", "Employee deleted")
	session.Save()
	ctx.Redirect(302, ctx.Request.Referer())
}

func (h *EmployeeHandler) EmployeeJobs(ctx *gin.Context) {
	middleware.PermissionMiddleware("read-employee")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	id := ctx.Param("id")

	session := sessions.Default(ctx)

	factory := usecase.FindByIdUseCaseFactory(h.Log)

	resp, err := factory.Execute(&usecase.IFindByIdUseCaseRequest{
		ID: uuid.MustParse(id),
	})
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err)
		// ctx.JSON(500, gin.H{"error": err.Error()})
		ctx.Redirect(302, "/employees")
		return
	}

	orgFactory := orgUsecase.FindAllOrganizationsUsecaseFactory(h.Log)

	orgResp, err := orgFactory.Execute()
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err)
		// ctx.JSON(500, gin.H{"error": err.Error()})
		ctx.Redirect(302, "/employees")
		return
	}

	jobFactory := jobUsecase.GetJobsByOrganizationIDUseCaseFactory(h.Log)

	jobResp, err := jobFactory.Execute(&jobUsecase.IGetJobsByOrganizationIDUseCaseRequest{
		OrganizationID: resp.Employee.OrganizationID,
	})
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err)
		// ctx.JSON(500, gin.H{"error": err.Error()})
		ctx.Redirect(302, "/employees")
		return
	}

	orgLocFactory := orgLocUsecase.FindByOrganizationIdUseCaseFactory(h.Log)

	orgLocResp, err := orgLocFactory.Execute(&orgLocUsecase.IFindByOrganizationIdUseCaseRequest{
		OrganizationID: resp.Employee.OrganizationID,
	})
	if err != nil {
		session.Set("error", err.Error())
		session.Save()
		h.Log.Error(err)
		// ctx.JSON(500, gin.H{"error": err.Error()})
		ctx.Redirect(302, "/employees")
		return
	}

	var tipe string = "create"

	if resp.Employee.EmployeeJob != nil {
		tipe = "update"
	}

	index := views.NewView("base", "views/employees/employee_job.html")
	data := map[string]interface{}{
		"Title":         "Go SSO | Employee Jobs",
		"Employee":      resp.Employee,
		"Organizations": orgResp.Organizations,
		"Jobs":          jobResp.Jobs,
		"OrgLocs":       orgLocResp.OrganizationLocations,
		"Tipe":          tipe,
	}

	index.Render(ctx, data)
}
