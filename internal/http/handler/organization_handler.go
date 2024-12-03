package handler

import (
	usecase "app/go-sso/internal/usecase/organization"
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

type OrganizationHandler struct {
	Config   *viper.Viper
	log      *logrus.Logger
	Validate *validator.Validate
}

type IOrganizationHandler interface {
	FindAllPaginated(ctx *gin.Context)
	FindById(ctx *gin.Context)
}

func NewOrganizationHandler(log *logrus.Logger, validate *validator.Validate) IOrganizationHandler {
	config := viper.New()
	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	return &OrganizationHandler{
		Config:   config,
		log:      log,
		Validate: validate,
	}
}

func OrganizationHandlerFactory(log *logrus.Logger, validate *validator.Validate) IOrganizationHandler {
	return NewOrganizationHandler(log, validate)
}

func (h *OrganizationHandler) FindAllPaginated(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := ctx.Query("search")
	if search == "" {
		search = ""
	}

	factory := usecase.FindAllPaginatedUseCaseFactory(h.log)
	res, err := factory.Execute(&usecase.IFindAllPaginatedRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	})
	if err != nil {
		h.log.Errorf("Error: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}

func (h *OrganizationHandler) FindById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	factory := usecase.FindByIdUseCaseFactory(h.log)
	res, err := factory.Execute(&usecase.IFindByIdUseCaseRequest{
		ID: uuid.MustParse(id),
	})
	if err != nil {
		h.log.Errorf("Error: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res.Organization)
}
