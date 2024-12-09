package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IOrganizationStructureRepository interface {
	FindAllPaginated(page int, pageSize int, search string) (*[]entity.OrganizationStructure, int64, error)
	FindById(id uuid.UUID) (*entity.OrganizationStructure, error)
	GetOrganizationSructuresByJobLevelID(jobLevelID uuid.UUID) (*[]entity.OrganizationStructure, error)
	FindAllChildren(parentID uuid.UUID) ([]entity.OrganizationStructure, error)
	// GetDB() *gorm.DB
}

type OrganizationStructureRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewOrganizationStructureRepository(log *logrus.Logger, db *gorm.DB) IOrganizationStructureRepository {
	return &OrganizationStructureRepository{
		Log: log,
		DB:  db,
	}
}

func (r *OrganizationStructureRepository) FindAllPaginated(page int, pageSize int, search string) (*[]entity.OrganizationStructure, int64, error) {
	var organizationStructures []entity.OrganizationStructure
	var total int64

	query := r.DB.Preload("Organization.OrganizationType").Preload("JobLevel")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&organizationStructures).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Model(&entity.OrganizationStructure{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return &organizationStructures, total, nil
}

func (r *OrganizationStructureRepository) FindAllChildren(parentID uuid.UUID) ([]entity.OrganizationStructure, error) {
	var children []entity.OrganizationStructure
	if err := r.DB.Preload("Organization.OrganizationType").Preload("JobLevel").Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
		return nil, err
	}

	for i := range children {
		subChildren, err := r.FindAllChildren(children[i].ID)
		if err != nil {
			return nil, err
		}
		children[i].Children = subChildren
	}

	return children, nil
}

func (r *OrganizationStructureRepository) FindById(id uuid.UUID) (*entity.OrganizationStructure, error) {
	var organizationStructure entity.OrganizationStructure
	err := r.DB.Preload("Organization.OrganizationType").Preload("JobLevel").Where("id = ?", id).First(&organizationStructure).Error
	if err != nil {
		return nil, err
	}
	return &organizationStructure, nil
}

func (r *OrganizationStructureRepository) GetOrganizationSructuresByJobLevelID(jobLevelID uuid.UUID) (*[]entity.OrganizationStructure, error) {
	var organizationStructures []entity.OrganizationStructure
	err := r.DB.Preload("Organization").Preload("JobLevel").Preload("Jobs").Where("job_level_id = ?", jobLevelID).Find(&organizationStructures).Error
	if err != nil {
		return nil, err
	}
	return &organizationStructures, nil
}

// func (r *OrganizationStructureRepository) GetDB() *gorm.DB {
// 	return r.DB
// }

func OrganizationStructureRepositoryFactory(log *logrus.Logger) IOrganizationStructureRepository {
	db := config.NewDatabase()
	return NewOrganizationStructureRepository(log, db)
}
