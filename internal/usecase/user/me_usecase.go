package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	messaging "app/go-sso/internal/messaging/user"
	"app/go-sso/internal/repository"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IMeUseCaseRequest struct {
	ID          uuid.UUID `json:"id"`
	ChoosedRole string    `json:"choosed_role"`
}

type IMeUseCase interface {
	Execute(request *IMeUseCaseRequest) (*IMeUseCaseResponse, error)
}

type IMeUseCaseResponse struct {
	User *response.UserResponse `json:"user"`
}

type MeUseCase struct {
	Log              *logrus.Logger
	UserRepository   repository.IUserRepository
	OrgStructureRepo repository.IOrganizationStructureRepository
}

func NewMeUseCase(log *logrus.Logger, userRepository repository.IUserRepository,
	orgStructureRepo repository.IOrganizationStructureRepository) IMeUseCase {
	return &MeUseCase{
		Log:              log,
		UserRepository:   userRepository,
		OrgStructureRepo: orgStructureRepo,
	}
}

func (u *MeUseCase) Execute(request *IMeUseCaseRequest) (*IMeUseCaseResponse, error) {
	user, err := u.UserRepository.FindById(request.ID)
	if err != nil {
		u.Log.Error("[MeUseCase.Execute] " + err.Error())
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	userProfileMessageFactory := messaging.SendFindUserProfileMessageFactory(u.Log)
	userProfileMessage, err := userProfileMessageFactory.Execute(&messaging.ISendFindUserProfileMessageRequest{
		UserID: user.ID.String(),
	})
	if err != nil {
		u.Log.Error("[MeUseCase.Execute] " + err.Error())
		return nil, err
	}

	var hasUserProfile bool = true
	var userProfileData map[string]interface{}

	userProfile, ok := userProfileMessage.UserProfile["user_profile"].(map[string]interface{})
	if !ok {
		hasUserProfile = false
		userProfileData = map[string]interface{}{}
	} else {
		userProfileData = userProfile
	}

	if hasUserProfile {
		user.VerifiedUserProfile = userProfile["status"].(string)
		user.UserProfile = userProfileData
	} else {
		user.VerifiedUserProfile = "false"
	}

	filteredRoles := []entity.Role{}
	for _, role := range user.Roles {
		if role.Name == request.ChoosedRole {
			filteredRoles = append(filteredRoles, role)
			break
		}
	}
	user.Roles = filteredRoles
	user.ChoosedRole = request.ChoosedRole

	u.Log.Info("[MeUseCase.Execute] User found", user.Employee.EmployeeJob)

	return &IMeUseCaseResponse{
		User: dto.ConvertToSingleUserResponse(user),
	}, nil
}

func MeUseCaseFactory(log *logrus.Logger) IMeUseCase {
	userRepository := repository.UserRepositoryFactory(log)
	orgStructureRepo := repository.OrganizationStructureRepositoryFactory(log)
	return NewMeUseCase(log, userRepository, orgStructureRepo)
}
