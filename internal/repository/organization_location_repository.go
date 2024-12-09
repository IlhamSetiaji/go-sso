package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IOrganizationLocationRepository interface {
	FindAllPaginated(page int, pageSize int, search string) (*[]entity.OrganizationLocation, int64, error)
	FindById(id uuid.UUID) (*entity.OrganizationLocation, error)
	FindByOrganizationID(organizationID uuid.UUID) (*[]entity.OrganizationLocation, error)
}

type OrganizationLocationRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewOrganizationLocationRepository(log *logrus.Logger, db *gorm.DB) IOrganizationLocationRepository {
	return &OrganizationLocationRepository{
		Log: log,
		DB:  db,
	}
}

func (r *OrganizationLocationRepository) FindAllPaginated(page int, pageSize int, search string) (*[]entity.OrganizationLocation, int64, error) {
	var organizationLocations []entity.OrganizationLocation
	var total int64

	query := r.DB

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&organizationLocations).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Model(&entity.OrganizationLocation{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return &organizationLocations, total, nil
}

func (r *OrganizationLocationRepository) FindById(id uuid.UUID) (*entity.OrganizationLocation, error) {
	var organizationLocation entity.OrganizationLocation
	err := r.DB.Where("id = ?", id).First(&organizationLocation).Error
	if err != nil {
		return nil, err
	}
	return &organizationLocation, nil
}

func (r *OrganizationLocationRepository) FindByOrganizationID(organizationID uuid.UUID) (*[]entity.OrganizationLocation, error) {
	var organizationLocations []entity.OrganizationLocation
	err := r.DB.Where("organization_id = ?", organizationID).Find(&organizationLocations).Error
	if err != nil {
		return nil, err
	}
	return &organizationLocations, nil
}

func OrganizationLocationRepositoryFactory(log *logrus.Logger) IOrganizationLocationRepository {
	db := config.NewDatabase()
	return NewOrganizationLocationRepository(log, db)
}
