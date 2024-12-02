package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IApplicationRepository interface {
	GetAllApplications() (*[]entity.Application, error)
}

type ApplicationRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewApplicationRepository(log *logrus.Logger, db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{
		Log: log,
		DB:  db,
	}
}

func ApplicationRepositoryFactory(log *logrus.Logger) *ApplicationRepository {
	db := config.NewDatabase()
	return NewApplicationRepository(log, db)
}

func (r *ApplicationRepository) GetAllApplications() (*[]entity.Application, error) {
	var applications []entity.Application
	if err := r.DB.Find(&applications).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}
	return &applications, nil
}
