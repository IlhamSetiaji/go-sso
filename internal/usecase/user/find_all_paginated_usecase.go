package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"
	"errors"

	"github.com/sirupsen/logrus"
)

type IFindAllPaginatedRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type IFindAllPaginatedResponse struct {
	Users *[]entity.User `json:"users"`
	Total int64          `json:"total"`
}

type IFindAllPaginated interface {
	Execute(request *IFindAllPaginatedRequest) (*IFindAllPaginatedResponse, error)
}

type FindAllPaginated struct {
	Log            *logrus.Logger
	UserRepository repository.UserRepositoryInterface
}

func NewFindAllPaginated(
	log *logrus.Logger,
	userRepository repository.UserRepositoryInterface,
) IFindAllPaginated {
	return &FindAllPaginated{
		Log:            log,
		UserRepository: userRepository,
	}
}

func (uc *FindAllPaginated) Execute(req *IFindAllPaginatedRequest) (*IFindAllPaginatedResponse, error) {
	users, total, err := uc.UserRepository.FindAllPaginated(req.Page, req.PageSize)
	if err != nil {
		return nil, errors.New("[FindAllPaginatedUseCase.FindAllPaginated] " + err.Error())
	}

	return &IFindAllPaginatedResponse{
		Users: users,
		Total: total,
	}, nil
}

func FindAllPaginatedUseCaseFactory(log *logrus.Logger) IFindAllPaginated {
	userRepository := repository.UserRepositoryFactory(log)
	return NewFindAllPaginated(log, userRepository)
}
