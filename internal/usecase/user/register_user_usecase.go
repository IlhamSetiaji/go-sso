package usecase

import (
	"app/go-sso/internal/entity"
	"app/go-sso/internal/http/dto"
	"app/go-sso/internal/http/request"
	"app/go-sso/internal/http/response"
	"app/go-sso/internal/messaging"
	"app/go-sso/internal/repository"
	"app/go-sso/utils"
	"errors"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type IRegisterUserUseCaseRequest struct {
	Username    string            `json:"username"`
	Email       string            `json:"email"`
	Name        string            `json:"name"`
	Password    string            `json:"password"`
	Gender      entity.UserGender `form:"gender"`
	MobilePhone string            `form:"mobile_phone"`
	BirthDate   string            `form:"birth_date"`
	BirthPlace  string            `form:"birth_place"`
	Address     string            `form:"address"`
	NoKTP       string            `form:"no_ktp"`
	KTP         string            `form:"ktp"`
}

type IRegisterUserUseCaseResponse struct {
	User *response.UserResponse `json:"user"`
}

type IRegisterUserUseCase interface {
	Execute(payload IRegisterUserUseCaseRequest) (*IRegisterUserUseCaseResponse, error)
}

type RegisterUserUseCase struct {
	Log            *logrus.Logger
	Repository     repository.IUserRepository
	RoleRepository repository.IRoleRepository
	MailMessage    messaging.IMailMessage
}

func NewRegisterUserUseCase(
	log *logrus.Logger,
	repository repository.IUserRepository,
	roleRepository repository.IRoleRepository,
	mailMessage messaging.IMailMessage,
) IRegisterUserUseCase {
	return &RegisterUserUseCase{
		Log:            log,
		Repository:     repository,
		RoleRepository: roleRepository,
		MailMessage:    mailMessage,
	}
}

func (uc *RegisterUserUseCase) Execute(payload IRegisterUserUseCaseRequest) (*IRegisterUserUseCaseResponse, error) {
	user, err := uc.Repository.FindByEmail(payload.Email)
	if err != nil {
		uc.Log.Error("[UserUseCase.Register] " + err.Error())
		return nil, err
	}

	if user != nil {
		uc.Log.Warn("[UserUseCase.Register] User already registered")
		return nil, errors.New("user already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.Log.Error("[UserUseCase.Register] " + err.Error())
		return nil, err
	}

	user = &entity.User{
		Username: payload.Username,
		Email:    payload.Email,
		Name:     payload.Name,
		Password: string(hashedPassword),
		Gender:   payload.Gender,
		Status:   entity.USER_PENDING,
	}

	role, err := uc.RoleRepository.FindByName("Applicant")
	if err != nil {
		uc.Log.Error("[UserUseCase.Register] " + err.Error())
		return nil, err
	}

	if role == nil {
		uc.Log.Error("[UserUseCase.Register] Role not found")
		return nil, errors.New("role not found")
	}

	roleIDs := []uuid.UUID{role.ID}

	if _, err := uc.Repository.CreateUser(user, roleIDs); err != nil {
		uc.Log.Error("[UserUseCase.Register] " + err.Error())
		return nil, err
	}

	randomIntToken, err := utils.GenerateRandomIntToken(6)
	if err != nil {
		log.Fatalf("Failed to generate random token: %v", err)
	}
	if err := uc.Repository.CreateUserToken(payload.Email, int(randomIntToken), entity.UserTokenVerification); err != nil {
		uc.Log.Error("[UserUseCase.Register] " + err.Error())
		return nil, err
	}

	if _, err := uc.MailMessage.SendMail(&request.MailRequest{
		Email:   user.Email,
		Subject: "Email Verification",
		Body:    "Your verification code is " + strconv.Itoa(int(randomIntToken)),
		From:    "ilham.ahmadz18@gmail.com",
		To:      user.Email,
	}); err != nil {
		uc.Log.Error("[UserUseCase.Register] " + err.Error())
		return nil, err
	}

	return &IRegisterUserUseCaseResponse{
		User: dto.ConvertToSingleUserResponse(user),
	}, nil
}

func RegisterUserUseCaseFactory(log *logrus.Logger) IRegisterUserUseCase {
	userRepository := repository.UserRepositoryFactory(log)
	roleRepository := repository.RoleRepositoryFactory(log)
	mailMessage := messaging.MailMessageFactory(log)
	return NewRegisterUserUseCase(log, userRepository, roleRepository, mailMessage)
}
