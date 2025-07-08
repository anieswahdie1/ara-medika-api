package main

import (
	"github.com/anieswahdie1/ara-medika-api.git/internal/configs"
	"github.com/anieswahdie1/ara-medika-api.git/internal/controllers"
	"github.com/anieswahdie1/ara-medika-api.git/internal/repositories"
	"github.com/anieswahdie1/ara-medika-api.git/internal/routes"
	"github.com/anieswahdie1/ara-medika-api.git/internal/services"
	"github.com/anieswahdie1/ara-medika-api.git/internal/utils"
	"github.com/anieswahdie1/ara-medika-api.git/pkg/validators"
)

func main() {
	// Load configuration
	cfg := configs.LoadConfig()

	// setup logger
	logger := utils.SetupLogger()

	// Connect to database
	db, err := utils.ConnectDB(cfg)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	// Connect to Redis
	redisClient, err := utils.ConnectRedis(cfg)
	if err != nil {
		logger.Fatalf("Failed to connect to redis: %v", err)
	}

	// Auto migrate models
	// db.AutoMigrate(&entities.User{}, &entities.MasterData{}, ...)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, logger)

	// Initialize controllers
	userController := controllers.NewUserController(userService, logger)

	// Initialize validator
	validators.Init() // Ini akan menginisialisasi validators.Validate

	// Setup router
	router := routes.InitRouter(
		cfg,
		redisClient,
		logger,
		userController,
	)

	// Start server
	logger.Infof("Server is running on port %s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
