package middleware

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a configured CORS middleware.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	cfg := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Timezone"},
		ExposeHeaders:    []string{"X-Total-Count"},
		AllowCredentials: true,
	}
	// Validate config; fall back to a permissive default in development.
	if err := cfg.Validate(); err != nil {
		return func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			c.Next()
		}
	}
	return cors.New(cfg)
}
