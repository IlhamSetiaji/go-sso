package repository

import (
	"app/go-sso/internal/config"
	"app/go-sso/internal/entity"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IEmployeeRepository interface {
	FindAllPaginated(page int, pageSize int, search string, isOnboarding string) (*[]entity.Employee, int64, error)
	FindAllEmployees() (*[]entity.Employee, error)
	FindAllEmployeesNotInUsers() (*[]entity.Employee, error)
	Store(employee *entity.Employee) (*entity.Employee, error)
	Update(employee *entity.Employee) (*entity.Employee, error)
	Delete(id uuid.UUID) error
	StoreEmployeeJob(employeeJob *entity.EmployeeJob) (*entity.EmployeeJob, error)
	UpdateEmployeeJob(employeeJob *entity.EmployeeJob) (*entity.EmployeeJob, error)
	FindById(id uuid.UUID) (*entity.Employee, error)
	FindByMidsuitID(midsuitID string) (*entity.Employee, error)
	CountEmployeeRetiredEndByDateRange(startDate string, endDate string) (int64, error)
	GetOrganizationStructureIdDistinct() ([]uuid.UUID, error)
	CountByOrganizationStructureID(organizationStructureID uuid.UUID) (int, error)
	UpdateEmployeeMidsuitID(id uuid.UUID, midsuitID string) (*entity.Employee, error)
	FindEmployeeRecruitmentManager() (*[]entity.Employee, error)
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

func (r *EmployeeRepository) FindAllPaginated(page int, pageSize int, search string, isOnboarding string) (*[]entity.Employee, int64, error) {
	var employees []entity.Employee
	var total int64

	query := r.DB.Preload("EmployeeJob").Preload("User").Preload("Organization")

	if isOnboarding != "" {
		query = query.Where("is_onboarding = ?", isOnboarding)
	}

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
	err := r.DB.Preload("EmployeeJob.Job").Preload("EmployeeJob.Grade").Preload("User").Preload("Organization.OrganizationType").Find(&employees).Error
	if err != nil {
		return nil, err
	}

	return &employees, nil
}

func (r *EmployeeRepository) FindAllEmployeesNotInUsers() (*[]entity.Employee, error) {
	var users []entity.User
	var employeeIDs []uuid.UUID

	// Fetch users where employee_id is not NULL
	if err := r.DB.Where("employee_id IS NOT NULL").Find(&users).Error; err != nil {
		return nil, err
	}

	// Collect employee IDs while checking for nil values
	for _, user := range users {
		if user.EmployeeID != nil { // Avoid dereferencing nil pointer
			employeeIDs = append(employeeIDs, *user.EmployeeID)
		}
	}

	// Debugging: Check for invalid UUIDs
	for _, id := range employeeIDs {
		if _, err := uuid.Parse(id.String()); err != nil {
			fmt.Printf("Invalid UUID: %s - Error: %s\n", id, err)
		} else {
			fmt.Printf("Valid UUID: %s\n", id)
		}
	}

	var employees []entity.Employee

	if len(employeeIDs) > 0 {
		err := r.DB.Where("id NOT IN (?)", employeeIDs).
			Preload("EmployeeJob.Job").
			Preload("User").
			Preload("Organization.OrganizationType").
			Find(&employees).Error
		if err != nil {
			r.Log.Error("[EmployeeRepository] FindAllEmployeesNotInUsers: ", err)
			return nil, err
		}
	} else {
		err := r.DB.Preload("EmployeeJob.Job").
			Preload("User").
			Preload("Organization.OrganizationType").
			Find(&employees).Error
		if err != nil {
			return nil, err
		}
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

func (r *EmployeeRepository) Update(employee *entity.Employee) (*entity.Employee, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error(tx.Error)
		return nil, tx.Error
	}

	if err := tx.Where("id = ?", employee.ID).Updates(employee).Error; err != nil {
		tx.Rollback()
		r.Log.Error(err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error(err)
		return nil, err
	}

	return employee, nil
}

func (r *EmployeeRepository) Delete(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var employee entity.Employee

	if err := tx.Where("id = ?", id).First(&employee).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", employee.ID).Delete(&employee).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *EmployeeRepository) StoreEmployeeJob(employeeJob *entity.EmployeeJob) (*entity.EmployeeJob, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(employeeJob).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return employeeJob, nil
}

func (r *EmployeeRepository) UpdateEmployeeJob(employeeJob *entity.EmployeeJob) (*entity.EmployeeJob, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Where("id = ?", employeeJob.ID).Updates(employeeJob).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return employeeJob, nil
}

func (r *EmployeeRepository) GetOrganizationStructureIdDistinct() ([]uuid.UUID, error) {
	var organizationIDs []uuid.UUID
	err := r.DB.Raw(`
		SELECT DISTINCT organization_structure_id
		FROM employee_jobs
	`).Scan(&organizationIDs).Error
	if err != nil {
		r.Log.Errorf("[EmployeeRepository.GetOrganizationStructureIdDistinct] error when querying distinct organization structure id: %v", err)
		return nil, err
	}
	return organizationIDs, nil
}

func (r *EmployeeRepository) CountByOrganizationStructureID(organizationStructureID uuid.UUID) (int, error) {
	var count int
	err := r.DB.Raw(`
		SELECT COUNT(*)
		FROM employee_jobs
		WHERE organization_structure_id = ?
	`, organizationStructureID).Scan(&count).Error
	if err != nil {
		r.Log.Errorf("[EmployeeRepository.CountByOrganizationStructureID] error when querying count by job level id: %v", err)
		return 0, err
	}
	return count, nil
}

func (r *EmployeeRepository) FindByMidsuitID(midsuitID string) (*entity.Employee, error) {
	var employee entity.Employee
	err := r.DB.Preload("EmployeeJob.Job").Preload("User").Preload("Organization").Where("midsuit_id = ?", midsuitID).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *EmployeeRepository) UpdateEmployeeMidsuitID(id uuid.UUID, midsuitID string) (*entity.Employee, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var employee entity.Employee

	if err := tx.Where("id = ?", id).First(&employee).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	employee.MidsuitID = midsuitID

	if err := tx.Save(&employee).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &employee, nil
}

func (r *EmployeeRepository) FindEmployeeRecruitmentManager() (*[]entity.Employee, error) {
	var employees []entity.Employee
	err := r.DB.Raw(`
			SELECT 
					e.*
			FROM 
					employees e
			JOIN 
					employee_jobs ej ON e.id = ej.employee_id
			JOIN 
					jobs j ON ej.job_id = j.id
			JOIN 
					organizations o ON e.organization_id = o.id
			WHERE 
					j.name = 'Recruitment Manager' 
					AND o.name = 'Head Office'
	`).Scan(&employees).Error
	if err != nil {
		return nil, err
	}
	return &employees, nil
}

func EmployeeRepositoryFactory(log *logrus.Logger) IEmployeeRepository {
	db := config.NewDatabase()
	return NewEmployeeRepository(log, db)
}
