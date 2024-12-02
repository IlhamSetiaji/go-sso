package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	GetAllRoles() (*[]entity.Role, error)
}

type RoleRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewRoleRepository(log *logrus.Logger, db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		Log: log,
		DB:  db,
	}
}

func RoleRepositoryFactory(log *logrus.Logger) *RoleRepository {
	db := config.NewDatabase()
	return NewRoleRepository(log, db)
}

func (r *RoleRepository) GetAllRoles() (*[]entity.Role, error) {
	var roles []entity.Role
	if err := r.DB.Preload("Application").Find(&roles).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}
	return &roles, nil
}
