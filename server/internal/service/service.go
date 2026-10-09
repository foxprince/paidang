package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/foxprince/paidang/server/internal/config"
	"github.com/foxprince/paidang/server/internal/model"
	"github.com/foxprince/paidang/server/internal/repository"
)

// Service 业务逻辑层
type Service struct {
	store *repository.Store
	cfg   *config.Config
}

func New(store *repository.Store, cfg *config.Config) *Service {
	return &Service{store: store, cfg: cfg}
}

// GetOrCreateCoach 按 openid 找陪练，没有则建档
func (s *Service) GetOrCreateCoach(openid string) (*model.Coach, bool, error) {
	var coach model.Coach
	res := s.store.DB().Where("openid = ?", openid).First(&coach)
	if res.Error == nil {
		return &coach, false, nil
	}
	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, false, res.Error
	}
	coach = model.Coach{OpenID: openid, Nickname: "微信用户"}
	if err := s.store.DB().Create(&coach).Error; err != nil {
		return nil, false, err
	}
	return &coach, true, nil
}
