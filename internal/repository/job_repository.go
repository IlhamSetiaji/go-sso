package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IJobRepository interface {
	FindAllPaginated(page int, pageSize int, search string) (*[]entity.Job, int64, error)
	FindById(id uuid.UUID) (*entity.Job, error)
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

func (r *JobRepository) FindAllPaginated(page int, pageSize int, search string) (*[]entity.Job, int64, error) {
	var jobs []entity.Job
	var total int64

	query := r.DB

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
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
	err := r.DB.Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func JobRepositoryFactory(log *logrus.Logger) IJobRepository {
	db := config.NewDatabase()
	return NewJobRepository(log, db)
}
