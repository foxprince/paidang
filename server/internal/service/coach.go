package service

import (
	"errors"
	"time"

	"github.com/foxprince/paidang/server/internal/model"
)

// GetCoach 取陪练档案
func (s *Service) GetCoach(id int64) (*model.Coach, error) {
	var c model.Coach
	if err := s.store.DB().First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// profile 可编辑字段白名单
var profileFields = []string{
	"nickname", "avatar_url", "phone", "title", "bio",
	"home_court", "video_url", "price_per_hour", "pay_qr_url", "schedule",
}

// UpdateProfile 编辑主页，只允许白名单字段
func (s *Service) UpdateProfile(id int64, patch map[string]any) (*model.Coach, error) {
	allowed := map[string]any{}
	for _, f := range profileFields {
		if v, ok := patch[f]; ok {
			allowed[f] = v
		}
	}
	if len(allowed) == 0 {
		return nil, errors.New("没有可更新的字段")
	}
	allowed["updated_at"] = time.Now()
	if err := s.store.DB().Model(&model.Coach{}).Where("id = ?", id).Updates(allowed).Error; err != nil {
		return nil, err
	}
	return s.GetCoach(id)
}

// Publish 发布/下架主页。发布时必须有 30 秒对打视频。
func (s *Service) Publish(id int64, published bool) error {
	var c model.Coach
	if err := s.store.DB().First(&c, id).Error; err != nil {
		return err
	}
	if published && c.VideoURL == "" {
		return errors.New("请先上传 30 秒对打视频")
	}
	return s.store.DB().Model(&model.Coach{}).Where("id = ?", id).
		Updates(map[string]any{"published": published, "updated_at": time.Now()}).Error
}

// TodaySchedule 今日课表：待确认置顶，再按时间排
func (s *Service) TodaySchedule(coachID int64) ([]model.Booking, error) {
	today := time.Now().Format("2006-01-02")
	var list []model.Booking
	err := s.store.DB().
		Where("coach_id = ? AND play_date = ? AND status IN ('pending','confirmed')", coachID, today).
		Order("CASE WHEN status='pending' THEN 0 ELSE 1 END, start_time").
		Find(&list).Error
	return list, err
}
