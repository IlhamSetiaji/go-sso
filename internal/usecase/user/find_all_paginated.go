package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"
	"errors"
	"log"
)

type FindAllPaginated struct {
	Log            *log.Logger
	UserRepository repository.UserRepositoryInterface
}

type FindAllPaginatedResponse struct {
	Users *[]entity.User `json:"users"`
	Total int64          `json:"total"`
}

type FindAllPaginatedRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func FindAllPaginatedFactory(
	log *log.Logger,
) *FindAllPaginated {
	return &FindAllPaginated{
		Log:            log,
		UserRepository: repository.UserRepositoryFactory(log),
	}
}

func (uc *FindAllPaginated) FindAllPaginated(req *FindAllPaginatedRequest) (*FindAllPaginatedResponse, error) {
	users, total, err := uc.UserRepository.FindAllPaginated(req.Page, req.PageSize)
	if err != nil {
		return nil, errors.New("[FindAllPaginated.FindAllPaginated] " + err.Error())
	}

	return &FindAllPaginatedResponse{
		Users: users,
		Total: total,
	}, nil
}
