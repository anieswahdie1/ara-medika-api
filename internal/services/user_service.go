package services

import (
	"github.com/anieswahdie1/ara-medika-api.git/internal/models/entities"
	"github.com/anieswahdie1/ara-medika-api.git/internal/repositories"
	"github.com/anieswahdie1/ara-medika-api.git/internal/utils"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	CreateUser(user *entities.Users) error
}

type userService struct {
	userRepo repositories.UserRepository
	logger   *logrus.Logger
}

func NewUserService(userRepo repositories.UserRepository, logger *logrus.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *userService) CreateUser(user *entities.Users) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		s.logger.Errorf("Failed to hash password: %v", err)
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Create(user)
}
