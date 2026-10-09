package service

import (
	"gorm.io/gorm"

	"github.com/foxprince/paidang/server/internal/model"
)

// logAdmin 写审计日志（可用普通 DB 或事务）
func (s *Service) logAdmin(db *gorm.DB, adminOpenID, action, targetType string, targetID int64, detail string) error {
	return db.Create(&model.AdminLog{
		AdminOpenID: adminOpenID, Action: action,
		TargetType: targetType, TargetID: targetID, Detail: detail,
	}).Error
}
