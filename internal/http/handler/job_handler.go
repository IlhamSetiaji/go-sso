package handler

import (
	"app/go-sso/internal/http/middleware"
	usecase "app/go-sso/internal/usecase/job"
	jobLevelUsecase "app/go-sso/internal/usecase/job_level"
	userUsecase "app/go-sso/internal/usecase/user"
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
	FindAllJobLevelsPaginated(ctx *gin.Context)
	FindJobLevelById(ctx *gin.Context)
	GetJobsByJobLevelId(ctx *gin.Context)
	FindJobLevelsByOrganizationId(ctx *gin.Context)
	GetJobsByOrganizationId(ctx *gin.Context)
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
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.Log.Errorf("Permission denied")
		return
	}

	user, err := middleware.GetUser(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when getting user: %v", err)
		return
	}
	if user == nil {
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		h.Log.Errorf("User not found")
		return
	}

	userFactory := userUsecase.MeUseCaseFactory(h.Log)
	me, err := userFactory.Execute(&userUsecase.IMeUseCaseRequest{
		ID:          uuid.MustParse(user["id"].(string)),
		ChoosedRole: user["choosed_role"].(string),
	})

	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when finding user by ID: %v", err)
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

	organizationId := ctx.Query("organization_id")
	if organizationId == "" {
		organizationId = ""
	}

	if me.User.Employee.OrganizationID != uuid.Nil {
		organizationId = me.User.Employee.OrganizationID.String()
	}

	factory := usecase.FindAllPaginatedUseCaseFactory(h.Log)
	response, err := factory.Execute(&usecase.IFindAllPaginatedUseCaseRequest{
		Page:           page,
		PageSize:       pageSize,
		Search:         search,
		OrganizationID: organizationId,
	})

	if err != nil {
		h.Log.Errorf("Error FindAllPaginated: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", response)
}

func (h *JobHandler) FindById(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-job")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.Log.Errorf("Permission denied")
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

func (h *JobHandler) FindAllJobLevelsPaginated(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-job")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.Log.Errorf("Permission denied")
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

	factory := jobLevelUsecase.FindAllPaginatedUseCaseFactory(h.Log)
	response, err := factory.Execute(&jobLevelUsecase.IFindAllPaginatedUseCaseRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	})

	if err != nil {
		h.Log.Errorf("Error FindAllJobLevelsPaginated: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", response)
}

func (h *JobHandler) FindJobLevelById(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-job")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.Log.Errorf("Permission denied")
		return
	}

	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	factory := jobLevelUsecase.FindByIdUseCaseFactory(h.Log)
	response, err := factory.Execute(&jobLevelUsecase.IFindByIdUseCaseRequest{
		ID: uuid.MustParse(id),
	})

	if err != nil {
		h.Log.Errorf("Error FindJobLevelById: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", response.Job)
}

func (h *JobHandler) GetJobsByJobLevelId(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-job")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.Log.Errorf("Permission denied")
		return
	}

	jobLevelId := ctx.Param("job_level_id")
	if jobLevelId == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "job_level_id is required")
		return
	}

	organizationId := ctx.Query("organization_id")
	if organizationId == "" {
		organizationId = ""
	}

	factory := usecase.GetJobsByJobLevelIDUseCaseFactory(h.Log)
	response, err := factory.Execute(&usecase.IGetJobsByJobLevelIDUseCaseRequest{
		JobLevelID:     jobLevelId,
		OrganizationID: organizationId,
	})

	if err != nil {
		h.Log.Errorf("Error GetJobsByJobLevelId: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", response.Jobs)
}

func (h *JobHandler) FindJobLevelsByOrganizationId(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-job")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.Log.Errorf("Permission denied")
		return
	}

	organizationId := ctx.Param("organization_id")
	if organizationId == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "organization_id is required")
		return
	}

	factory := jobLevelUsecase.FindByOrganizationIDUseCaseFactory(h.Log)
	response, err := factory.Execute(&jobLevelUsecase.IFindByOrganizationIDUseCaseRequest{
		OrganizationID: organizationId,
	})

	if err != nil {
		h.Log.Errorf("Error FindJobLevelsByOrganizationId: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", response.JobLevels)
}

func (h *JobHandler) GetJobsByOrganizationId(ctx *gin.Context) {
	middleware.PermissionApiMiddleware("read-job")(ctx)
	if denied, exists := ctx.Get("permission_denied"); exists && denied.(bool) {
		h.Log.Errorf("Permission denied")
		return
	}

	organizationId := ctx.Param("organization_id")
	if organizationId == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "organization_id is required")
		return
	}

	factory := usecase.GetJobsByOrganizationIDUseCaseFactory(h.Log)
	response, err := factory.Execute(&usecase.IGetJobsByOrganizationIDUseCaseRequest{
		OrganizationID: uuid.MustParse(organizationId),
	})

	if err != nil {
		h.Log.Errorf("Error GetJobsByOrganizationId: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", response.Jobs)
}
