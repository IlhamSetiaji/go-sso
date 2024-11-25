package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IAuthTokenRepository interface {
	StoreAuthToken(user *entity.User, token *entity.AuthToken) error
	FindAuthToken(userID uuid.UUID, token string) (*entity.AuthToken, error)
	DeleteAuthToken(userID string, token string) error
}

type AuthTokenRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewAuthTokenRepository(log *logrus.Logger, db *gorm.DB) IAuthTokenRepository {
	return &AuthTokenRepository{
		Log: log,
		DB:  db,
	}
}

func (r *AuthTokenRepository) StoreAuthToken(user *entity.User, token *entity.AuthToken) error {
	if err := r.DB.Model(user).Association("AuthTokens").Append(token); err != nil {
		r.Log.Error("[AuthTokenRepository.StoreAuthToken] " + err.Error())
		return err
	}
	return nil
}

func (r *AuthTokenRepository) FindAuthToken(userID uuid.UUID, token string) (*entity.AuthToken, error) {
	var authToken entity.AuthToken
	err := r.DB.Preload("User").Where("user_id = ? AND token = ?", userID, token).First(&authToken).Error
	if err != nil {
		r.Log.Error("[AuthTokenRepository.FindAuthToken] " + err.Error())
		return nil, err
	}
	return &authToken, nil
}

func (r *AuthTokenRepository) DeleteAuthToken(userID string, token string) error {
	err := r.DB.Where("user_id = ? AND token = ?", userID, token).Delete(&entity.AuthToken{}).Error
	if err != nil {
		r.Log.Error("[AuthTokenRepository.DeleteAuthToken] " + err.Error())
		return err
	}
	return nil
}

func AuthTokenRepositoryFactory(log *logrus.Logger) IAuthTokenRepository {
	db := config.NewDatabase()
	return NewAuthTokenRepository(log, db)
}
