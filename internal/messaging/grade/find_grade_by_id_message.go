package messaging

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IFindGradeByIDMessageRequest struct {
	ID string `json:"id"`
}

type IFindGradeByIDMessageResponse struct {
	Grade *response.GradeResponse `json:"grade"`
}

type IFindGradeByIDMessage interface {
	Execute(req *IFindGradeByIDMessageRequest) (*IFindGradeByIDMessageResponse, error)
}

type FindGradeByIDMessage struct {
	Log             *logrus.Logger
	GradeRepository repository.IGradeRepository
}

func NewFindGradeByIDMessage(
	log *logrus.Logger,
	gradeRepository repository.IGradeRepository,
) IFindGradeByIDMessage {
	return &FindGradeByIDMessage{
		Log:             log,
		GradeRepository: gradeRepository,
	}
}

func (uc *FindGradeByIDMessage) Execute(req *IFindGradeByIDMessageRequest) (*IFindGradeByIDMessageResponse, error) {
	grade, err := uc.GradeRepository.FindByKeys(map[string]interface{}{
		"id": req.ID,
	})
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}

	if grade == nil {
		uc.Log.Error("grade not found")
		return nil, nil
	}

	gradeResp := dto.ConvertToSingleGradeResponse(grade)

	return &IFindGradeByIDMessageResponse{
		Grade: gradeResp,
	}, nil
}

func FindGradeByIDMessageFactory(log *logrus.Logger) *FindGradeByIDMessage {
	gradeRepository := repository.GradeRepositoryFactory(log)
	return &FindGradeByIDMessage{
		GradeRepository: gradeRepository,
	}
}
