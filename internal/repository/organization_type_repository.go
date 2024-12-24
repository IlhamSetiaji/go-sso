package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IOrganizationTypeRepository interface {
	FindAllPaginated(page int, pageSize int, search string) (*[]entity.OrganizationType, int64, error)
	FindById(id uuid.UUID) (*entity.OrganizationType, error)
}

type OrganizationTypeRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewOrganizationTypeRepository(log *logrus.Logger, db *gorm.DB) IOrganizationTypeRepository {
	return &OrganizationTypeRepository{
		Log: log,
		DB:  db,
	}
}

func (r *OrganizationTypeRepository) FindAllPaginated(page int, pageSize int, search string) (*[]entity.OrganizationType, int64, error) {
	var organizationTypes []entity.OrganizationType
	var total int64

	query := r.DB

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&organizationTypes).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Model(&entity.OrganizationType{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return &organizationTypes, total, nil
}

func (r *OrganizationTypeRepository) FindById(id uuid.UUID) (*entity.OrganizationType, error) {
	var organizationType entity.OrganizationType
	err := r.DB.Where("id = ?", id).First(&organizationType).Error
	if err != nil {
		return nil, err
	}
	return &organizationType, nil
}

func OrganizationTypeRepositoryFactory(log *logrus.Logger) IOrganizationTypeRepository {
	db := config.NewDatabase()
	return NewOrganizationTypeRepository(log, db)
}
