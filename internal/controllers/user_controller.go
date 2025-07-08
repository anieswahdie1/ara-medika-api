package controllers

import (
	"net/http"

	"github.com/anieswahdie1/ara-medika-api.git/internal/models/entities"
	"github.com/anieswahdie1/ara-medika-api.git/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	userService services.UserService
	logger      *logrus.Logger
}

func NewUserController(userService services.UserService, logger *logrus.Logger) *UserController {
	return &UserController{userService: userService, logger: logger}
}

func (controller *UserController) CreateUser(ctx *gin.Context) {
	var user entities.Users
	if err := ctx.ShouldBindJSON(&user); err != nil {
		controller.logger.Errorf("Failed to bind user data: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := controller.userService.CreateUser(&user); err != nil {
		controller.logger.Errorf("Failed to create user: %v", err)
		ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}
