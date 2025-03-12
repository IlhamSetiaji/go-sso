package usecase

import (
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/repository"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IUploadLogoOrganizationUseCaseRequest struct {
	ID       string                `form:"id"`
	Logo     *multipart.FileHeader `form:"logo"`
	LogoPath string                `form:"logo_path"`
}

type IUploadLogoOrganizationUseCaseResponse struct {
	Organization *response.OrganizationResponse `json:"organization"`
}

type IUploadLogoOrganizationUseCase interface {
	Execute(request *IUploadLogoOrganizationUseCaseRequest) (*IUploadLogoOrganizationUseCaseResponse, error)
}

type UploadLogoOrganizationUseCase struct {
	Log                    *logrus.Logger
	OrganizationRepository repository.IOrganizationRepository
}

func NewUploadLogoOrganizationUseCase(
	log *logrus.Logger,
	organizationRepository repository.IOrganizationRepository,
) IUploadLogoOrganizationUseCase {
	return &UploadLogoOrganizationUseCase{
		Log:                    log,
		OrganizationRepository: organizationRepository,
	}
}

func (uc *UploadLogoOrganizationUseCase) Execute(req *IUploadLogoOrganizationUseCaseRequest) (*IUploadLogoOrganizationUseCaseResponse, error) {
	parsedOrgID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, err
	}

	organization, err := uc.OrganizationRepository.FindByIdOnly(parsedOrgID)
	if err != nil {
		return nil, err
	}
	if organization == nil {
		return nil, nil
	}

	organization.Logo = req.LogoPath
	_, err = uc.OrganizationRepository.UpdateOrganization(organization)
	if err != nil {
		return nil, err
	}

	findByID, err := uc.OrganizationRepository.FindByIdOnly(parsedOrgID)
	if err != nil {
		return nil, err
	}
	if findByID == nil {
		return nil, nil
	}

	resp := dto.ConvertToSingleOrganizationResponse(findByID)

	return &IUploadLogoOrganizationUseCaseResponse{
		Organization: resp,
	}, nil
}

func UploadLogoOrganizationUseCaseFactory(log *logrus.Logger) IUploadLogoOrganizationUseCase {
	organizationRepository := repository.OrganizationRepositoryFactory(log)
	return NewUploadLogoOrganizationUseCase(log, organizationRepository)
}
