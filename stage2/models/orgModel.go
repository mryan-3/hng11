package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Organisation models
type Organisation struct {
	gorm.Model
	ID            uuid.UUID       `json:"orgId" gorm:"type:uuid;default:gen_random_uuid();primary_key" validate:"unique"`
	Name          string          `json:"name" gorm:"type:varchar(255);not null" validate:"required"`
	Description   string          `json:"description" gorm:"type:varchar(255)"`
	Users []*User `gorm:"many2many:user_organizations;"`
}
