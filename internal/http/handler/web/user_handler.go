package web

import (
	"app/go-sso/internal/http/middleware"
	usecase "app/go-sso/internal/usecase/user"
	"app/go-sso/views"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	Log      *logrus.Logger
	Validate *validator.Validate
}

type UserHandlerInterface interface {
	Index(ctx *gin.Context)
}

func UserHandlerFactory(log *logrus.Logger, validator *validator.Validate) UserHandlerInterface {
	return &UserHandler{
		Log:      log,
		Validate: validator,
	}
}

func (h *UserHandler) Index(ctx *gin.Context) {
	middleware.PermissionMiddleware("read-user")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	factory := usecase.FindAllPaginatedFactory(h.Log)

	req := &usecase.FindAllPaginatedRequest{
		Page:     page,
		PageSize: pageSize,
	}
	resp, err := factory.FindAllPaginated(req)
	if err != nil {
		h.Log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	index := views.NewView("base", "views/users/index.html")
	data := map[string]interface{}{
		"Title": "Go SSO | Users",
		"Users": resp.Users,
	}

	index.Render(ctx, data)
}
