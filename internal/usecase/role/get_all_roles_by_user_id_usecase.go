package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IGetAllRolesByUserIdUsecaseRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

type IGetAllRolesByUserIdUsecaseResponse struct {
	Roles *[]entity.Role `json:"roles"`
}

type IGetAllRolesByUserIdUsecase interface {
	Execute(request *IGetAllRolesByUserIdUsecaseRequest) (*IGetAllRolesByUserIdUsecaseResponse, error)
}

type GetAllRolesByUserIdUsecase struct {
	Log      *logrus.Logger
	UserRepo repository.IUserRepository
	RoleRepo repository.IRoleRepository
}

func GetAllRolesByUserIdUsecaseFactory(log *logrus.Logger) IGetAllRolesByUserIdUsecase {
	userRepo := repository.UserRepositoryFactory(log)
	roleRepo := repository.RoleRepositoryFactory(log)
	return &GetAllRolesByUserIdUsecase{
		Log:      log,
		UserRepo: userRepo,
		RoleRepo: roleRepo,
	}
}

func (u *GetAllRolesByUserIdUsecase) Execute(request *IGetAllRolesByUserIdUsecaseRequest) (*IGetAllRolesByUserIdUsecaseResponse, error) {
	roles, err := u.RoleRepo.GetAllRolesInUserID(uuid.MustParse(request.UserID))
	if err != nil {
		u.Log.Error(err)
		return nil, err
	}

	return &IGetAllRolesByUserIdUsecaseResponse{
		Roles: roles,
	}, nil
}
