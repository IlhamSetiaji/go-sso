package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IOrganizationRepository interface {
	FindAllPaginated(page int, pageSize int, search string) (*[]entity.Organization, int64, error)
	FindById(id uuid.UUID) (*entity.Organization, error)
}

type OrganizationRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewOrganizationRepository(log *logrus.Logger, db *gorm.DB) IOrganizationRepository {
	return &OrganizationRepository{
		Log: log,
		DB:  db,
	}
}

func (r *OrganizationRepository) FindAllPaginated(page int, pageSize int, search string) (*[]entity.Organization, int64, error) {
	var organizations []entity.Organization
	var total int64

	query := r.DB.Preload("OrganizationLocations").Preload("OrganizationStructures").Preload("OrganizationType")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&organizations).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Model(&entity.Organization{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return &organizations, total, nil
}

func (r *OrganizationRepository) FindById(id uuid.UUID) (*entity.Organization, error) {
	var organization entity.Organization
	err := r.DB.Preload("OrganizationLocations").Preload("OrganizationStructures").Preload("OrganizationType").Where("id = ?", id).First(&organization).Error
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func OrganizationRepositoryFactory(log *logrus.Logger) IOrganizationRepository {
	db := config.NewDatabase()
	return NewOrganizationRepository(log, db)
}
