package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IEmployeeRepository interface {
	FindAllPaginated(page int, pageSize int, search string) (*[]entity.Employee, int64, error)
	FindAllEmployees() (*[]entity.Employee, error)
	Store(employee *entity.Employee) (*entity.Employee, error)
	FindById(id uuid.UUID) (*entity.Employee, error)
	CountEmployeeRetiredEndByDateRange(startDate string, endDate string) (int64, error)
}

type EmployeeRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewEmployeeRepository(log *logrus.Logger, db *gorm.DB) IEmployeeRepository {
	return &EmployeeRepository{
		Log: log,
		DB:  db,
	}
}

func (r *EmployeeRepository) FindAllPaginated(page int, pageSize int, search string) (*[]entity.Employee, int64, error) {
	var employees []entity.Employee
	var total int64

	query := r.DB.Preload("EmployeeJob").Preload("User").Preload("Organization")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Model(&entity.Employee{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return &employees, total, nil
}

func (r *EmployeeRepository) FindAllEmployees() (*[]entity.Employee, error) {
	var employees []entity.Employee
	err := r.DB.Preload("EmployeeJob.Job").Preload("User").Preload("Organization.OrganizationType").Find(&employees).Error
	if err != nil {
		return nil, err
	}

	return &employees, nil
}

func (r *EmployeeRepository) FindById(id uuid.UUID) (*entity.Employee, error) {
	var employee entity.Employee
	err := r.DB.Preload("EmployeeJob.Job").Preload("User").Preload("Organization.OrganizationType").Where("id = ?", id).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *EmployeeRepository) CountEmployeeRetiredEndByDateRange(startDate string, endDate string) (int64, error) {
	var total int64
	err := r.DB.Model(&entity.Employee{}).
		Where("retirement_date BETWEEN ? AND ?", startDate, endDate).Or("end_date BETWEEN ? AND ?", startDate, endDate).
		Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *EmployeeRepository) Store(employee *entity.Employee) (*entity.Employee, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(employee).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return employee, nil
}

func EmployeeRepositoryFactory(log *logrus.Logger) IEmployeeRepository {
	db := config.NewDatabase()
	return NewEmployeeRepository(log, db)
}
