package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User models
type User struct {
	gorm.Model

	UserID        uuid.UUID       `json:"userId" gorm:"type:uuid;default:gen_random_uuid();primary_key" validate:"unique"`
	FirstName     string          `json:"firstName" gorm:"type:varchar(255);not null" validate:"required"`
	LastName      string          `json:"lastName" gorm:"type:varchar(255);not null" validate:"required"`
	Email         string          `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password      string          `json:"-" gorm:"not null" validate:"required"` // "-" exclude field from json response
	Phone         string          `json:"phone" gorm:"type:varchar(255)"`
	Organisations []*Organisation `gorm:"many2many:user_organizations;"`
}
