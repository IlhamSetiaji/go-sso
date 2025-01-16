package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IUserRepository interface {
	FindByEmail(email string) (*entity.User, error)
	FindAllPaginated(page int, pageSize int, search string) (*[]entity.User, int64, error)
	FindById(id uuid.UUID) (*entity.User, error)
	FindByIdOnly(id uuid.UUID) (*entity.User, error)
	FindByEmailOnly(email string) (*entity.User, error)
	GetAllUsers() (*[]entity.User, error)
	GetAllUsersDoesNotHaveEmployee() (*[]entity.User, error)
	CreateUser(user *entity.User, roleIDs []uuid.UUID) (*entity.User, error)
	UpdateUser(user *entity.User, roleIDs []uuid.UUID) (*entity.User, error)
	UpdateUserOnly(user *entity.User) (*entity.User, error)
	UpdateEmployeeIdToNull(user *entity.User) (*entity.User, error)
	DeleteUser(id uuid.UUID) error
	CreateUserToken(email string, token int, tokenType entity.UserTokenType) error
	AcknowledgeUserToken(email string, token int) error
	VerifyUserEmail(email string) error
	FindUserTokenByEmail(email string) (*entity.UserToken, error)
	DeleteUserToken(email string, tokenType entity.UserTokenType) error
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

func (r *UserRepository) FindAllPaginated(page int, pageSize int, search string) (*[]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := r.DB.Preload("Employee.Organization").Preload("Employee.EmployeeJob.Job").Preload("Employee.EmployeeJob.EmpOrganization").Preload("Employee.EmployeeJob.OrganizationLocation")

	if search != "" {
		query = query.Where("email ILIKE ?", "%"+search+"%").Or("name ILIKE ?", "%"+search+"%").Or("username ILIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		r.Log.Error("[UserRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[UserRepository.FindAllPaginated] " + err.Error())
	}

	if err := query.Model(&entity.User{}).Count(&total).Error; err != nil {
		r.Log.Error("[UserRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[UserRepository.FindAllPaginated] " + err.Error())
	}

	return &users, total, nil
}

func (r *UserRepository) GetAllUsers() (*[]entity.User, error) {
	var users []entity.User

	if err := r.DB.Preload("Roles.Application").Preload("Employee.Organization").Find(&users).Error; err != nil {
		r.Log.Error("[UserRepository.GetAllUsers] " + err.Error())
		return nil, errors.New("[UserRepository.GetAllUsers] " + err.Error())
	}

	return &users, nil
}

func (r *UserRepository) GetAllUsersDoesNotHaveEmployee() (*[]entity.User, error) {
	var users []entity.User

	if err := r.DB.Preload("Roles.Application").Where("employee_id IS NULL").Find(&users).Error; err != nil {
		r.Log.Error("[UserRepository.GetAllUsersDoesNotHaveEmployee] " + err.Error())
		return nil, errors.New("[UserRepository.GetAllUsersDoesNotHaveEmployee] " + err.Error())
	}

	return &users, nil
}

func (r *UserRepository) FindById(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.DB.Preload("Roles.Permissions").Preload("Employee.Organization.OrganizationStructures.JobLevel").Preload("Employee.Organization.OrganizationType").Preload("Employee.EmployeeJob.Job.OrganizationStructure").Preload("Employee.EmployeeJob.EmpOrganization").Preload("Employee.EmployeeJob.OrganizationLocation").Where("id = ?", id).First(&user).Error
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

func (r *UserRepository) FindByIdOnly(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.DB.Preload("Roles").Where("id = ?", id).First(&user).Error
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

func (r *UserRepository) FindByEmailOnly(email string) (*entity.User, error) {
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

func (r *UserRepository) CreateUser(user *entity.User, roleIDs []uuid.UUID) (*entity.User, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("[UserRepository.CreateUser] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.CreateUser] " + err.Error())
		return nil, errors.New("[UserRepository.CreateUser] " + err.Error())
	}

	for _, roleID := range roleIDs {
		var role entity.Role
		if err := tx.First(&role, "id = ?", roleID).Error; err != nil {
			tx.Rollback()
			r.Log.Error("[UserRepository.CreateUser] Role not found: " + err.Error())
			return nil, errors.New("[UserRepository.CreateUser] Role not found: " + err.Error())
		}

		var userRole = entity.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		}

		if err := tx.Create(&userRole).Error; err != nil {
			tx.Rollback()
			r.Log.Error("[UserRepository.AppendUser] " + err.Error())
			return nil, errors.New("[UserRepository.AppendUser] " + err.Error())
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.CreateUser] failed to commit transaction: " + err.Error())
		return nil, errors.New("[UserRepository.CreateUser] failed to commit transaction: " + err.Error())
	}

	if err := r.DB.Preload("Roles").First(user, user.ID).Error; err != nil {
		r.Log.Error("[UserRepository.CreateUser] Failed to reload user: " + err.Error())
		return nil, errors.New("[UserRepository.CreateUser] Failed to reload user: " + err.Error())
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(user *entity.User, roleIDs []uuid.UUID) (*entity.User, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("[UserRepository.UpdateUser] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Model(&user).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.UpdateUser] " + err.Error())
		return nil, errors.New("[UserRepository.UpdateUser] " + err.Error())
	}

	if len(roleIDs) > 0 {
		// delete user role by user id
		if err := tx.Where("user_id = ?", user.ID).Delete(&entity.UserRole{}).Error; err != nil {
			tx.Rollback()
			r.Log.Error("[UserRepository.UpdateUser] " + err.Error())
			return nil, errors.New("[UserRepository.UpdateUser] " + err.Error())
		}
		for _, roleID := range roleIDs {
			var role entity.Role
			if err := tx.First(&role, "id = ?", roleID).Error; err != nil {
				tx.Rollback()
				r.Log.Error("[UserRepository.UpdateUser] Role not found: " + err.Error())
				return nil, errors.New("[UserRepository.UpdateUser] Role not found: " + err.Error())
			}

			var userRole = entity.UserRole{
				UserID: user.ID,
				RoleID: role.ID,
			}

			if err := tx.Create(&userRole).Error; err != nil {
				tx.Rollback()
				r.Log.Error("[UserRepository.UpdateUser] " + err.Error())
				return nil, errors.New("[UserRepository.UpdateUser] " + err.Error())
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.CreateUser] failed to commit transaction: " + err.Error())
		return nil, errors.New("[UserRepository.CreateUser] failed to commit transaction: " + err.Error())
	}

	if err := r.DB.Preload("Roles").First(user, user.ID).Error; err != nil {
		r.Log.Error("[UserRepository.CreateUser] Failed to reload user: " + err.Error())
		return nil, errors.New("[UserRepository.CreateUser] Failed to reload user: " + err.Error())
	}

	return user, nil
}

func (r *UserRepository) UpdateUserOnly(user *entity.User) (*entity.User, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("[UserRepository.UpdateUserOnly] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Model(&user).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.UpdateUserOnly] " + err.Error())
		return nil, errors.New("[UserRepository.UpdateUserOnly] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.UpdateUserOnly] failed to commit transaction: " + err.Error())
		return nil, errors.New("[UserRepository.UpdateUserOnly] failed to commit transaction: " + err.Error())
	}

	if err := r.DB.Preload("Roles").First(user, user.ID).Error; err != nil {
		r.Log.Error("[UserRepository.UpdateUserOnly] Failed to reload user: " + err.Error())
		return nil, errors.New("[UserRepository.UpdateUserOnly] Failed to reload user: " + err.Error())
	}

	return user, nil
}

func (r *UserRepository) UpdateEmployeeIdToNull(user *entity.User) (*entity.User, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("[UserRepository.UpdateEmployeeIdToNull] failed to begin transaction: " + tx.Error.Error())
	}

	user.EmployeeID = nil

	if err := tx.Model(&user).Where("id = ?", user.ID).Update("employee_id", nil).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.UpdateEmployeeIdToNull] " + err.Error())
		return nil, errors.New("[UserRepository.UpdateEmployeeIdToNull] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.UpdateEmployeeIdToNull] failed to commit transaction: " + err.Error())
		return nil, errors.New("[UserRepository.UpdateEmployeeIdToNull] failed to commit transaction: " + err.Error())
	}

	if err := r.DB.Preload("Roles").First(user, user.ID).Error; err != nil {
		r.Log.Error("[UserRepository.UpdateEmployeeIdToNull] Failed to reload user: " + err.Error())
		return nil, errors.New("[UserRepository.UpdateEmployeeIdToNull] Failed to reload user: " + err.Error())
	}

	return user, nil
}

func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return errors.New("[UserRepository.DeleteUser] failed to begin transaction: " + tx.Error.Error())
	}

	var user entity.User
	if err := tx.First(&user, "id = ?", id).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.DeleteUser] User not found: " + err.Error())
		return errors.New("[UserRepository.DeleteUser] User not found: " + err.Error())
	}

	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.DeleteUser] " + err.Error())
		return errors.New("[UserRepository.DeleteUser] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.DeleteUser] failed to commit transaction: " + err.Error())
		return errors.New("[UserRepository.DeleteUser] failed to commit transaction: " + err.Error())
	}

	return nil
}

