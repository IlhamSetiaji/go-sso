package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IStoreRoleUseCaseRequest struct {
	Role          *entity.Role `json:"role"`
	ApplicationID uuid.UUID    `json:"application_id"`
}

type IStoreRoleUseCaseResponse struct {
	Role *entity.Role `json:"role"`
}

type IStoreRoleUseCase interface {
	Execute(request *IStoreRoleUseCaseRequest) (*IStoreRoleUseCaseResponse, error)
}

type StoreRoleUseCase struct {
	Log                  *logrus.Logger
	RoleRepository       repository.IRoleRepository
	PermissionRepository repository.IPermissionRepository
}

func NewStoreRoleUseCase(log *logrus.Logger, roleRepository repository.IRoleRepository, permissionRepo repository.IPermissionRepository) IStoreRoleUseCase {
	return &StoreRoleUseCase{
		Log:                  log,
		RoleRepository:       roleRepository,
		PermissionRepository: permissionRepo,
	}
}

func (uc *StoreRoleUseCase) Execute(request *IStoreRoleUseCaseRequest) (*IStoreRoleUseCaseResponse, error) {
	uc.Log.Info("StoreRoleUseCase.Execute")

	role, err := uc.RoleRepository.StoreRole(&entity.Role{
		Name:          request.Role.Name,
		ApplicationID: request.ApplicationID,
		GuardName:     request.Role.GuardName,
		Status:        request.Role.Status,
	})
	if err != nil {
		return nil, err
	}

	permissionNames := []string{
		"read-user",
		"read-role",
		"read-permission",
		"read-application",
		"read-client",
		"assign-role",
		"assign-permission",
		"read-organization",
		"read-organization-location",
		"read-job-level",
		"read-organization-structure",
		"read-job",
		"read-employee",
		"read-employee-job",
	}

	permissions, err := uc.PermissionRepository.GetAllPermissionsByNames(permissionNames)
	if err != nil {
		return nil, err
	}

	var permissionIDs []string

	for _, permission := range *permissions {
		permissionIDs = append(permissionIDs, permission.ID.String())
	}

	_, err = uc.RoleRepository.AssignRoleToPermissions(role, permissionIDs)
	if err != nil {
		uc.Log.Error(err)
		return nil, err
	}

	return &IStoreRoleUseCaseResponse{
		Role: role,
	}, nil
}

func StoreRoleUseCaseFactory(log *logrus.Logger) IStoreRoleUseCase {
	roleRepository := repository.RoleRepositoryFactory(log)
	permissionRepo := repository.PermissionRepositoryFactory(log)
	return &StoreRoleUseCase{
		Log:                  log,
		RoleRepository:       roleRepository,
		PermissionRepository: permissionRepo,
	}
}
