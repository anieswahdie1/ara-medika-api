package controllers

import (
	"net/http"

	"github.com/anieswahdie1/ara-medika-api.git/internal/models/entities"
	"github.com/anieswahdie1/ara-medika-api.git/internal/services"
	"github.com/anieswahdie1/ara-medika-api.git/pkg/validators"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	userService services.UserService
	logger      *logrus.Logger
	validator   *validator.Validate
}

func NewUserController(userService services.UserService, logger *logrus.Logger) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
		validator:   &validator.Validate{},
	}
}

func (controller *UserController) CreateUser(ctx *gin.Context) {
	var req entities.UserCreateRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		controller.logger.Errorf("Failed to bind user data: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Validasi input menggunakan validators.Validate
	if err := validators.Validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		// Handle validation errors
		for _, e := range validationErrors {
			controller.logger.Errorf("Validation error on field %s: %v", e.Field(), e.Tag())
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
		return
	}

	user := entities.Users{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		Active:   true,
	}
	if err := controller.userService.CreateUser(&user); err != nil {
		controller.logger.Errorf("Failed to create user: %v", err)
		ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}
