package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
	DB  *gorm.DB
}

type UserRepositoryInterface interface {
	FindByEmail(email string) (*entity.User, error)
	FindAllPaginated(page int, pageSize int) (*[]entity.User, int64, error)
}

// Repository: Repository[entity.User]{DB: db},
// , db *gorm.DB
func UserRepositoryFactory(log *logrus.Logger) UserRepositoryInterface {
	return &UserRepository{
		Log: log,
		DB:  config.NewDatabase(),
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
