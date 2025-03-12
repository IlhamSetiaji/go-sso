package handler

import (
	"app/go-sso/internal/http/middleware"
	usecase "app/go-sso/internal/usecase/organization"
	locationUsecase "app/go-sso/internal/usecase/organization_location"
	structureUsecase "app/go-sso/internal/usecase/organization_structure"
	orgTypeUsecase "app/go-sso/internal/usecase/organization_type"
	userUsecase "app/go-sso/internal/usecase/user"
	"app/go-sso/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	FindOrganizationStructurePaginated(ctx *gin.Context)
	FindOrganizationStructureById(ctx *gin.Context)
	FindOrganizationStructureByIdWithParents(ctx *gin.Context)
	FindOrganizationLocationsPaginated(ctx *gin.Context)
	FindOrganizationLocationById(ctx *gin.Context)
	FindOrganizationLocationByOrganizationId(ctx *gin.Context)
	FindOrganizationTypesPaginated(ctx *gin.Context)
	FindOrganizationTypeById(ctx *gin.Context)
	UploadLogoOrganization(ctx *gin.Context)
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
	// middleware.PermissionApiMiddleware("read-organization")(ctx)
	// if ctx.IsAborted() {
	// 	ctx.Abort()
	// 	return
	// }
	middleware.PermissionApiMiddleware("read-organization")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.log.Errorf("Permission denied")
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
	middleware.PermissionApiMiddleware("read-organization")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.log.Errorf("Permission denied")
		return
	}

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
	res.Organization.Logo = h.Config.GetString("app.url") + res.Organization.Logo

	utils.SuccessResponse(ctx, http.StatusOK, "success", res.Organization)
}

func (h *OrganizationHandler) FindOrganizationStructureByIdWithParents(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-organization-structure")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.log.Errorf("Permission denied")
		return
	}

	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	factory := structureUsecase.FindByIDWithParentUseCaseFactory(h.log)
	res, err := factory.Execute(&structureUsecase.IFindByIDWithParentUseCaseRequest{
		ID: uuid.MustParse(id),
	})
	if err != nil {
		h.log.Errorf("Error: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res.OrganizationStructure)
}

func (h *OrganizationHandler) FindOrganizationStructurePaginated(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-organization-structure")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.log.Errorf("Permission denied")
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
	if search == "" {
		search = ""
	}

	user, err := middleware.GetUser(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.log.Errorf("Error when getting user: %v", err)
		return
	}
	if user == nil {
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		h.log.Errorf("User not found")
		return
	}

	userFactory := userUsecase.MeUseCaseFactory(h.log)
	resp, err := userFactory.Execute(&userUsecase.IMeUseCaseRequest{
		ID:          uuid.MustParse(user["id"].(string)),
		ChoosedRole: user["choosed_role"].(string),
	})

	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.log.Errorf("Error when finding user by ID: %v", err)
		return
	}

	factory := structureUsecase.FindAllPaginatedUseCaseFactory(h.log)
	res, err := factory.Execute(&structureUsecase.IFindAllPaginatedUseCaseRequest{
		Page:           page,
		PageSize:       pageSize,
		Search:         search,
		OrganizationID: resp.User.Employee.OrganizationID,
	})
	if err != nil {
		h.log.Errorf("Error: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res)
}

func (h *OrganizationHandler) FindOrganizationStructureById(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-organization-structure")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.log.Errorf("Permission denied")
		return
	}

	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	factory := structureUsecase.FindByIdUseCaseFactory(h.log)
	res, err := factory.Execute(&structureUsecase.IFindByIdUseCaseRequest{
		ID: uuid.MustParse(id),
	})
	if err != nil {
		h.log.Errorf("Error: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res.OrganizationStructure)
}

func (h *OrganizationHandler) FindOrganizationLocationsPaginated(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-organization-location")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.log.Errorf("Permission denied")
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
	if search == "" {
		search = ""
	}

	factory := locationUsecase.FindAllPaginatedUseCaseFactory(h.log)
	res, err := factory.Execute(&locationUsecase.IFindAllPaginatedUseCaseRequest{
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

func (h *OrganizationHandler) FindOrganizationLocationById(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-organization-location")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.log.Errorf("Permission denied")
		return
	}

	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	factory := locationUsecase.FindByIdUseCaseFactory(h.log)
	res, err := factory.Execute(&locationUsecase.IFindByIdUseCaseRequest{
		ID: uuid.MustParse(id),
	})
	if err != nil {
		h.log.Errorf("Error: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res.OrganizationLocation)
}

func (h *OrganizationHandler) FindOrganizationLocationByOrganizationId(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-organization-location")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.log.Errorf("Permission denied")
		return
	}

	organizationId := ctx.Param("organization_id")
	if organizationId == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "organizationId is required")
		return
	}

	factory := locationUsecase.FindByOrganizationIdUseCaseFactory(h.log)
	res, err := factory.Execute(&locationUsecase.IFindByOrganizationIdUseCaseRequest{
		OrganizationID: uuid.MustParse(organizationId),
	})

	if err != nil {
		h.log.Errorf("Error: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res.OrganizationLocations)
}

func (h *OrganizationHandler) FindOrganizationTypesPaginated(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-organization-type")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.log.Errorf("Permission denied")
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
	if search == "" {
		search = ""
	}

	factory := orgTypeUsecase.FindAllPaginatedUseCaseFactory(h.log)
	res, err := factory.Execute(&orgTypeUsecase.IFindAllPaginatedRequest{
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

func (h *OrganizationHandler) FindOrganizationTypeById(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-organization-type")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.log.Errorf("Permission denied")
		return
	}

	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	factory := orgTypeUsecase.FindByIdUseCaseFactory(h.log)
	res, err := factory.Execute(&orgTypeUsecase.IFindByIdUseCaseRequest{
		ID: uuid.MustParse(id),
	})
	if err != nil {
		h.log.Errorf("Error: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", res.OrganizationType)
}

func (h *OrganizationHandler) UploadLogoOrganization(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("update-organization")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.log.Errorf("Permission denied")
		return
	}

	var req usecase.IUploadLogoOrganizationUseCaseRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.log.Error("[OrganizationHandler.UploadLogoOrganization] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.log.Error("[OrganizationHandler.UploadLogoOrganization] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	// handle logo upload
	if req.Logo != nil {
		timestamp := time.Now().UnixNano()
		filePath := "storage/logo/" + strconv.FormatInt(timestamp, 10) + "_" + req.Logo.Filename
		if err := ctx.SaveUploadedFile(req.Logo, filePath); err != nil {
			h.log.Error("failed to save logo file: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save logo file", err.Error())
			return
		}

		req.Logo = nil
		req.LogoPath = filePath
	}

	factory := usecase.UploadLogoOrganizationUseCaseFactory(h.log)
	res, err := factory.Execute(&req)
	if err != nil {
		h.log.Errorf("Error: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	res.Organization.Logo = h.Config.GetString("app.url") + res.Organization.Logo

	utils.SuccessResponse(ctx, http.StatusOK, "success", res.Organization)
}
