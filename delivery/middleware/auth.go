package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"myapp/application/serviceInterface"
	"myapp/domain"
)

// AdminClaimsKey is the context key under which the parsed token claims are stored.
const AdminClaimsKey = "adminClaims"

// Auth validates the Bearer token and injects the claims into the Gin context.
func Auth(tokenSvc serviceInterface.TokenService, tokenRepo domain.TokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Token tidak valid.",
			}})
			c.Abort()
			return
		}
		raw := strings.TrimPrefix(header, "Bearer ")

		claims, err := tokenSvc.Parse(c.Request.Context(), raw)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Token tidak valid atau sudah habis masa berlakunya.",
			}})
			c.Abort()
			return
		}

		revoked, err := tokenRepo.IsRevoked(c.Request.Context(), claims.JTI)
		if err != nil || revoked {
			c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Token tidak valid.",
			}})
			c.Abort()
			return
		}

		c.Set(AdminClaimsKey, claims)
		c.Next()
	}
}
