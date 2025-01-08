package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IJobLevelRepository interface {
	FindAllPaginated(page int, pageSize int, search string) (*[]entity.JobLevel, int64, error)
	FindById(id uuid.UUID) (*entity.JobLevel, error)
	FindByIds(ids []uuid.UUID) (*[]entity.JobLevel, error)
}

type JobLevelRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewJobLevelRepository(log *logrus.Logger, db *gorm.DB) IJobLevelRepository {
	return &JobLevelRepository{
		Log: log,
		DB:  db,
	}
}

func (r *JobLevelRepository) FindAllPaginated(page int, pageSize int, search string) (*[]entity.JobLevel, int64, error) {
	var jobLevels []entity.JobLevel
	var total int64

	query := r.DB.Preload("OrganizationStructures")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&jobLevels).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Model(&entity.JobLevel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return &jobLevels, total, nil
}

func (r *JobLevelRepository) FindById(id uuid.UUID) (*entity.JobLevel, error) {
	var jobLevel entity.JobLevel
	err := r.DB.Preload("OrganizationStructures").Where("id = ?", id).First(&jobLevel).Error
	if err != nil {
		return nil, err
	}
	return &jobLevel, nil
}

func (r *JobLevelRepository) FindByIds(ids []uuid.UUID) (*[]entity.JobLevel, error) {
	var jobLevels []entity.JobLevel
	err := r.DB.Preload("OrganizationStructures").Where("id IN (?)", ids).Order("level ASC").Find(&jobLevels).Error
	if err != nil {
		return nil, err
	}
	return &jobLevels, nil
}

func JobLevelRepositoryFactory(log *logrus.Logger) IJobLevelRepository {
	db := config.NewDatabase()
	return NewJobLevelRepository(log, db)
}
