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
	Name     string `gorm:"not null" validate:"required,min=3,max=50"`
	Email    string `gorm:"unique;not null" validate:"required,email"`
	Password string `gorm:"not null" validate:"required,min=8"`
	Role     Role   `gorm:"type:role;not null" validate:"required,role"`
	Active   bool   `gorm:"default:true"`
}

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,strong_password"`
	Role     Role   `json:"role" validate:"required,role"`
}

type UserUpdateRequest struct {
	Name  string `json:"name" validate:"omitempty,min=3,max=50"`
	Email string `json:"email" validate:"omitempty,email"`
}
