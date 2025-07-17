package routes

import (
	"github.com/anieswahdie1/ara-medika-api.git/internal/configs"
	"github.com/anieswahdie1/ara-medika-api.git/internal/controllers"
	"github.com/anieswahdie1/ara-medika-api.git/internal/middlewares"
	"github.com/anieswahdie1/ara-medika-api.git/internal/models/entities"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupUserRoutes(
	router *gin.Engine,
	cfg *configs.Config,
	redisClient *redis.Client,
	userController *controllers.UserController,
) {
	userGroup := router.Group("/users")
	userGroup.Use(middlewares.AuthMiddleware(cfg, redisClient))
	{
		// Routes untuk semua user terautentikasi
		userGroup.GET("/me", userController.GetUserByID)
		userGroup.PUT("/me", userController.UpdateUser)
		// Routes untuk admin dan super_admin
		userGroup.Use(middlewares.RoleMiddleware(string(entities.Admin), string(entities.SuperAdmin)))
		{
			userGroup.POST("/", userController.CreateUser)
			userGroup.GET("/", userController.GetListUser)
		}

	}
}
