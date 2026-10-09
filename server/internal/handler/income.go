package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/foxprince/paidang/server/pkg/response"
)

// IncomeSummary 收入看板 ?month=2026-10
func (h *Handler) IncomeSummary(c *gin.Context) {
	sum, err := h.svc.IncomeSummary(h.coachID(c), c.Query("month"))
	if err != nil {
		response.Err(c, 50001, "服务异常")
		return
	}
	response.OK(c, sum)
}
