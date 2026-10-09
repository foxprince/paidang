package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/foxprince/paidang/server/internal/service"
	"github.com/foxprince/paidang/server/pkg/mask"
	"github.com/foxprince/paidang/server/pkg/response"
)

func (h *Handler) adminOpenID(c *gin.Context) string {
	if v, ok := c.Get("admin_openid"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// AdminPing 是否管理员（能进到这里就是）
func (h *Handler) AdminPing(c *gin.Context) {
	response.OK(c, gin.H{"is_admin": true})
}

// AdminListCoaches 陪练列表 ?verified=&status=
func (h *Handler) AdminListCoaches(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.AdminListCoaches(c.Query("verified"), c.Query("status"), page, size)
	if err != nil {
		response.Err(c, 50001, "服务异常")
		return
	}
	for i := range list {
		list[i].Phone = mask.Phone(list[i].Phone)
	}
	response.OK(c, pageData(list, total))
}

// AdminCoachDetail 陪练详情
func (h *Handler) AdminCoachDetail(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	d, err := h.svc.AdminCoachDetail(id)
	if err != nil {
		response.Err(c, 40401, err.Error())
		return
	}
	d.Coach.Phone = mask.Phone(d.Coach.Phone)
	response.OK(c, d)
}

// AdminVerifyCoach 审核
func (h *Handler) AdminVerifyCoach(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		Verified bool   `json:"verified"`
		Note     string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 40001, "参数错误")
		return
	}
	if err := h.svc.AdminVerifyCoach(h.adminOpenID(c), id, req.Verified, req.Note); err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	response.OK(c, gin.H{"verified": req.Verified})
}

// AdminSetCoachStatus 改状态
func (h *Handler) AdminSetCoachStatus(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		Status string `json:"status" binding:"required"`
		Note   string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 40001, "参数错误")
		return
	}
	if err := h.svc.AdminSetCoachStatus(h.adminOpenID(c), id, req.Status, req.Note); err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	response.OK(c, gin.H{"status": req.Status})
}

// AdminListComplaints 投诉列表
func (h *Handler) AdminListComplaints(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.AdminListComplaints(c.Query("status"), page, size)
	if err != nil {
		response.Err(c, 50001, "服务异常")
		return
	}
	now := time.Now()
	out := make([]gin.H, 0, len(list))
	for _, it := range list {
		overdue := it.Complaint.Status == "open" && now.Sub(it.Complaint.CreatedAt) > 48*time.Hour
		out = append(out, gin.H{
			"id": it.Complaint.ID, "booking_id": it.Complaint.BookingID,
			"coach_name": it.CoachName, "booking_date": it.BookingDate,
			"status": it.Complaint.Status, "overdue": overdue,
			"created_at": it.Complaint.CreatedAt,
		})
	}
	response.OK(c, gin.H{"list": out, "total": total})
}

// AdminComplaintDetail 投诉详情
func (h *Handler) AdminComplaintDetail(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	d, err := h.svc.AdminComplaintDetail(id)
	if err != nil {
		response.Err(c, 40401, err.Error())
		return
	}
	d.Booking.StudentPhone = mask.Phone(d.Booking.StudentPhone)
	response.OK(c, d)
}

// AdminJudgeComplaint 裁决
func (h *Handler) AdminJudgeComplaint(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		Result string `json:"result" binding:"required"` // upheld / rejected
		Note   string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 40001, "参数错误")
		return
	}
	if err := h.svc.AdminJudgeComplaint(h.adminOpenID(c), id, req.Result, req.Note); err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	response.OK(c, gin.H{"result": req.Result})
}

// AdminStats 数据看板
func (h *Handler) AdminStats(c *gin.Context) {
	st, err := h.svc.AdminStats()
	if err != nil {
		response.Err(c, 50001, "服务异常")
		return
	}
	response.OK(c, st)
}

// AdminCreateCoach 管理员手动建档，返回认领码
func (h *Handler) AdminCreateCoach(c *gin.Context) {
	var req struct {
		Name      string  `json:"name" binding:"required"`
		Phone     string  `json:"phone"`
		Title     string  `json:"title"`
		PriceYuan float64 `json:"price_yuan"`
		Bio       string  `json:"bio"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 40001, "姓名必填")
		return
	}
	coach, err := h.svc.AdminCreateCoach(h.adminOpenID(c), service.CreateCoachInput{
		Name: req.Name, Phone: req.Phone, Title: req.Title,
		PriceYuan: req.PriceYuan, Bio: req.Bio,
	})
	if err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	response.OK(c, coach)
}
