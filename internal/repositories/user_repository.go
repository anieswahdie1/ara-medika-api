package repositories

import (
	"github.com/anieswahdie1/ara-medika-api.git/internal/models/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entities.Users) error
	// FindByID(id uint) (*entities.Users, error)
	// FindByEmail(email string) (*entities.Users, error)
	// Update(user *entities.Users) error
	// Delete(id uint) error
	// FindAll() ([]entities.Users, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entities.Users) error {
	return r.db.Create(user).Error
}
