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
	FindAllOrgStructuresByOrganizationID(organizationID uuid.UUID) (*[]entity.OrganizationStructure, error)
	FindById(id uuid.UUID) (*entity.OrganizationStructure, error)
	GetOrganizationSructuresByJobLevelID(jobLevelID uuid.UUID) (*[]entity.OrganizationStructure, error)
	GetOrganizationSructuresByJobLevelIDAndOrganizationID(jobLevelID uuid.UUID, organizationID uuid.UUID) (*[]entity.OrganizationStructure, error)
	FindByOrganizationId(organizationID uuid.UUID) (*[]entity.OrganizationStructure, error)
	FindAllChildren(parentID uuid.UUID) ([]entity.OrganizationStructure, error)
	FindAllParents(parentID uuid.UUID) ([]entity.OrganizationStructure, error)
	FindAllChildrenIDs(parentID uuid.UUID) ([]uuid.UUID, error)
	FindParent(id uuid.UUID) (*entity.OrganizationStructure, error)
	GetOrganizationSructuresByOrganizationID(organizationID uuid.UUID) (*[]entity.OrganizationStructure, error)
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
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&organizationStructures).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Model(&entity.OrganizationStructure{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return &organizationStructures, total, nil
}

func (r *OrganizationStructureRepository) FindAllOrgStructuresByOrganizationID(organizationID uuid.UUID) (*[]entity.OrganizationStructure, error) {
	var organizationStructures []entity.OrganizationStructure
	err := r.DB.Preload("Organization.OrganizationType").Preload("JobLevel").Where("organization_id = ?", organizationID).Find(&organizationStructures).Error
	if err != nil {
		return nil, err
	}
	return &organizationStructures, nil
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

func (r *OrganizationStructureRepository) FindAllParents(parentID uuid.UUID) ([]entity.OrganizationStructure, error) {
	var parents []entity.OrganizationStructure
	if err := r.DB.Preload("Organization.OrganizationType").Preload("JobLevel").Where("id = ?", parentID).Find(&parents).Error; err != nil {
		return nil, err
	}

	for i := range parents {
		if parents[i].ParentID != nil {
			subParents, err := r.FindAllParents(*parents[i].ParentID)
			if err != nil {
				return nil, err
			}
			parents[i].Parents = subParents
		}
	}

	return parents, nil
}

func (r *OrganizationStructureRepository) FindAllChildrenIDs(parentID uuid.UUID) ([]uuid.UUID, error) {
	var children []entity.OrganizationStructure
	childrenIDMap := make(map[uuid.UUID]bool)
	var childrenIDs []uuid.UUID

	if err := r.DB.Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
		return nil, err
	}

	for _, child := range children {
		if !childrenIDMap[child.ID] {
			childrenIDMap[child.ID] = true
			childrenIDs = append(childrenIDs, child.ID)
		}
		subChildrenIDs, err := r.FindAllChildrenIDs(child.ID)
		if err != nil {
			return nil, err
		}
		for _, subChildID := range subChildrenIDs {
			if !childrenIDMap[subChildID] {
				childrenIDMap[subChildID] = true
				childrenIDs = append(childrenIDs, subChildID)
			}
		}
	}

	return childrenIDs, nil
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

func (r *OrganizationStructureRepository) GetOrganizationSructuresByJobLevelIDAndOrganizationID(jobLevelID uuid.UUID, organizationID uuid.UUID) (*[]entity.OrganizationStructure, error) {
	var organizationStructures []entity.OrganizationStructure
	err := r.DB.Preload("Organization").Preload("JobLevel").Preload("Jobs").Where("job_level_id = ? AND organization_id = ?", jobLevelID, organizationID).Find(&organizationStructures).Error
	if err != nil {
		return nil, err
	}
	return &organizationStructures, nil
}

func (r *OrganizationStructureRepository) FindByOrganizationId(organizationID uuid.UUID) (*[]entity.OrganizationStructure, error) {
	var organizationStructures []entity.OrganizationStructure
	err := r.DB.Preload("Organization.OrganizationType").Preload("JobLevel").Where("organization_id = ?", organizationID).Find(&organizationStructures).Error
	if err != nil {
		return nil, err
	}
	return &organizationStructures, nil
}

func (r *OrganizationStructureRepository) FindParent(id uuid.UUID) (*entity.OrganizationStructure, error) {
	var organizationStructure entity.OrganizationStructure
	if err := r.DB.Where("id = ?", id).First(&organizationStructure).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &organizationStructure, nil
}

func (r *OrganizationStructureRepository) GetOrganizationSructuresByOrganizationID(organizationID uuid.UUID) (*[]entity.OrganizationStructure, error) {
	var organizationStructures []entity.OrganizationStructure
	err := r.DB.Preload("Organization.OrganizationType").Preload("JobLevel").Where("organization_id = ?", organizationID).Find(&organizationStructures).Error
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
