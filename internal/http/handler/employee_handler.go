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
	CountEmployeeRetiredEndByDateRange(ctx *gin.Context)
	FindEmployeeRecruitmentManager(ctx *gin.Context)
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
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.Log.Errorf("Permission denied")
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	isOnboarding := ctx.Query("is_onboarding")
	if isOnboarding != "" && isOnboarding != "YES" && isOnboarding != "NO" {
		isOnboarding = ""
	}

	search := ctx.Query("search")

	req := &usecase.IFindAllPaginatedUseCaseRequest{
		Page:         page,
		PageSize:     pageSize,
		Search:       search,
		IsOnboarding: isOnboarding,
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
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.Log.Errorf("Permission denied")
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

func (h *EmployeeHandler) CountEmployeeRetiredEndByDateRange(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-employee")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.Log.Errorf("Permission denied")
		return
	}

	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	if startDate == "" || endDate == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "start_date and end_date is required")
		return
	}

	req := &usecase.ICountEmployeeRetiredEndByDateRangeUseCaseRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

	uc := usecase.CountEmployeeRetiredEndByDateRangeUseCaseFactory(h.Log)
	res, err := uc.Execute(req)
	if err != nil {
		h.Log.Errorf("Error CountEmployeeRetiredEndByDateRangeUseCaseFactory: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}

func (h *EmployeeHandler) FindEmployeeRecruitmentManager(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-employee")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.Log.Errorf("Permission denied")
		return
	}

	uc := usecase.FindEmployeeRecruitmentManagerUsecaseFactory(h.Log)
	res, err := uc.Execute()
	if err != nil {
		h.Log.Errorf("Error FindEmployeeRecruitmentManagerUsecaseFactory: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res.Employee)
}
