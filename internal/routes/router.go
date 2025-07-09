package routes

import (
	"github.com/anieswahdie1/ara-medika-api.git/internal/configs"
	"github.com/anieswahdie1/ara-medika-api.git/internal/controllers"
	"github.com/anieswahdie1/ara-medika-api.git/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func InitRouter(
	cfg *configs.Config,
	redisClient *redis.Client,
	logger *logrus.Logger,
	userController *controllers.UserController,
	authController *controllers.AuthController,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	// Global middleware
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.RequestLoggerMiddleware(logger))

	// Setup routes
	SetupUserRoutes(router, cfg, redisClient, userController)
	SetupAuthRoutes(router, cfg, redisClient, authController)

	return router
}
