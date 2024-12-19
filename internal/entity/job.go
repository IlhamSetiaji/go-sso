package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model              `json:"-"`
	ID                      uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey"`
	Name                    string     `json:"name"`
	OrganizationStructureID uuid.UUID  `json:"organization_structure_id" gorm:"type:char(36)"`
	ParentID                *uuid.UUID `json:"parent_id" gorm:"type:char(36)"`
	Level                   int        `json:"level" gorm:"index"`    // Add level for hierarchy depth
	Path                    string     `json:"path" gorm:"type:text"` // Store full path for easy traversal
	Existing                int        `json:"existing" gorm:"default:0"`
	MidsuitID               string     `json:"midsuit_id"`

	// JobLevelData          JobLevel              `json:"job_level_data" gorm:"-"`
	OrganizationStructure OrganizationStructure `json:"organization_structure" gorm:"foreignKey:OrganizationStructureID;references:ID;constraint:OnDelete:CASCADE"`
	Parent                *Job                  `json:"parent" gorm:"foreignKey:ParentID;references:ID;constraint:OnDelete:CASCADE"`
	Children              []Job                 `json:"children" gorm:"foreignKey:ParentID;references:ID"`
	EmployeeJobs          []EmployeeJob         `json:"employee_jobs" gorm:"foreignKey:JobID;references:ID"`
}

func (job *Job) BeforeCreate(tx *gorm.DB) (err error) {
	job.ID = uuid.New()
	job.CreatedAt = time.Now().Add(time.Hour * 7)
	job.UpdatedAt = time.Now().Add(time.Hour * 7)

	// Set level and path based on parent
	if job.ParentID != nil {
		var parent Job
		if err := tx.First(&parent, "id = ?", job.ParentID).Error; err != nil {
			return err
		}
		job.Level = parent.Level + 1
		job.Path = fmt.Sprintf("%s/%s", parent.Path, job.ID.String())
	} else {
		job.Level = 0
		job.Path = job.ID.String()
	}
	return nil
}

func (job *Job) BeforeUpdate(tx *gorm.DB) (err error) {
	job.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (Job) TableName() string {
	return "jobs"
}

// Helper methods for tree operations
func (j *Job) GetAncestorsJob(db *gorm.DB) ([]Job, error) {
	var ancestors []Job
	if j.Path == "" {
		return ancestors, nil
	}

	ancestorIDs := strings.Split(j.Path, "/")
	return ancestors, db.Where("id IN ?", ancestorIDs).Order("level asc").Find(&ancestors).Error
}

func (j *Job) GetDescendantsJob(db *gorm.DB) ([]Job, error) {
	var descendants []Job
	return descendants, db.Where("path LIKE ?", j.Path+"%").
		Where("id != ?", j.ID).
		Order("level asc").
		Find(&descendants).Error
}

// Get siblings (nodes at same level with same parent)
func (j *Job) GetSiblingsJob(db *gorm.DB) ([]Job, error) {
	var siblings []Job
	query := db.Where("parent_id = ? AND id != ?", j.ParentID, j.ID)
	return siblings, query.Find(&siblings).Error
}

// Check if node is ancestor of another node
func (j *Job) IsAncestorOfJob(child *Job) bool {
	return strings.HasPrefix(child.Path, j.Path+"/")
}

// Check if node is descendant of another node
func (j *Job) IsDescendantOfJob(parent *Job) bool {
	return strings.HasPrefix(j.Path, parent.Path+"/")
}

// Get the root node of the tree
func GetRootJob(db *gorm.DB, jobID uuid.UUID) (*Job, error) {
	var root Job
	return &root, db.Where("job_id = ? AND parent_id IS NULL", jobID).First(&root).Error
}

// Example usage for moving a node in the tree
func (j *Job) MoveToJob(db *gorm.DB, newParentID *uuid.UUID) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	oldPath := j.Path

	// Update parent
	j.ParentID = newParentID

	// Recalculate level and path
	if err := j.BeforeCreate(tx); err != nil {
		tx.Rollback()
		return err
	}

	// Update current node
	if err := tx.Save(j).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update all descendants
	if err := tx.Model(&Job{}).
		Where("path LIKE ?", oldPath+"/%").
		Update("path", gorm.Expr("REPLACE(path, ?, ?)", oldPath, j.Path)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