func (r *UserRepository) CreateUserToken(email string, token int, tokenType entity.UserTokenType) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return errors.New("[UserRepository.CreateUserToken] failed to begin transaction: " + tx.Error.Error())
	}

	var user entity.User
	if err := tx.First(&user, "email = ?", email).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.CreateUserToken] User not found: " + err.Error())
		return errors.New("[UserRepository.CreateUserToken] User not found: " + err.Error())
	}

	var userToken = entity.UserToken{
		Email:     email,
		Token:     token,
		TokenType: tokenType,
	}

	if err := tx.Create(&userToken).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.CreateUserToken] " + err.Error())
		return errors.New("[UserRepository.CreateUserToken] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.CreateUserToken] failed to commit transaction: " + err.Error())
		return errors.New("[UserRepository.CreateUserToken] failed to commit transaction: " + err.Error())
	}

	return nil
}

func (r *UserRepository) AcknowledgeUserToken(email string, token int) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return errors.New("[UserRepository.AcknowledgeUserToken] failed to begin transaction: " + tx.Error.Error())
	}

	var userToken entity.UserToken
	if err := tx.First(&userToken, "email = ? AND token = ?", email, token).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.AcknowledgeUserToken] User token not found: " + err.Error())
		return errors.New("[UserRepository.AcknowledgeUserToken] User token not found: " + err.Error())
	}

	if err := tx.Where("email = ? AND token = ?", userToken.Email, userToken.Token).Delete(&userToken).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.AcknowledgeUserToken] " + err.Error())
		return errors.New("[UserRepository.AcknowledgeUserToken] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.AcknowledgeUserToken] failed to commit transaction: " + err.Error())
		return errors.New("[UserRepository.AcknowledgeUserToken] failed to commit transaction: " + err.Error())
	}

	return nil
}

