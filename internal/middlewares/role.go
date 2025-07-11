package middlewares

import (
	"net/http"

	"github.com/anieswahdie1/ara-medika-api.git/internal/models/responses"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var tagMiddlewareRole = "internal.middlewares.role."

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	var (
		tag        = tagMiddlewareRole + "RoleMiddleware."
		errMessage responses.ErrorResponse
	)

	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("role")
		if !exists {

			logrus.WithFields(logrus.Fields{
				"tag":   tag + "01",
				"error": "role not found",
			})

			errMessage.Error = "Role not found"
			ctx.AbortWithStatusJSON(http.StatusForbidden, responses.Responses{
				Code:        http.StatusForbidden,
				Description: "FORBIDDEN",
				Data:        errMessage,
			})
			return
		}

		allowed := false
		for _, role := range allowedRoles {
			if role == userRole {
				allowed = true
				break
			}
		}

		if !allowed {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "02",
				"error": "user not permitted!",
			})

			errMessage.Error = "Insufficent permissions"
			ctx.AbortWithStatusJSON(http.StatusForbidden, responses.Responses{
				Code:        http.StatusForbidden,
				Description: "FORBIDDEN",
				Data:        errMessage,
			})
			return
		}

		ctx.Next()
	}
}
