package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/mahesh-yadav/go-recipes-api/models"
	"github.com/rs/zerolog/log"
)

func GlobalErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log error with stack trace
				log.Error().
					Str("module", "recovery").
					Interface("error", err).
					Bytes("stack_trace", debug.Stack()).
					Msg("PANIC RECOVERED")

				// Return JSON error response
				if gin.Mode() == gin.ReleaseMode {
					c.JSON(http.StatusInternalServerError, models.ErrorResponse{
						Code:    http.StatusInternalServerError,
						Message: "Internal Server Error",
					})
				} else {
					c.JSON(http.StatusInternalServerError, models.ErrorResponse{
						Code:    http.StatusInternalServerError,
						Message: fmt.Sprintf("An unexpected error occurred: %v", err),
					})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}
