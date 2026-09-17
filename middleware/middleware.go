package middleware

import (
	"net/http"
	"strings"

	"github.com/arunkumar-1311/mongo-db/auth"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserId    = "userId"
	ContextEmail     = "email"
	ContextRoleLevel = "roleLevel"
)

// Authenticate validates the Bearer JWT on the request and stores the
// claims in the Gin context for downstream handlers/middleware.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "authorization header is required"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "authorization header must be in the form: Bearer <token>"})
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
			return
		}

		c.Set(ContextUserId, claims.UserId)
		c.Set(ContextEmail, claims.Email)
		c.Set(ContextRoleLevel, claims.RoleLevel)
		c.Next()
	}
}

// RequireRole restricts access to callers whose token role level is in the
// allowed set. Must run after Authenticate.
func RequireRole(allowed ...string) gin.HandlerFunc {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, r := range allowed {
		allowedSet[strings.ToLower(r)] = struct{}{}
	}

	return func(c *gin.Context) {
		roleLevel, _ := c.Get(ContextRoleLevel)
		roleStr, _ := roleLevel.(string)

		if _, ok := allowedSet[strings.ToLower(roleStr)]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "you do not have permission to perform this action"})
			return
		}
		c.Next()
	}
}
