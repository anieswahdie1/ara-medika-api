package middlewares

import (
	"log"

	"github.com/anieswahdie1/ara-medika-api.git/internal/errors"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next() // Process request

		// Handle errors from handlers
		errs := ctx.Errors
		if len(errs) > 0 {
			err := errs[0].Err

			switch e := err.(type) {
			case *errors.APIError:
				ctx.JSON(e.Status, e)
			default:
				log.Printf("Unhandled error: %v", err)
				apiErr := errors.NewInternalServerError(errors.CodeInternalError, "INTERNAL_SERVER_ERROR")
				ctx.JSON(apiErr.Status, apiErr)
			}
			return
		}
	}
}
