package services

import (
	"errors"

	"github.com/anieswahdie1/ara-medika-api.git/internal/models/entities"
	"github.com/anieswahdie1/ara-medika-api.git/internal/repositories"
	"github.com/anieswahdie1/ara-medika-api.git/internal/utils"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	CreateUser(user *entities.Users) error
	GetUserByID(id uint) (*entities.Users, error)
	GetUserByEmail(email string) (*entities.Users, error)
	UpdateUser(user *entities.Users) error
	DeleteUser(id uint) error
	ListUsers(limit, offset int) ([]entities.Users, error)
	ChangePassword(userID uint, oldPassword, newPassword string) error
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
	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err != nil {
		s.logger.Errorf("Error checking email existence: %v", err)
		return errors.New("failed to check email availability")
	}
	if existingUser != nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		s.logger.Errorf("Failed to hash password: %v", err)
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Create(user)
}

func (s *userService) GetUserByID(id uint) (*entities.Users, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		s.logger.Errorf("Failed to get user by ID %d: %v", id, err)
		return nil, errors.New("failed to get user")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) GetUserByEmail(email string) (*entities.Users, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		s.logger.Errorf("Failed to get user by email %s: %v", email, err)
		return nil, errors.New("failed to get user")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) UpdateUser(user *entities.Users) error {
	existingUser, err := s.userRepo.FindByID(user.ID)
	if err != nil {
		s.logger.Errorf("Failed to find user %d for update: %v", user.ID, err)
		return errors.New("failed to find user")
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// Prevent updating certain fields
	user.Password = existingUser.Password
	user.Role = existingUser.Role

	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	// Check if user exists first
	existingUser, err := s.userRepo.FindByID(id)
	if err != nil {
		s.logger.Errorf("Failed to find user %d for deletion: %v", id, err)
		return errors.New("failed to find user")
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}

func (s *userService) ListUsers(limit, offset int) ([]entities.Users, error) {
	users, err := s.userRepo.FindAll(limit, offset)
	if err != nil {
		s.logger.Errorf("Failed to list users: %v", err)
		return nil, errors.New("failed to list users")
	}
	return users, nil
}

func (s *userService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Errorf("Failed to find user %d for password change: %v", userID, err)
		return errors.New("failed to find user")
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Verify old password
	if !utils.CheckPassword(oldPassword, user.Password) {
		return errors.New("incorrect old password")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		s.logger.Errorf("Failed to hash new password: %v", err)
		return errors.New("failed to process password")
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}
