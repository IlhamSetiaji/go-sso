package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/repository"
	"errors"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type ILoginUseCaseRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ILoginUseCaseResponse struct {
	User entity.User `json:"user"`
}

type ILoginUseCase interface {
	Execute(request ILoginUseCaseRequest) (*ILoginUseCaseResponse, error)
}

type LoginUseCase struct {
	Log            *logrus.Logger
	UserRepository repository.UserRepositoryInterface
}

func NewLoginUseCase(log *logrus.Logger, userRepository repository.UserRepositoryInterface) ILoginUseCase {
	return &LoginUseCase{
		Log:            log,
		UserRepository: userRepository,
	}
}

func (uc *LoginUseCase) Execute(request ILoginUseCaseRequest) (*ILoginUseCaseResponse, error) {
	user, err := uc.UserRepository.FindByEmail(request.Email)
	if err != nil {
		return nil, errors.New("[LoginUseCase.Execute] " + err.Error())
	}

	if user == nil {
		uc.Log.Error("User not found")
		return nil, errors.New("Email or password is incorrect")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		uc.Log.Error("Password not match")
		return nil, errors.New("Email or password is incorrect")
	}

	return &ILoginUseCaseResponse{
		User: *user,
	}, nil
}

func LoginUseCaseFactory(log *logrus.Logger) ILoginUseCase {
	userRepository := repository.UserRepositoryFactory(log)
	return NewLoginUseCase(log, userRepository)
}
