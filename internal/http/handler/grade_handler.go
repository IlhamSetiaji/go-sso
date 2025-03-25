package handler

import (
	usecase "app/go-sso/internal/usecase/grade"
	"app/go-sso/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IGradeHandler interface {
	FindAllByJobLevelID(ctx *gin.Context)
}

type GradeHandler struct {
	Viper    *viper.Viper
	Log      *logrus.Logger
	Validate *validator.Validate
}

func NewGradeHandler(viper *viper.Viper, log *logrus.Logger, validate *validator.Validate) IGradeHandler {
	return &GradeHandler{
		Viper:    viper,
		Log:      log,
		Validate: validate,
	}
}

func GradeHandlerFactory(viper *viper.Viper, log *logrus.Logger, validate *validator.Validate) IGradeHandler {
	return NewGradeHandler(viper, log, validate)
}

func (h *GradeHandler) FindAllByJobLevelID(ctx *gin.Context) {
	jobLevelID := ctx.Param("job_level_id")
	if jobLevelID == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "bad request", "job_level_id is required")
		return
	}

	req := &usecase.IGetAllByJobLevelIDUseCaseRequest{
		JobLevelID: jobLevelID,
	}

	ucFactory := usecase.GetAllByJobLevelIDUseCaseFactory(h.Log)
	resp, err := ucFactory.Execute(req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", resp.Grade)
}
