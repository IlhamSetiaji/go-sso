package handler

import (
	"app/go-sso/internal/http/middleware"
	usecase "app/go-sso/internal/usecase/job"
	"app/go-sso/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type JobHandler struct {
	Log      *logrus.Logger
	Validate *validator.Validate
}

type IJobHandler interface {
	FindAllPaginated(ctx *gin.Context)
	FindById(ctx *gin.Context)
}

func NewJobHandler(log *logrus.Logger, validate *validator.Validate) IJobHandler {
	return &JobHandler{
		Log:      log,
		Validate: validate,
	}
}

func JobHandlerFactory(log *logrus.Logger, validate *validator.Validate) IJobHandler {
	return NewJobHandler(log, validate)
}

func (h *JobHandler) FindAllPaginated(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-job")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := ctx.Query("search")
	if search == "" {
		search = ""
	}

	factory := usecase.FindAllPaginatedUseCaseFactory(h.Log)
	response, err := factory.Execute(&usecase.IFindAllPaginatedUseCaseRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	})

	if err != nil {
		h.Log.Errorf("Error FindAllPaginated: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", response)
}

func (h *JobHandler) FindById(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-job")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	factory := usecase.FindByIdUseCaseFactory(h.Log)
	response, err := factory.Execute(&usecase.IFindByIdUseCaseRequest{
		ID: uuid.MustParse(id),
	})

	if err != nil {
		h.Log.Errorf("Error FindById: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", response.Job)
}