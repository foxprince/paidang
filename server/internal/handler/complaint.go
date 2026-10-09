package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/foxprince/paidang/server/pkg/response"
)

// CreateComplaint 投诉
func (h *Handler) CreateComplaint(c *gin.Context) {
	var req struct {
		BookingID int64    `json:"booking_id" binding:"required"`
		Evidence  []string `json:"evidence"`
		Detail    string   `json:"detail"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 40001, "参数错误")
		return
	}
	cp, err := h.svc.CreateComplaint(req.BookingID, req.Evidence, req.Detail)
	if err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	response.OK(c, cp)
}
