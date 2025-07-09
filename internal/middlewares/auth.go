package middlewares

import (
	"net/http"
	"strings"

	"github.com/anieswahdie1/ara-medika-api.git/internal/configs"
	"github.com/anieswahdie1/ara-medika-api.git/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(cfg *configs.Config, redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]

		// Check if token is blacklisted in Redis
		val, err := redisClient.Get(ctx, tokenString).Result()
		if err == nil && val == "blacklisted" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been invalidated"})
			return
		}

		claims, err := utils.ValidateToken(cfg, tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			return
		}

		ctx.Set("userID", claims.UserID)
		ctx.Set("email", claims.Email)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
