package handler

import (
	"app/go-sso/internal/http/middleware"
	usecase "app/go-sso/internal/usecase/employee"
	"app/go-sso/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IEmployeeHandler interface {
	FindAllPaginated(ctx *gin.Context)
	FindById(ctx *gin.Context)
}

type EmployeeHandler struct {
	Config   *viper.Viper
	Log      *logrus.Logger
	Validate *validator.Validate
}

func NewEmployeeHandler(log *logrus.Logger, validate *validator.Validate) IEmployeeHandler {
	config := viper.New()
	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	return &EmployeeHandler{
		Config:   config,
		Log:      log,
		Validate: validate,
	}
}

func EmployeeHandlerFactory(log *logrus.Logger, validate *validator.Validate) IEmployeeHandler {
	return NewEmployeeHandler(log, validate)
}

func (h *EmployeeHandler) FindAllPaginated(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-employee")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := ctx.Query("search")

	req := &usecase.IFindAllPaginatedUseCaseRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	uc := usecase.FindAllPaginatedUseCaseFactory(h.Log)
	res, err := uc.Execute(req)
	if err != nil {
		h.Log.Errorf("Error FindAllPaginatedUseCaseFactory: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}

func (h *EmployeeHandler) FindById(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-employee")(ctx)
	if ctx.IsAborted() {
		ctx.Abort()
		return
	}

	id := ctx.Param("id")

	req := &usecase.IFindByIdUseCaseRequest{
		ID: uuid.MustParse(id),
	}

	uc := usecase.FindByIdUseCaseFactory(h.Log)
	res, err := uc.Execute(req)
	if err != nil {
		h.Log.Errorf("Error FindByIdUseCaseFactory: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res.Employee)
}