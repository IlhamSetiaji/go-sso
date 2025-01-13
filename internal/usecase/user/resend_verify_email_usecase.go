package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/request"
	"app/go-sso/internal/messaging"
	"app/go-sso/internal/repository"
	"app/go-sso/utils"
	"errors"
	"log"
	"strconv"

	"github.com/sirupsen/logrus"
)

type IResendVerfiyEmailUseCaseRequest struct {
	Email string `json:"email"`
}

type IResendVerfiyEmailUseCaseResponse struct {
	Message string `json:"message"`
}

type IResendVerfiyEmailUseCase interface {
	Execute(payload IResendVerfiyEmailUseCaseRequest) (*IResendVerfiyEmailUseCaseResponse, error)
}

type ResendVerfiyEmailUseCase struct {
	Log         *logrus.Logger
	Repository  repository.IUserRepository
	MailMessage messaging.IMailMessage
}

func NewResendVerfiyEmailUseCase(
	log *logrus.Logger,
	repository repository.IUserRepository,
	mailMessage messaging.IMailMessage,
) IResendVerfiyEmailUseCase {
	return &ResendVerfiyEmailUseCase{
		Log:         log,
		Repository:  repository,
		MailMessage: mailMessage,
	}
}

func (u *ResendVerfiyEmailUseCase) Execute(payload IResendVerfiyEmailUseCaseRequest) (*IResendVerfiyEmailUseCaseResponse, error) {
	user, err := u.Repository.FindByEmailOnly(payload.Email)
	if err != nil {
		u.Log.Error("[UserUseCase.VerifyUserEmail] " + err.Error())
		return nil, err
	}

	if user == nil {
		u.Log.Warn("[UserUseCase.VerifyUserEmail] User not found")
		return nil, errors.New("user not found")
	}

	err = u.Repository.DeleteUserToken(user.Email, entity.UserTokenVerification)
	if err != nil {
		u.Log.Error("[UserUseCase.VerifyUserEmail] " + err.Error())
		return nil, err
	}

	randomIntToken, err := utils.GenerateRandomIntToken(6)
	if err != nil {
		log.Fatalf("Failed to generate random token: %v", err)
	}
	randomIntTokenInt, err := strconv.Atoi(randomIntToken)
	if err != nil {
		u.Log.Error("[UserUseCase.VerifyUserEmail] " + err.Error())
		return nil, err
	}
	if err := u.Repository.CreateUserToken(payload.Email, randomIntTokenInt, entity.UserTokenVerification); err != nil {
		u.Log.Error("[UserUseCase.Register] " + err.Error())
		return nil, err
	}

	if _, err := u.MailMessage.SendMail(&request.MailRequest{
		Email:   user.Email,
		Subject: "Email Verification",
		Body:    "Your verification code is " + randomIntToken,
		From:    "ilham.ahmadz18@gmail.com",
		To:      user.Email,
	}); err != nil {
		u.Log.Error("[UserUseCase.Register] " + err.Error())
		return nil, err
	}

	return &IResendVerfiyEmailUseCaseResponse{
		Message: "Email verification code has been sent",
	}, nil
}

func ResendVerfiyEmailUseCaseFactory(log *logrus.Logger) IResendVerfiyEmailUseCase {
	userRepository := repository.UserRepositoryFactory(log)
	mailMessage := messaging.MailMessageFactory(log)
	return NewResendVerfiyEmailUseCase(log, userRepository, mailMessage)
}
