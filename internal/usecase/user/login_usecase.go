package usecase

import (
	"app/go-sso/internal/entity"
	request "app/go-sso/internal/http/request/user"
	"app/go-sso/internal/repository"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type LoginUseCase struct {
	Log            *log.Logger
	UserRepository repository.UserRepositoryInterface
}

type LoginUseCaseResponse struct {
	User entity.User `json:"user"`
}

type LoginUseCaseInterface interface {
	Login(request request.LoginRequest) (*LoginUseCaseResponse, error)
}

func LoginUseCaseFactory(
	log *log.Logger,
) LoginUseCaseInterface {
	return &LoginUseCase{
		Log:            log,
		UserRepository: repository.UserRepositoryFactory(log),
	}
}

func (uc *LoginUseCase) Login(request request.LoginRequest) (*LoginUseCaseResponse, error) {
	user, err := uc.UserRepository.FindByEmail(request.Email)
	if err != nil {
		return nil, errors.New("[LoginUseCase.Login] " + err.Error())
	}

	if user == nil {
		// uc.Log.Panicf("User not found")
		return nil, errors.New("User not found")
	}

	checkedPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if checkedPassword != nil {
		// uc.Log.Fatalf("Password not match")
		// return nil, errors.New("Password not match:" + checkedPassword.Error())
		return nil, errors.New("Password not match")
	}

	return &LoginUseCaseResponse{
		User: *user,
	}, nil
}
