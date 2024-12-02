package web

import (
	"app/go-sso/internal/http/middleware"
	appUsecase "app/go-sso/internal/usecase/application"
	usecase "app/go-sso/internal/usecase/role"
	"app/go-sso/views"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type RoleHandler struct {
	Log      *logrus.Logger
	Validate *validator.Validate
}

type RoleHandlerInterface interface {
	Index(ctx *gin.Context)
}

func RoleHandlerFactory(log *logrus.Logger, validator *validator.Validate) RoleHandlerInterface {
	return &RoleHandler{
		Log:      log,
		Validate: validator,
	}
}

func (h *RoleHandler) Index(ctx *gin.Context) {
	middleware.PermissionMiddleware("read-role")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	factory := usecase.GetAllRolesUseCaseFactory(h.Log)

	resp, err := factory.Execute()
	if err != nil {
		h.Log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	appFactory := appUsecase.GetAllApplicationsUseCaseFactory(h.Log)

	appResp, err := appFactory.Execute()
	if err != nil {
		h.Log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	index := views.NewView("base", "views/roles/index.html")
	data := map[string]interface{}{
		"Title":        "Julong Portal | Roles",
		"Roles":        resp.Roles,
		"Applications": appResp.Applications,
	}

	index.Render(ctx, data)
}
