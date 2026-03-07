package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/synapse/server/internal/apierror"
)

// ErrorHandler is a Gin middleware that recovers from panics and maps domain
// errors to appropriate HTTP status codes. It must be registered BEFORE routes.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered", "panic", r)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}()

		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var apiErr *apierror.APIError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.Code, gin.H{"error": apiErr.Message})
			return
		}

		slog.Error("unhandled error", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
