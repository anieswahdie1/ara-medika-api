package routes

import (
	"github.com/anieswahdie1/ara-medika-api.git/internal/configs"
	"github.com/anieswahdie1/ara-medika-api.git/internal/controllers"
	"github.com/anieswahdie1/ara-medika-api.git/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupAuthRoutes(
	router *gin.Engine,
	cfg *configs.Config,
	redisClient *redis.Client,
	authController *controllers.AuthController,
) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/refresh", authController.RefreshToken)

		// protected routes
		authGroup.Use(middlewares.AuthMiddleware(cfg, redisClient))
		{
			authGroup.POST("/logout", authController.Logout)
			authGroup.GET("/me", authController.GetCurrentUser)
		}
	}
}
