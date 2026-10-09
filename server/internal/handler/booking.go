package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/foxprince/paidang/server/internal/wechat"
	"github.com/foxprince/paidang/server/pkg/mask"
	"github.com/foxprince/paidang/server/pkg/response"
)

// ListBookings 预约列表 ?status=&date=
func (h *Handler) ListBookings(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.ListBookings(h.coachID(c), c.Query("status"), c.Query("date"), page, size)
	if err != nil {
		response.Err(c, 50001, "服务异常")
		return
	}
	for i := range list {
		list[i].StudentPhone = mask.Phone(list[i].StudentPhone)
	}
	response.OK(c, pageData(list, total))
}

func (h *Handler) bookingID(c *gin.Context) int64 {
	n, err := parseIntParam(c.Param("id"))
	if err != nil || n <= 0 {
		return 0
	}
	return n
}

func parseIntParam(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// Confirm 确认预约
func (h *Handler) ConfirmBooking(c *gin.Context) {
	b, err := h.svc.Transition(h.bookingID(c), h.coachID(c), "confirmed", "", "")
	if err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	h.notify(b.StudentOpenID, "result", map[string]map[string]string{
		"thing1": wechat.Str("预约已确认"),
		"time2":  wechat.Str(b.PlayDate + " " + b.StartTime),
	})
	response.OK(c, b)
}

// RejectBooking 拒绝预约
func (h *Handler) RejectBooking(c *gin.Context) {
	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)
	b, err := h.svc.Transition(h.bookingID(c), h.coachID(c), "cancelled", req.Reason, "coach")
	if err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	h.notify(b.StudentOpenID, "result", map[string]map[string]string{
		"thing1": wechat.Str("预约未通过"),
		"time2":  wechat.Str(b.PlayDate + " " + b.StartTime),
	})
	response.OK(c, b)
}

// CompleteBooking 标记完成
func (h *Handler) CompleteBooking(c *gin.Context) {
	b, err := h.svc.Transition(h.bookingID(c), h.coachID(c), "completed", "", "")
	if err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	h.notify(b.StudentOpenID, "done", map[string]map[string]string{
		"thing1": wechat.Str("课程已完成，欢迎评价"),
		"time2":  wechat.Str(b.PlayDate + " " + b.StartTime),
	})
	response.OK(c, b)
}

// CancelBooking 取消预约（记爽约方）
func (h *Handler) CancelBooking(c *gin.Context) {
	var req struct {
		Reason   string `json:"reason"`
		NoShowBy string `json:"no_show_by"` // coach / student
	}
	_ = c.ShouldBindJSON(&req)
	if req.NoShowBy == "" {
		req.NoShowBy = "coach"
	}
	b, err := h.svc.Transition(h.bookingID(c), h.coachID(c), "cancelled", req.Reason, req.NoShowBy)
	if err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	response.OK(c, b)
}

// MarkPaid 确认收款
func (h *Handler) MarkPaid(c *gin.Context) {
	b, err := h.svc.MarkPaid(h.bookingID(c), h.coachID(c))
	if err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	response.OK(c, b)
}
