package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var allowedOrigins = map[string]struct{}{
	// "https://amigoscare.club":           {},
	// "https://amigoscareclub.vercel.app": {},
	"*": {}, // Allow all origins (for development only, remove in production)
}

func isAllowedOrigin(origin string) bool {
	if origin == "" {
		return true
	}
	if _, ok := allowedOrigins[origin]; ok {
		return true
	}
	if strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "https://localhost") {
		return true
	}
	return false
}

// CORS sets basic CORS headers for known origins.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if isAllowedOrigin(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
