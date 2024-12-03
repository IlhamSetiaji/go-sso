package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationStructure struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID               `json:"id" gorm:"type:char(36);primaryKey"`
	OrganizationID uuid.UUID               `json:"organization_id" gorm:"type:char(36)"`
	Name           string                  `json:"name"`
	JobLevelID     uuid.UUID               `json:"job_level_id" gorm:"type:char(36)"`
	ParentID       *uuid.UUID              `json:"parent_id" gorm:"type:char(36)"`
	Level          int                     `json:"level" gorm:"index"`    // Add level for hierarchy depth
	Path           string                  `json:"path" gorm:"type:text"` // Store full path for easy traversal
	Organization   Organization            `json:"organization" gorm:"foreignKey:OrganizationID;references:ID;constraint:OnDelete:CASCADE"`
	JobLevel       JobLevel                `json:"job_level" gorm:"foreignKey:JobLevelID;references:ID;constraint:OnDelete:CASCADE"`
	Parent         *OrganizationStructure  `json:"parent" gorm:"foreignKey:ParentID;references:ID;constraint:OnDelete:CASCADE"`
	Children       []OrganizationStructure `json:"children" gorm:"foreignKey:ParentID;references:ID"`
}

func (organizationStructure *OrganizationStructure) BeforeCreate(tx *gorm.DB) (err error) {
	organizationStructure.ID = uuid.New()
	organizationStructure.CreatedAt = time.Now().Add(time.Hour * 7)
	organizationStructure.UpdatedAt = time.Now().Add(time.Hour * 7)

	// Set level and path based on parent
	if organizationStructure.ParentID != nil {
		var parent OrganizationStructure
		if err := tx.First(&parent, "id = ?", organizationStructure.ParentID).Error; err != nil {
			return err
		}
		organizationStructure.Level = parent.Level + 1
		organizationStructure.Path = fmt.Sprintf("%s/%s", parent.Path, organizationStructure.ID.String())
	} else {
		organizationStructure.Level = 0
		organizationStructure.Path = organizationStructure.ID.String()
	}
	return nil
}

func (organizationStructure *OrganizationStructure) BeforeUpdate(tx *gorm.DB) (err error) {
	organizationStructure.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (OrganizationStructure) TableName() string {
	return "organization_structures"
}

// Helper methods for tree operations
func (o *OrganizationStructure) GetAncestors(db *gorm.DB) ([]OrganizationStructure, error) {
	var ancestors []OrganizationStructure
	if o.Path == "" {
		return ancestors, nil
	}

	ancestorIDs := strings.Split(o.Path, "/")
	return ancestors, db.Where("id IN ?", ancestorIDs).Order("level asc").Find(&ancestors).Error
}

func (o *OrganizationStructure) GetDescendants(db *gorm.DB) ([]OrganizationStructure, error) {
	var descendants []OrganizationStructure
	return descendants, db.Where("path LIKE ?", o.Path+"%").
		Where("id != ?", o.ID).
		Order("level asc").
		Find(&descendants).Error
}

// Get siblings (nodes at same level with same parent)
func (o *OrganizationStructure) GetSiblings(db *gorm.DB) ([]OrganizationStructure, error) {
	var siblings []OrganizationStructure
	query := db.Where("parent_id = ? AND id != ?", o.ParentID, o.ID)
	return siblings, query.Find(&siblings).Error
}

// Check if node is ancestor of another node
func (o *OrganizationStructure) IsAncestorOf(child *OrganizationStructure) bool {
	return strings.HasPrefix(child.Path, o.Path+"/")
}

// Check if node is descendant of another node
func (o *OrganizationStructure) IsDescendantOf(parent *OrganizationStructure) bool {
	return strings.HasPrefix(o.Path, parent.Path+"/")
}

// Get the root node of the tree
func GetRoot(db *gorm.DB, organizationID uuid.UUID) (*OrganizationStructure, error) {
	var root OrganizationStructure
	return &root, db.Where("organization_id = ? AND parent_id IS NULL", organizationID).First(&root).Error
}

// Example usage for moving a node in the tree
func (o *OrganizationStructure) MoveTo(db *gorm.DB, newParentID *uuid.UUID) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	oldPath := o.Path

	// Update parent
	o.ParentID = newParentID

	// Recalculate level and path
	if err := o.BeforeCreate(tx); err != nil {
		tx.Rollback()
		return err
	}

	// Update current node
	if err := tx.Save(o).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update all descendants
	if err := tx.Model(&OrganizationStructure{}).
		Where("path LIKE ?", oldPath+"/%").
		Update("path", gorm.Expr("REPLACE(path, ?, ?)", oldPath, o.Path)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
