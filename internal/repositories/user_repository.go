package repositories

import (
	"errors"
	"strings"

	"github.com/anieswahdie1/ara-medika-api.git/internal/models/entities"
	"github.com/anieswahdie1/ara-medika-api.git/internal/models/requests"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entities.Users) error
	FindByID(id uint) (*entities.Users, error)
	FindByEmail(email string) (*entities.Users, error)
	Update(user *entities.Users) error
	Delete(id uint) error
	FindUsers(request requests.BaseGetListRequest) ([]entities.Users, error)
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

func (r *userRepository) FindUsers(request requests.BaseGetListRequest) ([]entities.Users, error) {
	var (
		users  []entities.Users
		offset int
	)

	offset = (request.Page - 1) * request.Limit

	queryBuilder := r.db.
		Limit(request.Limit).
		Offset(offset).
		Order("created_at DESC").
		Where("active = ?", "true").
		Find(&users)

	if request.Search != "" {
		queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+request.Search+"%")
	}

	if err := queryBuilder.Error; err != nil {
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
