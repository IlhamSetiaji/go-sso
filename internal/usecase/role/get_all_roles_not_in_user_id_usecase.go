package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IGetAllRolesNotInUserIdUsecaseRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

type IGetAllRolesNotInUserIdUsecaseResponse struct {
	Roles *[]entity.Role `json:"roles"`
}

type IGetAllRolesNotInUserIdUsecase interface {
	Execute(request *IGetAllRolesNotInUserIdUsecaseRequest) (*IGetAllRolesNotInUserIdUsecaseResponse, error)
}

type GetAllRolesNotInUserIdUsecase struct {
	Log      *logrus.Logger
	UserRepo repository.IUserRepository
	RoleRepo repository.IRoleRepository
}

func GetAllRolesNotInUserIdUsecaseFactory(log *logrus.Logger) IGetAllRolesNotInUserIdUsecase {
	userRepo := repository.UserRepositoryFactory(log)
	roleRepo := repository.RoleRepositoryFactory(log)
	return &GetAllRolesNotInUserIdUsecase{
		Log:      log,
		UserRepo: userRepo,
		RoleRepo: roleRepo,
	}
}

func (u *GetAllRolesNotInUserIdUsecase) Execute(request *IGetAllRolesNotInUserIdUsecaseRequest) (*IGetAllRolesNotInUserIdUsecaseResponse, error) {
	roles, err := u.RoleRepo.GetAllRolesNotInUserID(uuid.MustParse(request.UserID))
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IGetAllRolesNotInUserIdUsecaseResponse{
		Roles: roles,
	}, nil
}
