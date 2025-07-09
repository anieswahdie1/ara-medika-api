package controllers

import (
	"net/http"
	"strings"

	"github.com/anieswahdie1/ara-medika-api.git/internal/models/responses"
	"github.com/anieswahdie1/ara-medika-api.git/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
	logger      *logrus.Logger
}

func NewAuthController(
	authService services.AuthService,
	userService services.UserService,
	logger *logrus.Logger,
) *AuthController {
	return &AuthController{
		authService: authService,
		userService: userService,
		logger:      logger,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// @Summary Login user
// @Description Authenticate user and get JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Login credentials"
// @Success 200 {object} responses.TokenResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Errorf("Invalid login request: %v", err)
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error: "Invalid request format",
		})
		return
	}

	accessToken, refreshToken, err := c.authService.Login(req.Email, req.Password)
	if err != nil {
		c.logger.Warnf("Login failed for email %s: %v", req.Email, err)
		ctx.JSON(http.StatusUnauthorized, responses.ErrorResponse{
			Error: "Invalid email or password",
		})
		return
	}

	ctx.JSON(http.StatusOK, responses.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// @Summary Refresh JWT token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} responses.TokenResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error: "Invalid request format",
		})
		return
	}

	accessToken, refreshToken, err := c.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.logger.Warnf("Refresh token failed: %v", err)
		ctx.JSON(http.StatusUnauthorized, responses.ErrorResponse{
			Error: "Invalid refresh token",
		})
		return
	}

	ctx.JSON(http.StatusOK, responses.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// @Summary Logout user
// @Description Invalidate user's JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Error: "Authorization header is required",
		})
		return
	}

	tokenString := strings.Split(authHeader, " ")[1]
	userID := ctx.MustGet("userID").(uint)

	if err := c.authService.Logout(tokenString, userID); err != nil {
		c.logger.Errorf("Logout failed for user %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error: "Failed to logout",
		})
		return
	}

	ctx.JSON(http.StatusOK, responses.SuccessResponse{
		Message: "Successfully logged out",
	})
}

// GetCurrentUser godoc
// @Summary Get current user profile
// @Description Get profile of currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.UserResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /auth/me [get]
func (c *AuthController) GetCurrentUser(ctx *gin.Context) {
	// Dapatkan userID dari context yang sudah diset oleh middleware
	userID, exists := ctx.Get("userID")
	if !exists {
		c.logger.Error("UserID not found in context")
		ctx.JSON(http.StatusUnauthorized, responses.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	// Dapatkan data user dari service
	user, err := c.userService.GetUserByID(userID.(uint))
	if err != nil {
		if err.Error() == "user not found" {
			c.logger.Errorf("User not found for ID %d", userID)
			ctx.JSON(http.StatusUnauthorized, responses.ErrorResponse{
				Error: "User not found",
			})
			return
		}
		c.logger.Errorf("Failed to get user %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error: "Failed to get user data",
		})
		return
	}

	// Format response (tampilkan hanya data yang diperlukan)
	ctx.JSON(http.StatusOK, responses.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	})
}
