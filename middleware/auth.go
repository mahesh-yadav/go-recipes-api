package middleware

import "github.com/gin-gonic/gin"

const APIKeyHeader = "X-API-KEY"
const X_API_KEY = "WwpPG9VkonfCOp6jZZUuIA=="

func AuthMiddlewareAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(APIKeyHeader)
		if apiKey != X_API_KEY {
			c.AbortWithStatus(401)
		}
		c.Next()
	}
}
