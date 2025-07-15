package middlewares

import (
	"net/http"
	"strings"

	"github.com/anieswahdie1/ara-medika-api.git/internal/configs"
	"github.com/anieswahdie1/ara-medika-api.git/internal/errors"
	"github.com/anieswahdie1/ara-medika-api.git/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(cfg *configs.Config, redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.Abort()
			ctx.Error(errors.NewUnauthorizedError(
				errors.CodeUnauthorized,
				"Authorization header is required",
			))
			return
		}

		// Periksa format header dengan lebih hati-hati
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":           "Invalid authorization header format",
				"expected_format": "Bearer <token>",
				"received":        authHeader,
			})
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			ctx.Abort()
			ctx.Error(errors.NewUnauthorizedError(
				errors.CodeUnauthorized,
				"Token cannot be empty",
			))
			return
		}

		// Check if token is blacklisted in Redis
		val, err := redisClient.Get(ctx, tokenString).Result()
		if err == nil && val == "blacklisted" {
			ctx.Abort()
			ctx.Error(errors.NewUnauthorizedError(
				errors.CodeUnauthorized,
				"Token has been invalidated",
			))
			return
		}

		claims, err := utils.ValidateToken(cfg, tokenString)
		if err != nil {
			ctx.Abort()
			ctx.Error(errors.NewUnauthorizedError(
				errors.CodeUnauthorized,
				"Invalid Token",
			))
			return
		}

		ctx.Set("userID", claims.UserID)
		ctx.Set("email", claims.Email)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
