package web

import (
	"app/go-sso/views"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DashboardHandler struct {
	Log      *log.Logger
	validate *validator.Validate
}

type DashboardHandlerInterface interface {
	Index(ctx *gin.Context)
}

func DashboardHandlerFactory(log *log.Logger, validate *validator.Validate) DashboardHandlerInterface {
	return &DashboardHandler{
		Log:      log,
		validate: validate,
	}
}

func (h *DashboardHandler) Index(ctx *gin.Context) {
	index := views.NewView("base", "views/index.html")
	data := map[string]interface{}{
		"Title": "Go SSO | Dashboard",
	}
	index.Render(ctx, data)
}
