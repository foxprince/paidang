package service

import (
	"errors"

	"github.com/foxprince/paidang/server/internal/model"
)

type PublicCoachPage struct {
	Coach     model.Coach   `json:"coach"`
	AvgRating float64       `json:"avg_rating"`
	Reviews   []model.Review `json:"reviews"`
}

// PublicCoach 学员端主页：只展示已发布的主页
func (s *Service) PublicCoach(id int64) (*PublicCoachPage, error) {
	var c model.Coach
	if err := s.store.DB().First(&c, id).Error; err != nil {
		return nil, errors.New("陪练不存在")
	}
	if !c.Published || c.Status != "active" {
		return nil, errors.New("该陪练主页未发布")
	}
	var reviews []model.Review
	_ = s.store.DB().Where("coach_id = ?", id).Order("created_at DESC").Limit(20).Find(&reviews).Error
	var avg float64
	_ = s.store.DB().Model(&model.Review{}).Where("coach_id = ?", id).
		Select("COALESCE(AVG(rating),0)").Scan(&avg).Error
	// 主页不暴露敏感字段
	c.Phone = ""
	return &PublicCoachPage{Coach: c, AvgRating: avg, Reviews: reviews}, nil
}
