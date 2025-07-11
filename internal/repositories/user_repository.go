package repositories

import (
	"errors"
	"strings"

	"github.com/anieswahdie1/ara-medika-api.git/internal/models/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entities.Users) error
	FindByID(id uint) (*entities.Users, error)
	FindByEmail(email string) (*entities.Users, error)
	Update(user *entities.Users) error
	Delete(id uint) error
	FindAll(limit, offset int) ([]entities.Users, error)
	FindAllMenus(roles string) ([]entities.Menus, error)
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

func (r *userRepository) FindByID(id uint) (*entities.Users, error) {
	var user entities.Users
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*entities.Users, error) {
	var user entities.Users
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *entities.Users) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&entities.Users{}, id).Error
}

func (r *userRepository) FindAll(limit, offset int) ([]entities.Users, error) {
	var users []entities.Users
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repos *userRepository) FindAllMenus(roles string) ([]entities.Menus, error) {
	var (
		menus, newListMenus []entities.Menus
	)
	err := repos.db.Order("priority asc").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	if len(menus) > 0 {
		for _, menu := range menus {
			newMenu := strings.Split(menu.CanAccessBy, ",")
			for _, item := range newMenu {
				if item == roles {
					newListMenus = append(newListMenus, menu)
				}
			}
		}
	}
	return newListMenus, nil
}
