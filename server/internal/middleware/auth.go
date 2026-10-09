package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/foxprince/paidang/server/pkg/response"
)

// JWT 校验，coach_id 注入 context
func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		token := strings.TrimPrefix(h, "Bearer ")
		if token == "" || token == h {
			response.Err(c, 40101, "未登录")
			c.Abort()
			return
		}
		t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil || !t.Valid {
			response.Err(c, 40101, "登录已过期")
			c.Abort()
			return
		}
		if claims, ok := t.Claims.(jwt.MapClaims); ok {
			if id, ok := claims["coach_id"].(float64); ok {
				c.Set("coach_id", int64(id))
			}
		}
		c.Next()
	}
}
