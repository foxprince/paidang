package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/foxprince/paidang/server/pkg/mask"
	"github.com/foxprince/paidang/server/pkg/response"
)

// ListStudents 学员列表
func (h *Handler) ListStudents(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.ListStudents(h.coachID(c), page, size)
	if err != nil {
		response.Err(c, 50001, "服务异常")
		return
	}
	for i := range list {
		list[i].Phone = mask.Phone(list[i].Phone)
	}
	response.OK(c, pageData(list, total))
}

// CreateStudent 手动建档
func (h *Handler) CreateStudent(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		Phone  string `json:"phone"`
		Wechat string `json:"wechat"`
		Level  string `json:"level"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 40001, "学员姓名必填")
		return
	}
	st, err := h.svc.CreateStudent(h.coachID(c), req.Name, req.Phone, req.Wechat, req.Level)
	if err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	response.OK(c, st)
}

// GetStudent 学员详情 + 按次课备注时间线
func (h *Handler) GetStudent(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	d, err := h.svc.GetStudentDetail(h.coachID(c), id)
	if err != nil {
		response.Err(c, 40401, err.Error())
		return
	}
	d.Student.Phone = mask.Phone(d.Student.Phone)
	response.OK(c, d)
}

// AddNote 写按次课备注
func (h *Handler) AddNote(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		BookingID *int64 `json:"booking_id"`
		KeyPoints string `json:"key_points" binding:"required"`
		NextPlan  string `json:"next_plan"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 40001, "本节要点不能为空")
		return
	}
	note, err := h.svc.AddNote(h.coachID(c), id, req.BookingID, req.KeyPoints, req.NextPlan)
	if err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	response.OK(c, note)
}
