package middleware

import (
	"net/http"
	"strings"
	"task-api/internal/domain"
	"task-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func Auth(tokens usecase.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrInvalidToken.Error()})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := tokens.Parse(c.Request.Context(), tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrInvalidToken.Error()})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("tokenID", claims.TokenID)
		c.Set("mustChangePassword", claims.MustChangePassword)
		c.Next()
	}
}

func RequirePasswordChanged() gin.HandlerFunc {
	return func(c *gin.Context) {
		mustChange, ok := c.Get("mustChangePassword")
		if ok {
			if val, isBool := mustChange.(bool); isBool && val {
				c.JSON(http.StatusForbidden, gin.H{
					"error": domain.ErrPasswordChangeRequired.Error(),
					"code":  "PASSWORD_CHANGE_REQUIRED",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func RequireRole(roles ...domain.Role) gin.HandlerFunc {
	allowed := make(map[domain.Role]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrInvalidToken.Error()})
			c.Abort()
			return
		}

		role, ok := roleVal.(domain.Role)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrInvalidToken.Error()})
			c.Abort()
			return
		}

		if !allowed[role] {
			c.JSON(http.StatusForbidden, gin.H{"error": domain.ErrForbidden.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}
