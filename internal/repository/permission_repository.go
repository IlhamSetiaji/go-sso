package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IPermissionRepository interface {
	GetAllPermissions() (*[]entity.Permission, error)
	GetAllPermissionsByRoleID(roleID uuid.UUID) (*[]entity.Permission, error)
	GetAllPermissionsNotInRoleID(roleID uuid.UUID) (*[]entity.Permission, error)
	GetAllPermissionsByNames(names []string) (*[]entity.Permission, error)
	FindById(id uuid.UUID) (*entity.Permission, error)
	StorePermission(permission *entity.Permission) (*entity.Permission, error)
	UpdatePermission(permission *entity.Permission) (*entity.Permission, error)
	DeletePermission(id uuid.UUID) error
}

type PermissionRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewPermissionRepository(log *logrus.Logger, db *gorm.DB) IPermissionRepository {
	return &PermissionRepository{
		Log: log,
		DB:  db,
	}
}

func PermissionRepositoryFactory(log *logrus.Logger) IPermissionRepository {
	db := config.NewDatabase()
	return NewPermissionRepository(log, db)
}

func (r *PermissionRepository) GetAllPermissions() (*[]entity.Permission, error) {
	var permissions []entity.Permission
	if err := r.DB.Preload("Application").Preload("Roles").Find(&permissions).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}
	return &permissions, nil
}

func (r *PermissionRepository) GetAllPermissionsByRoleID(roleID uuid.UUID) (*[]entity.Permission, error) {
	var permissions []entity.Permission

	if err := r.DB.Preload("Application").Preload("Roles").Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").Where("role_permissions.role_id = ?", roleID).Find(&permissions).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}

	return &permissions, nil
}

func (r *PermissionRepository) GetAllPermissionsNotInRoleID(roleID uuid.UUID) (*[]entity.Permission, error) {
	var permissions []entity.Permission

	if err := r.DB.Preload("Application").Preload("Roles").Not("id IN (?)", r.DB.Table("role_permissions").Select("permission_id").Where("role_id = ?", roleID)).Find(&permissions).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}

	return &permissions, nil
}

func (r *PermissionRepository) FindById(id uuid.UUID) (*entity.Permission, error) {
	var permission entity.Permission
	if err := r.DB.Where("id = ?", id).First(&permission).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) StorePermission(permission *entity.Permission) (*entity.Permission, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("[PermissionRepository.StorePermission] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Create(permission).Error; err != nil {
		tx.Rollback()
		r.Log.Error(err)
		return nil, err
	}

	tx.Commit()
	return permission, nil
}

func (r *PermissionRepository) UpdatePermission(permission *entity.Permission) (*entity.Permission, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("[PermissionRepository.UpdatePermission] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Model(permission).Where("id = ?", permission.ID).Updates(permission).Error; err != nil {
		tx.Rollback()
		r.Log.Error(err)
		return nil, err
	}

	tx.Commit()
	return permission, nil
}

func (r *PermissionRepository) DeletePermission(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return errors.New("[PermissionRepository.DeletePermission] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Where("id = ?", id).Delete(&entity.Permission{}).Error; err != nil {
		tx.Rollback()
		r.Log.Error(err)
		return err
	}

	tx.Commit()
	return nil
}

func (r *PermissionRepository) GetAllPermissionsByNames(names []string) (*[]entity.Permission, error) {
	var permissions []entity.Permission
	if err := r.DB.Where("name IN (?)", names).Find(&permissions).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}
	return &permissions, nil
}
