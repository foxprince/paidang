package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/foxprince/paidang/server/internal/service"
	"github.com/foxprince/paidang/server/internal/sms"
	"github.com/foxprince/paidang/server/internal/wechat"
	"github.com/foxprince/paidang/server/pkg/response"
)

// PublicCoach 学员端：陪练主页展示
func (h *Handler) PublicCoach(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	page, err := h.svc.PublicCoach(id)
	if err != nil {
		response.Err(c, 40401, err.Error())
		return
	}
	response.OK(c, page)
}

// PublicSlots 学员端：?date= 可约时段
func (h *Handler) PublicSlots(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	slots, err := h.svc.AvailableSlots(id, c.Query("date"))
	if err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	response.OK(c, slots)
}

// SmsCode 发送短信验证码（限流 1/分钟由网关层控制，这里做发送间隔保护）
func (h *Handler) SmsCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Phone) != 11 {
		response.Err(c, 40001, "手机号格式错误")
		return
	}
	code := sms.NewCode()
	h.codes.Set(req.Phone, code, 5*time.Minute)
	if err := h.sms.Send(req.Phone, code); err != nil {
		response.Err(c, 50001, "发送失败")
		return
	}
	response.OK(c, gin.H{"sent": true})
}

// PublicBooking 学员下单
func (h *Handler) PublicBooking(c *gin.Context) {
	var req struct {
		CoachID  int64  `json:"coach_id" binding:"required"`
		Date     string `json:"date" binding:"required"`
		Start    string `json:"start" binding:"required"`
		End      string `json:"end" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		SmsCode  string `json:"sms_code" binding:"required"`
		IdemKey  string `json:"idem_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 40001, "参数错误")
		return
	}
	if !h.codes.Verify(req.Phone, req.SmsCode) {
		response.Err(c, 40001, "验证码错误或已过期")
		return
	}
	var coach, _ = h.svc.GetCoach(req.CoachID)
	if coach == nil {
		response.Err(c, 40401, "陪练不存在")
		return
	}
	b, err := h.svc.CreateBooking(service.CreateBookingInput{
		CoachID: req.CoachID, Date: req.Date, Start: req.Start, End: req.End,
		Name: req.Name, Phone: req.Phone,
		StudentOpenID: h.openidFromRequest(c),
		IdemKey:       req.IdemKey, Price: coach.PricePerHour,
	})
	if err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	// 通知陪练：有新预约
	h.notify(&coach.OpenID, "new_booking", map[string]map[string]string{
		"thing1": wechat.Str("收到新预约"),
		"time2":  wechat.Str(b.PlayDate + " " + b.StartTime + "-" + b.EndTime),
		"thing3": wechat.Str(req.Name),
	})
	response.OK(c, b)
}

// PublicReview 学员评价（订单完成后）
func (h *Handler) PublicReview(c *gin.Context) {
	var req struct {
		BookingID int64  `json:"booking_id" binding:"required"`
		Rating    int    `json:"rating" binding:"required"`
		Comment   string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 40001, "参数错误")
		return
	}
	rv, err := h.svc.CreateReview(req.BookingID, req.Rating, req.Comment)
	if err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	response.OK(c, rv)
}
