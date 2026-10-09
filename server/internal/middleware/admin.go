package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/foxprince/paidang/server/internal/model"
	"github.com/foxprince/paidang/server/internal/repository"
	"github.com/foxprince/paidang/server/pkg/response"
)

// AdminOnly 管理员鉴权：必须在 JWT 中间件之后使用
func AdminOnly(store *repository.Store, admins map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("coach_id")
		if !ok {
			response.Err(c, 40301, "无权限")
			c.Abort()
			return
		}
		var coach model.Coach
		if err := store.DB().First(&coach, v).Error; err != nil {
			response.Err(c, 40301, "无权限")
			c.Abort()
			return
		}
		if !admins[coach.OpenID] {
			response.Err(c, 40301, "无权限")
			c.Abort()
			return
		}
		c.Set("admin_openid", coach.OpenID)
		c.Next()
	}
}
