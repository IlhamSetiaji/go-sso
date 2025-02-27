package messaging

import (
	"app/go-sso/internal/repository"

	"github.com/sirupsen/logrus"
)

type IChartEmployeeOrganizationStructureMessageResponse struct {
	Labels   []string `json:"labels"`
	Datasets []int    `json:"datasets"`
}

type IChartEmployeeOrganizationStructureMessage interface {
	Execute() (*IChartEmployeeOrganizationStructureMessageResponse, error)
}

type ChartEmployeeOrganizationStructureMessage struct {
	Log                             *logrus.Logger
	Repository                      repository.IEmployeeRepository
	OrganizationStructureRepository repository.IOrganizationStructureRepository
}

func NewChartEmployeeOrganizationStructureMessage(
	log *logrus.Logger,
	repository repository.IEmployeeRepository,
	orgStructureRepository repository.IOrganizationStructureRepository,
) IChartEmployeeOrganizationStructureMessage {
	return &ChartEmployeeOrganizationStructureMessage{
		Log:                             log,
		Repository:                      repository,
		OrganizationStructureRepository: orgStructureRepository,
	}
}

func (m *ChartEmployeeOrganizationStructureMessage) Execute() (*IChartEmployeeOrganizationStructureMessageResponse, error) {
	organizationStructureIDs, err := m.Repository.GetOrganizationStructureIdDistinct()
	if err != nil {
		return nil, err
	}

	var labels []string
	var datasets []int

	for _, organizationStructureID := range organizationStructureIDs {
		orgStructure, err := m.OrganizationStructureRepository.FindByIdOnly(organizationStructureID)
		if err != nil {
			return nil, err
		}
		if orgStructure == nil {
			continue
		}
		count, err := m.Repository.CountByOrganizationStructureID(organizationStructureID)
		if err != nil {
			return nil, err
		}

		labels = append(labels, orgStructure.Name)
		datasets = append(datasets, count)
	}

	return &IChartEmployeeOrganizationStructureMessageResponse{
		Labels:   labels,
		Datasets: datasets,
	}, nil
}

func ChartEmployeeOrganizationStructureMessageFactory(log *logrus.Logger) IChartEmployeeOrganizationStructureMessage {
	empRepository := repository.EmployeeRepositoryFactory(log)
	orgStructureRepository := repository.OrganizationStructureRepositoryFactory(log)
	return NewChartEmployeeOrganizationStructureMessage(log, empRepository, orgStructureRepository)
}
