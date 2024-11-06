package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"
	"errors"
	"log"

	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *log.Logger
	DB  *gorm.DB
}

type UserRepositoryInterface interface {
	FindByEmail(email string) (*entity.User, error)
}

// Repository: Repository[entity.User]{DB: db},
// , db *gorm.DB
func UserRepositoryFactory(log *log.Logger) UserRepositoryInterface {
	return &UserRepository{
		Log: log,
		DB:  config.NewDatabase(log),
	}
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Preload("Roles.Permissions").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, errors.New("[UserRepository.FindByEmail] " + err.Error())
		}
	}
	return &user, nil
}
