package entities

import "gorm.io/gorm"

type Role string

const (
	SuperAdmin Role = "super_admin"
	Admin      Role = "admin"
	User       Role = "user"
)

type Users struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     Role   `gorm:"type:role;not null"`
	Active   bool   `gorm:"default:true"`
}
