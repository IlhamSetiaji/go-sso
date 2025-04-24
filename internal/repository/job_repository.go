package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IJobRepository interface {
	FindAllPaginated(page int, pageSize int, search string, orgStructureIds []string, filter map[string]interface{}) (*[]entity.Job, int64, error)
	FindById(id uuid.UUID) (*entity.Job, error)
	GetAll() (*[]entity.Job, error)
	FindAllJobs(includedIDs []string) (*[]entity.Job, error)
	GetJobsByOrganizationStructureIDs(organizationStructureIDs []uuid.UUID) (*[]entity.Job, error)
	FindAllChildren(parentID uuid.UUID) ([]entity.Job, error)
	FindParent(id uuid.UUID) (*entity.Job, error)
	FindAllByKeys(keys map[string]interface{}) (*[]entity.Job, error)
}

type JobRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewJobRepository(log *logrus.Logger, db *gorm.DB) IJobRepository {
	return &JobRepository{
		Log: log,
		DB:  db,
	}
}

func (r *JobRepository) GetAll() (*[]entity.Job, error) {
	var jobs []entity.Job
	if err := r.DB.Preload("OrganizationStructure.Organization").Preload("OrganizationStructure.JobLevel").Preload("JobLevel").Find(&jobs).Error; err != nil {
		return nil, err
	}

	return &jobs, nil
}

func (r *JobRepository) FindAllJobs(includedIDs []string) (*[]entity.Job, error) {
	var jobs []entity.Job
	if err := r.DB.Preload("OrganizationStructure.Organization").Preload("OrganizationStructure.JobLevel").Preload("JobLevel").Where("id IN ?", includedIDs).Find(&jobs).Error; err != nil {
		return nil, err
	}

	return &jobs, nil
}

func (r *JobRepository) FindAllPaginated(page int, pageSize int, search string, orgStructureIds []string, filter map[string]interface{}) (*[]entity.Job, int64, error) {
	var jobs []entity.Job
	var total int64

	query := r.DB.Preload("OrganizationStructure.Organization").Preload("OrganizationStructure.JobLevel").Preload("JobLevel")

	if search != "" {
		query = query.Where("jobs.name ILIKE ?", "%"+search+"%")
	}
	if filter["name"] != nil {
		query = query.Where("jobs.name ILIKE ?", "%"+filter["name"].(string)+"%")
	}
	if filter["parent.name"] != nil {
		query = query.Joins("JOIN jobs AS parent ON jobs.parent_id = parent.id").Where("parent.name ILIKE ?", "%"+filter["parent.name"].(string)+"%")
	}
	if filter["organization_name"] != nil {
		query = query.Joins("JOIN organization_structures AS org_struct ON jobs.organization_structure_id = org_struct.id").Joins("JOIN organizations AS org ON org_struct.organization_id = org.id").Where("org.name ILIKE ?", "%"+filter["organization_name"].(string)+"%")
	}

	if len(orgStructureIds) > 0 {
		query = query.Where("jobs.organization_structure_id IN ?", orgStructureIds)
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Model(&entity.Job{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return &jobs, total, nil
}

func (r *JobRepository) FindById(id uuid.UUID) (*entity.Job, error) {
	var job entity.Job
	err := r.DB.Preload("OrganizationStructure.Organization").Preload("OrganizationStructure.JobLevel").Preload("JobLevel").Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepository) GetJobsByOrganizationStructureIDs(organizationStructureIDs []uuid.UUID) (*[]entity.Job, error) {
	var jobs []entity.Job
	err := r.DB.Preload("OrganizationStructure.Organization").Preload("OrganizationStructure.JobLevel").Preload("JobLevel").Where("organization_structure_id IN ?", organizationStructureIDs).Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return &jobs, nil
}

func (r *JobRepository) FindAllChildren(parentID uuid.UUID) ([]entity.Job, error) {
	var children []entity.Job
	if err := r.DB.Preload("OrganizationStructure.Organization").Preload("OrganizationStructure.JobLevel").Preload("JobLevel").Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
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

func (r *JobRepository) FindParent(id uuid.UUID) (*entity.Job, error) {
	var job entity.Job
	if err := r.DB.Where("id = ?", id).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &job, nil
}

func (r *JobRepository) FindAllByKeys(keys map[string]interface{}) (*[]entity.Job, error) {
	var jobs []entity.Job
	if err := r.DB.Preload("Organization").Preload("JobLevel.OrganizationStructures").Where(keys).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return &jobs, nil
}

func JobRepositoryFactory(log *logrus.Logger) IJobRepository {
	db := config.NewDatabase()
	return NewJobRepository(log, db)
}
