package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IUserRepository interface {
	FindByEmail(email string) (*entity.User, error)
	FindAllPaginated(page int, pageSize int) (*[]entity.User, int64, error)
	FindById(id string) (*entity.User, error)
}

type UserRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewUserRepository(log *logrus.Logger, db *gorm.DB) IUserRepository {
	return &UserRepository{
		Log: log,
		DB:  db,
	}
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Preload("Roles.Permissions").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warn("[UserRepository.FindByEmail] User not found")
			return nil, nil
		} else {
			r.Log.Error("[UserRepository.FindByEmail] " + err.Error())
			return nil, errors.New("[UserRepository.FindByEmail] " + err.Error())
		}
	}
	return &user, nil
}

func (r *UserRepository) FindAllPaginated(page int, pageSize int) (*[]entity.User, int64, error) {
	var users []entity.User
	var totalCount int64
	offset := (page - 1) * pageSize

	if err := r.DB.Model(&entity.User{}).Count(&totalCount).Error; err != nil {
		r.Log.Error("[UserRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[UserRepository.FindAllPaginated] " + err.Error())
	}

	if err := r.DB.Preload("Roles.Permissions").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		r.Log.Error("[UserRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[UserRepository.FindAllPaginated] " + err.Error())
	}

	return &users, totalCount, nil
}

func (r *UserRepository) FindById(id string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Preload("Roles.Permissions").Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warn("[UserRepository.FindById] User not found")
			return nil, nil
		} else {
			r.Log.Error("[UserRepository.FindById] " + err.Error())
			return nil, errors.New("[UserRepository.FindById] " + err.Error())
		}
	}
	return &user, nil
}

func UserRepositoryFactory(log *logrus.Logger) IUserRepository {
	db := config.NewDatabase()
	return NewUserRepository(log, db)
}
