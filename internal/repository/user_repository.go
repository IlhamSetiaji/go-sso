package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	FindByEmail(email string) (*entity.User, error)
	FindAllPaginated(page int, pageSize int) (*[]entity.User, int64, error)
	StoreAuthToken(user *entity.User, token *entity.AuthToken) error
	FindAuthToken(userID string, token string) (*entity.AuthToken, error)
	DeleteAuthToken(userID string, token string) error
}

type UserRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewUserRepository(log *logrus.Logger, db *gorm.DB) UserRepositoryInterface {
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

func (r *UserRepository) StoreAuthToken(user *entity.User, token *entity.AuthToken) error {
	if err := r.DB.Model(user).Association("AuthTokens").Append(token); err != nil {
		r.Log.Error("[UserRepository.StoreAuthToken] " + err.Error())
		return errors.New("[UserRepository.StoreAuthToken] " + err.Error())
	}
	return nil
}

func (r *UserRepository) FindAuthToken(userID string, token string) (*entity.AuthToken, error) {
	var authToken entity.AuthToken
	err := r.DB.Where("user_id = ? AND token = ?", userID, token).First(&authToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warn("[UserRepository.FindAuthToken] Auth token not found")
			return nil, nil
		} else {
			r.Log.Error("[UserRepository.FindAuthToken] " + err.Error())
			return nil, errors.New("[UserRepository.FindAuthToken] " + err.Error())
		}
	}
	return &authToken, nil
}

func (r *UserRepository) DeleteAuthToken(userID string, token string) error {
	authToken, err := r.FindAuthToken(userID, token)
	if err != nil {
		return err
	}
	if authToken == nil {
		return errors.New("[UserRepository.DeleteAuthToken] Auth token not found")
	}
	if err := r.DB.Delete(authToken).Error; err != nil {
		r.Log.Error("[UserRepository.DeleteAuthToken] " + err.Error())
		return errors.New("[UserRepository.DeleteAuthToken] " + err.Error())
	}
	return nil
}

func UserRepositoryFactory(log *logrus.Logger) UserRepositoryInterface {
	db := config.NewDatabase()
	return NewUserRepository(log, db)
}
