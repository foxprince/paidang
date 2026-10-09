package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/foxprince/paidang/server/internal/config"
	"github.com/foxprince/paidang/server/internal/service"
	"github.com/foxprince/paidang/server/internal/sms"
	"github.com/foxprince/paidang/server/internal/wechat"
	"github.com/foxprince/paidang/server/pkg/response"
)

type Handler struct {
	svc   *service.Service
	cfg   *config.Config
	sms   sms.Sender
	codes *sms.CodeStore
	wxsub *wechat.SubscribeClient
	tpl   map[string]string // 订阅消息模板 ID
}

func New(svc *service.Service, cfg *config.Config) *Handler {
	return &Handler{
		svc:   svc,
		cfg:   cfg,
		sms:   sms.LogSender{},
		codes: sms.NewCodeStore(),
		wxsub: wechat.NewSubscribeClient(cfg.WechatAppID, cfg.WechatSecret),
		tpl: map[string]string{
			"new_booking": cfg.TplNewBooking,
			"result":      cfg.TplBookingResult,
			"done":        cfg.TplBookingDone,
		},
	}
}

// coachID 当前登录的陪练
func (h *Handler) coachID(c *gin.Context) int64 {
	if v, ok := c.Get("coach_id"); ok {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

func pageParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return page, size
}

func pageData(list any, total int64) gin.H {
	return gin.H{"list": list, "total": total}
}

// openidFromRequest 尝试从 token 解析 openid（学员端可选登录）
func (h *Handler) openidFromRequest(c *gin.Context) *string {
	auth := c.GetHeader("Authorization")
	if len(auth) < 8 || auth[:7] != "Bearer " {
		return nil
	}
	t, err := jwt.Parse(auth[7:], func(t *jwt.Token) (any, error) {
		return []byte(h.cfg.JWTSecret), nil
	})
	if err != nil || !t.Valid {
		return nil
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}
	id, ok := claims["coach_id"].(float64)
	if !ok {
		return nil
	}
	coach, err := h.svc.GetCoach(int64(id))
	if err != nil {
		return nil
	}
	return &coach.OpenID
}

// notify 发订阅消息，失败只记日志
func (h *Handler) notify(openid *string, tplKey string, data map[string]map[string]string) {
	if openid == nil || *openid == "" || h.tpl[tplKey] == "" {
		return
	}
	if err := h.wxsub.Send(*openid, h.tpl[tplKey], data); err != nil {
		// 不阻断业务
		println("[notify] failed:", err.Error())
	}
}

func (h *Handler) NotImpl(c *gin.Context) { response.NotImpl(c) }
