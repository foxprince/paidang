package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/foxprince/paidang/server/internal/model"
)

type CreateCoachInput struct {
	Name      string
	Phone     string
	Title     string
	PriceYuan float64
	Bio       string
}

// genClaimCode 生成 6 位认领码，保证唯一
func (s *Service) genClaimCode(db *gorm.DB) (string, error) {
	for i := 0; i < 10; i++ {
		var b [3]byte
		if _, err := rand.Read(b[:]); err != nil {
			return "", err
		}
		n := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
		code := fmt.Sprintf("%06d", n%1000000)
		var cnt int64
		db.Model(&model.Coach{}).Where("claim_code = ?", code).Count(&cnt)
		if cnt == 0 {
			return code, nil
		}
	}
	return "", errors.New("认领码生成失败，重试")
}

// AdminCreateCoach 管理员手动建档：生成认领码，微信发给陪练认领
func (s *Service) AdminCreateCoach(adminOpenID string, in CreateCoachInput) (*model.Coach, error) {
	if in.Name == "" {
		return nil, errors.New("姓名必填")
	}
	db := s.store.DB()
	code, err := s.genClaimCode(db)
	if err != nil {
		return nil, err
	}
	c := model.Coach{
		OpenID:       "",
		Nickname:     in.Name,
		Phone:        in.Phone,
		Title:        in.Title,
		Bio:          in.Bio,
		PricePerHour: int(in.PriceYuan * 100),
		ClaimCode:    &code,
	}
	if err := db.Create(&c).Error; err != nil {
		return nil, err
	}
	_ = s.logAdmin(db, adminOpenID, "create_coach", "coach", c.ID, "认领码 "+code)
	return &c, nil
}

// ClaimCoach 陪练输入认领码，绑定微信 openid
func (s *Service) ClaimCoach(claimerID int64, code string) (*model.Coach, error) {
	var target model.Coach
	err := s.store.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("claim_code = ? AND openid = ''", code).First(&target).Error; err != nil {
			return errors.New("认领码不对，或已被认领")
		}
		var me model.Coach
		if err := tx.First(&me, claimerID).Error; err != nil {
			return errors.New("登录异常")
		}
		// 自动建档的空档案直接删掉，被认领的档案接管
		if me.Nickname == "微信用户" && !me.Published {
			if err := tx.Delete(&me).Error; err != nil {
				return err
			}
		} else {
			return errors.New("你已有档案，请联系管理员")
		}
		target.OpenID = me.OpenID
		target.ClaimCode = nil
		target.UpdatedAt = time.Now()
		return tx.Save(&target).Error
	})
	if err != nil {
		return nil, err
	}
	return &target, nil
}
