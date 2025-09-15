package middleware

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AdminOnly es un middleware que permite el acceso solo a administradores
func AdminOnly(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización requerido"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudieron leer los claims del token"})
			c.Abort()
			return
		}

		userType, ok := claims["user_type"].(string)
		if !ok || userType != "administrador" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acceso solo para administradores"})
			c.Abort()
			return
		}

		c.Next()
	}
}
