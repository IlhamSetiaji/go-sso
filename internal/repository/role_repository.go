package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	GetAllRoles() (*[]entity.Role, error)
	FindById(id uuid.UUID) (*entity.Role, error)
	StoreRole(role *entity.Role) (*entity.Role, error)
	UpdateRole(role *entity.Role) (*entity.Role, error)
	DeleteRole(id uuid.UUID) error
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
	if err := r.DB.Preload("Application").Preload("Permissions").Preload("Users").Find(&roles).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}
	return &roles, nil
}

func (r *RoleRepository) FindById(id uuid.UUID) (*entity.Role, error) {
	var role entity.Role
	if err := r.DB.Preload("Application").Preload("Permissions").Preload("Users").Where("id = ?", id).First(&role).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) StoreRole(role *entity.Role) (*entity.Role, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("[UserRepository.CreateUser] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Create(role).Error; err != nil {
		tx.Rollback()
		r.Log.Error(err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[RoleRepository.CreateRole] failed to commit transaction: " + err.Error())
		return nil, errors.New("[RoleRepository.CreateRole] failed to commit transaction: " + err.Error())
	}

	return role, nil
}

func (r *RoleRepository) UpdateRole(role *entity.Role) (*entity.Role, error) {
	if err := r.DB.Save(role).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) DeleteRole(id uuid.UUID) error {
	if err := r.DB.Where("id = ?", id).Delete(&entity.Role{}).Error; err != nil {
		r.Log.Error(err)
		return err
	}
	return nil
}
