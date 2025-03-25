package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IGradeRepository interface {
	FindAllByJobLevelID(jobLevelID uuid.UUID) (*[]entity.Grade, error)
	FindByKeys(keys map[string]interface{}) (*entity.Grade, error)
}

type GradeRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewGradeRepository(log *logrus.Logger, db *gorm.DB) IGradeRepository {
	return &GradeRepository{
		Log: log,
		DB:  db,
	}
}

func GradeRepositoryFactory(log *logrus.Logger) *GradeRepository {
	db := config.NewDatabase()
	return &GradeRepository{
		Log: log,
		DB:  db,
	}
}

func (g *GradeRepository) FindAllByJobLevelID(jobLevelID uuid.UUID) (*[]entity.Grade, error) {
	var grades []entity.Grade
	if err := g.DB.Preload("JobLevel").Where("job_level_id = ?", jobLevelID).Find(&grades).Error; err != nil {
		return nil, err
	}

	return &grades, nil
}

func (g *GradeRepository) FindByKeys(keys map[string]interface{}) (*entity.Grade, error) {
	var grade entity.Grade
	if err := g.DB.Preload("JobLevel").Where(keys).First(&grade).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &grade, nil
}