func (r *UserRepository) VerifyUserEmail(email string) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return errors.New("[UserRepository.VerifyUserEmail] failed to begin transaction: " + tx.Error.Error())
	}

	var user entity.User
	if err := tx.First(&user, "email = ?", email).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.VerifyUserEmail] User not found: " + err.Error())
		return errors.New("[UserRepository.VerifyUserEmail] User not found: " + err.Error())
	}

	if err := tx.Model(&user).Where("email = ?", email).Updates(map[string]interface{}{
		"email_verified_at": time.Now(),
		"status":            entity.USER_ACTIVE,
	}).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.VerifyUserEmail] " + err.Error())
		return errors.New("[UserRepository.VerifyUserEmail] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.VerifyUserEmail] failed to commit transaction: " + err.Error())
		return errors.New("[UserRepository.VerifyUserEmail] failed to commit transaction: " + err.Error())
	}

	return nil
}

func (r *UserRepository) FindUserTokenByEmail(email string) (*entity.UserToken, error) {
	var userToken entity.UserToken
	err := r.DB.Where("token_type = ?", entity.UserTokenVerification).First(&userToken, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warn("[UserRepository.FindUserTokenByEmail] User token not found")
			return nil, nil
		} else {
			r.Log.Error("[UserRepository.FindUserTokenByEmail] " + err.Error())
			return nil, errors.New("[UserRepository.FindUserTokenByEmail] " + err.Error())
		}
	}
	return &userToken, nil
}

func (r *UserRepository) DeleteUserToken(email string, tokenType entity.UserTokenType) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return errors.New("[UserRepository.DeleteUserToken] failed to begin transaction: " + tx.Error.Error())
	}

	var userToken entity.UserToken
	if err := tx.First(&userToken, "email = ? AND token_type = ?", email, tokenType).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.DeleteUserToken] User token not found: " + err.Error())
		return errors.New("[UserRepository.DeleteUserToken] User token not found: " + err.Error())
	}

	if err := tx.Where("email = ? AND token_type = ?", userToken.Email, userToken.TokenType).Delete(&userToken).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.DeleteUserToken] " + err.Error())
		return errors.New("[UserRepository.DeleteUserToken] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.DeleteUserToken] failed to commit transaction: " + err.Error())
		return errors.New("[UserRepository.DeleteUserToken] failed to commit transaction: " + err.Error())
	}

	return nil
}

func UserRepositoryFactory(log *logrus.Logger) IUserRepository {
	db := config.NewDatabase()
	return NewUserRepository(log, db)
}
