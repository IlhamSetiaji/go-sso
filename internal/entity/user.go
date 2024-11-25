package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatus string
type UserGender string

const (
	ACTIVE   UserStatus = "ACTIVE"
	INACTIVE UserStatus = "INACTIVE"
	PENDING  UserStatus = "PENDING"
)

const (
	MALE   UserGender = "MALE"
	FEMALE UserGender = "FEMALE"
)

type User struct {
	gorm.Model
	ID              uuid.UUID   `json:"id" gorm:"type:char(36);primaryKey"`
	OauthID         string      `json:"oauth_id" gorm:"unique; default:null"`
	Username        string      `json:"username" gorm:"unique;not null"`
	Email           string      `json:"email" gorm:"unique;not null"`
	Name            string      `json:"name"`
	Password        string      `json:"password"`
	Gender          UserGender  `json:"gender" gorm:"not null"`
	EmailVerifiedAt time.Time   `json:"email_verified_at" gorm:"default:null"`
	Status          UserStatus  `json:"status" gorm:"default:PENDING"`
	Roles           []Role      `json:"roles" gorm:"many2many:user_roles;"` // many to many relationship
	AuthTokens      []AuthToken `json:"auth_tokens" gorm:"foreignKey:UserID;references:ID"`
	// CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	// UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	// DeletedAt       time.Time  `json:"deleted_at" gorm:"index"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	user.CreatedAt = time.Now().Add(time.Hour * 7)
	user.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	user.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (User) TableName() string {
	return "users"
}
