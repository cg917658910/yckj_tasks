package middleware

import (
	"strings"
	"task-system-go/internal/util"

	"github.com/gin-gonic/gin"
)

func AdminAuth(jwt util.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := parseBearer(c, jwt)
		if !ok || claims["role"] != "admin" {
			util.JSONError(c, "未登录", 401)
			c.Abort()
			return
		}
		c.Set("admin_id", claims["admin_id"])
		c.Next()
	}
}

func UserAuth(jwt util.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := parseBearer(c, jwt)
		if !ok || claims["role"] != "user" {
			util.JSONError(c, "未登录", 401)
			c.Abort()
			return
		}
		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}

func parseBearer(c *gin.Context, jwt util.JWTManager) (map[string]interface{}, bool) {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return nil, false
	}
	token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
	claims, err := jwt.Parse(token)
	if err != nil {
		return nil, false
	}
	return claims, true
}
