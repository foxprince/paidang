package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/foxprince/paidang/server/pkg/mask"
	"github.com/foxprince/paidang/server/pkg/response"
)

// GetMe 我的档案
func (h *Handler) GetMe(c *gin.Context) {
	coach, err := h.svc.GetCoach(h.coachID(c))
	if err != nil {
		response.Err(c, 40401, "不存在")
		return
	}
	response.OK(c, coach)
}

// UpdateProfile 编辑主页
func (h *Handler) UpdateProfile(c *gin.Context) {
	var patch map[string]any
	if err := c.ShouldBindJSON(&patch); err != nil {
		response.Err(c, 40001, "参数错误")
		return
	}
	coach, err := h.svc.UpdateProfile(h.coachID(c), patch)
	if err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	response.OK(c, coach)
}

// Publish 发布/下架主页
func (h *Handler) Publish(c *gin.Context) {
	var req struct {
		Published bool `json:"published"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, 40001, "参数错误")
		return
	}
	if err := h.svc.Publish(h.coachID(c), req.Published); err != nil {
		response.Err(c, 40001, err.Error())
		return
	}
	response.OK(c, gin.H{"published": req.Published})
}

// Today 今日课表
func (h *Handler) Today(c *gin.Context) {
	list, err := h.svc.TodaySchedule(h.coachID(c))
	if err != nil {
		response.Err(c, 50001, "服务异常")
		return
	}
	// 手机号脱敏
	for i := range list {
		list[i].StudentPhone = mask.Phone(list[i].StudentPhone)
	}
	response.OK(c, list)
}
