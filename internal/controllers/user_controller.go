package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/anieswahdie1/ara-medika-api.git/internal/models/entities"
	"github.com/anieswahdie1/ara-medika-api.git/internal/models/responses"
	"github.com/anieswahdie1/ara-medika-api.git/internal/services"
	"github.com/anieswahdie1/ara-medika-api.git/pkg/validators"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept json
// @Produce json
// @Param input body entities.UserCreateRequest true "User data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /users [post]
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

	// Hash Password dengan cost yang sesuai
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost, // atau bisa juga dengan bcrypt.MinCost untuk development
	)
	if err != nil {
		fmt.Errorf("hashing failed: %w", err)
		return
	}

	user := entities.Users{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
		Role:     req.Role,
		Active:   true,
	}
	if err := controller.userService.CreateUser(&user); err != nil {
		controller.logger.Errorf("Failed to create user: %v", err)
		ctx.JSON(http.StatusCreated, responses.ErrorResponse{
			Error: err.Error(),
		})
	}

	ctx.JSON(http.StatusCreated, responses.SuccessResponse{
		Message: "User created successfully",
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user details by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entities.User
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	user, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, responses.ErrorResponse{
				Error: "User not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error: "Failed to get user",
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user details
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param input body entities.UserUpdateRequest true "User data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req entities.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Errorf("Failed to bind user data: %v", err)
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error: "Invalid request payload",
		})
		return
	}

	// Validate input
	if err := validators.Validate.Struct(req); err != nil {
		c.logger.Errorf("Validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error: "Validation failed",
		})
		return
	}

	user := &entities.Users{
		Model: entities.Model{ID: uint(id)},
		Name:  req.Name,
		Email: req.Email,
	}

	if err := c.userService.UpdateUser(user); err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, responses.ErrorResponse{
				Error: "User not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, responses.SuccessResponse{
		Message: "User updated successfully",
	})
}
