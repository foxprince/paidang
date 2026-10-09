package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/foxprince/paidang/server/internal/wechat"
	"github.com/foxprince/paidang/server/pkg/response"
)

// WechatLogin 小程序 wx.login 拿 code 换 token
func (h *Handler) WechatLogin(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 40001, "参数错误")
		return
	}
	sess, err := wechat.Code2Session(h.cfg.WechatAppID, h.cfg.WechatSecret, req.Code)
	if err != nil {
		response.Err(c, 40001, "微信登录失败")
		return
	}
	coach, isNew, err := h.svc.GetOrCreateCoach(sess.OpenID)
	if err != nil {
		response.Err(c, 50001, "服务异常")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"coach_id": coach.ID,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		response.Err(c, 50001, "服务异常")
		return
	}
	response.OK(c, gin.H{"token": signed, "is_new": isNew})
}
