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
	FindByIdOnly(id uuid.UUID) (*entity.Organization, error)
	FindAllOrganizations() (*[]entity.Organization, error)
	FindByIDs(ids []uuid.UUID) (*[]entity.Organization, error)
	UpdateOrganization(ent *entity.Organization) (*entity.Organization, error)
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
		query = query.Where("name ILIKE ?", "%"+search+"%")
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
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &organization, nil
}

func (r *OrganizationRepository) FindByIdOnly(id uuid.UUID) (*entity.Organization, error) {
	var organization entity.Organization
	err := r.DB.Preload("OrganizationType").Where("id = ?", id).First(&organization).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &organization, nil
}

func (r *OrganizationRepository) FindByIDs(ids []uuid.UUID) (*[]entity.Organization, error) {
	var organizations []entity.Organization
	err := r.DB.Preload("OrganizationLocations").Preload("OrganizationStructures").Preload("OrganizationType").Where("id IN ?", ids).Find(&organizations).Error
	if err != nil {
		return nil, err
	}
	return &organizations, nil
}

func (r *OrganizationRepository) FindAllOrganizations() (*[]entity.Organization, error) {
	var organizations []entity.Organization
	err := r.DB.Preload("OrganizationLocations").Preload("OrganizationStructures").Preload("OrganizationType").Find(&organizations).Error
	if err != nil {
		return nil, err
	}
	return &organizations, nil
}

func (r *OrganizationRepository) UpdateOrganization(ent *entity.Organization) (*entity.Organization, error) {
	err := r.DB.Where("id = ?", ent.ID).Updates(&ent).Error
	if err != nil {
		return nil, err
	}
	return ent, nil
}

func OrganizationRepositoryFactory(log *logrus.Logger) IOrganizationRepository {
	db := config.NewDatabase()
	return NewOrganizationRepository(log, db)
}
