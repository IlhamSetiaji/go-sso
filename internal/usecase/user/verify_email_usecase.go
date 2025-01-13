package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/messaging"
	"app/go-sso/internal/repository"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

type IVerifyEmailUseCaseRequest struct {
	Email string `json:"email"`
	Token int    `json:"token"`
}

type IVerifyEmailUseCaseResponse struct {
	User *entity.User `json:"user"`
}

type IVerifyEmailUseCase interface {
	Execute(payload IVerifyEmailUseCaseRequest) (*IVerifyEmailUseCaseResponse, error)
}

type VerifyEmailUseCase struct {
	Log         *logrus.Logger
	Repository  repository.IUserRepository
	MailMessage messaging.IMailMessage
}

func NewVerifyEmailUseCase(
	log *logrus.Logger,
	repository repository.IUserRepository,
	mailMessage messaging.IMailMessage,
) IVerifyEmailUseCase {
	return &VerifyEmailUseCase{
		Log:         log,
		Repository:  repository,
		MailMessage: mailMessage,
	}
}

func (u *VerifyEmailUseCase) Execute(payload IVerifyEmailUseCaseRequest) (*IVerifyEmailUseCaseResponse, error) {
	user, err := u.Repository.FindByEmailOnly(payload.Email)
	if err != nil {
		u.Log.Error("[UserUseCase.VerifyUserEmail] " + err.Error())
		return nil, err
	}

	if user == nil {
		u.Log.Warn("[UserUseCase.VerifyUserEmail] User not found")
		return nil, errors.New("user not found")
	}

	userToken, err := u.Repository.FindUserTokenByEmail(payload.Email)
	if err != nil {
		u.Log.Error("[UserUseCase.VerifyUserEmail] " + err.Error())
		return nil, err
	}

	if userToken == nil {
		u.Log.Warn("[UserUseCase.VerifyUserEmail] User token not found")
		return nil, errors.New("user token not found")
	}

	if userToken.ExpiredAt.Before(time.Now()) {
		u.Log.Warn("[UserUseCase.VerifyUserEmail] Token expired")
		return nil, errors.New("token expired")
	}

	if userToken.Token != payload.Token {
		u.Log.Warn("[UserUseCase.VerifyUserEmail] Invalid token")
		return nil, errors.New("invalid token")
	}

	if err := u.Repository.VerifyUserEmail(user.Email); err != nil {
		u.Log.Error("[UserUseCase.VerifyUserEmail] " + err.Error())
		return nil, err
	}

	if err := u.Repository.AcknowledgeUserToken(user.Email, userToken.Token); err != nil {
		u.Log.Error("[UserUseCase.VerifyUserEmail] " + err.Error())
		return nil, err
	}

	return &IVerifyEmailUseCaseResponse{
		User: user,
	}, nil
}

func VerifyEmailUseCaseFactory(log *logrus.Logger) IVerifyEmailUseCase {
	repository := repository.UserRepositoryFactory(log)
	mailMessage := messaging.MailMessageFactory(log)
	return NewVerifyEmailUseCase(log, repository, mailMessage)
}
