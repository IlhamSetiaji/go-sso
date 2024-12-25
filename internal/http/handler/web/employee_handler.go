package web

import (
	"app/go-sso/internal/http/middleware"
	usecase "app/go-sso/internal/usecase/employee"
	userUsecase "app/go-sso/internal/usecase/user"
	"app/go-sso/views"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type EmployeeHandler struct {
	Log      *logrus.Logger
	Validate *validator.Validate
}

type EmployeeHandlerInterface interface {
	Index(ctx *gin.Context)
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

	index := views.NewView("base", "views/employees/index.html")
	data := map[string]interface{}{
		"Title":     "Go SSO | Employee",
		"Employees": resp.Employees,
		"Users":     userResp.Users,
	}

	index.Render(ctx, data)
}
