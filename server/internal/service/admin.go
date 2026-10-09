package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/foxprince/paidang/server/internal/model"
)

// AdminListCoaches 陪练列表。verified: ""全部 "true"已验证 "false"待审核
func (s *Service) AdminListCoaches(verified, status string, page, pageSize int) ([]model.Coach, int64, error) {
	q := s.store.DB().Model(&model.Coach{})
	if verified == "true" {
		q = q.Where("verified = ?", true)
	} else if verified == "false" {
		q = q.Where("verified = ?", false).Where("published = ?", true)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	_ = q.Count(&total).Error
	var list []model.Coach
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

type AdminCoachDetail struct {
	Coach    model.Coach `json:"coach"`
	Bookings int64       `json:"bookings"`
	Reviews  int64       `json:"reviews"`
}

func (s *Service) AdminCoachDetail(id int64) (*AdminCoachDetail, error) {
	var c model.Coach
	if err := s.store.DB().First(&c, id).Error; err != nil {
		return nil, errors.New("陪练不存在")
	}
	var bookings, reviews int64
	_ = s.store.DB().Model(&model.Booking{}).Where("coach_id = ?", id).Count(&bookings).Error
	_ = s.store.DB().Model(&model.Review{}).Where("coach_id = ?", id).Count(&reviews).Error
	return &AdminCoachDetail{Coach: c, Bookings: bookings, Reviews: reviews}, nil
}

// AdminVerifyCoach 审核：通过给验证徽章 / 驳回下架
func (s *Service) AdminVerifyCoach(adminOpenID string, id int64, verified bool, note string) error {
	var c model.Coach
	if err := s.store.DB().First(&c, id).Error; err != nil {
		return errors.New("陪练不存在")
	}
	patch := map[string]any{"verified": verified, "updated_at": time.Now()}
	if !verified {
		patch["published"] = false
	}
	db := s.store.DB()
	if err := db.Model(&model.Coach{}).Where("id = ?", id).Updates(patch).Error; err != nil {
		return err
	}
	action := "verify_coach_pass"
	if !verified {
		action = "verify_coach_reject"
	}
	return s.logAdmin(db, adminOpenID, action, "coach", id, note)
}

// AdminSetCoachStatus 改状态：active / suspended / banned
func (s *Service) AdminSetCoachStatus(adminOpenID string, id int64, status, note string) error {
	if status != "active" && status != "suspended" && status != "banned" {
		return errors.New("状态不合法")
	}
	patch := map[string]any{"status": status, "updated_at": time.Now()}
	if status != "active" {
		patch["published"] = false
	}
	db := s.store.DB()
	if err := db.Model(&model.Coach{}).Where("id = ?", id).Updates(patch).Error; err != nil {
		return err
	}
	return s.logAdmin(db, adminOpenID, "set_status_"+status, "coach", id, note)
}

type AdminComplaintItem struct {
	model.Complaint
	CoachName   string `json:"coach_name"`
	BookingDate string `json:"booking_date"`
}

// AdminListComplaints 投诉列表
func (s *Service) AdminListComplaints(status string, page, pageSize int) ([]AdminComplaintItem, int64, error) {
	q := s.store.DB().Model(&model.Complaint{})
	if status != "" {
		q = q.Where("complaints.status = ?", status)
	}
	var total int64
	_ = q.Count(&total).Error
	var list []AdminComplaintItem
	err := q.Select("complaints.*, coaches.nickname as coach_name, bookings.play_date as booking_date").
		Joins("JOIN bookings ON bookings.id = complaints.booking_id").
		Joins("JOIN coaches ON coaches.id = bookings.coach_id").
		Order("complaints.created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Scan(&list).Error
	return list, total, err
}

type AdminComplaintDetail struct {
	Complaint model.Complaint `json:"complaint"`
	Booking   model.Booking   `json:"booking"`
	Coach     model.Coach     `json:"coach"`
}

func (s *Service) AdminComplaintDetail(id int64) (*AdminComplaintDetail, error) {
	var cp model.Complaint
	if err := s.store.DB().First(&cp, id).Error; err != nil {
		return nil, errors.New("投诉不存在")
	}
	var b model.Booking
	_ = s.store.DB().First(&b, cp.BookingID).Error
	var c model.Coach
	_ = s.store.DB().First(&c, b.CoachID).Error
	return &AdminComplaintDetail{Complaint: cp, Booking: b, Coach: c}, nil
}

// AdminJudgeComplaint 裁决投诉。成立则记污点：2 次下架、3 次封号，自动执行。
func (s *Service) AdminJudgeComplaint(adminOpenID string, id int64, result, note string) error {
	if result != "upheld" && result != "rejected" {
		return errors.New("裁决结果不合法")
	}
	return s.store.DB().Transaction(func(tx *gorm.DB) error {
		var cp model.Complaint
		if err := tx.First(&cp, id).Error; err != nil {
			return errors.New("投诉不存在")
		}
		if cp.Status != "open" {
			return errors.New("已经处理过了")
		}
		now := time.Now()
		if err := tx.Model(&cp).Updates(map[string]any{"status": result, "handled_at": now}).Error; err != nil {
			return err
		}
		if err := s.logAdmin(tx, adminOpenID, "judge_complaint_"+result, "complaint", id, note); err != nil {
			return err
		}
		if result == "rejected" {
			return nil
		}
		// 成立：污点 +1
		var b model.Booking
		if err := tx.First(&b, cp.BookingID).Error; err != nil {
			return err
		}
		var c model.Coach
		if err := tx.First(&c, b.CoachID).Error; err != nil {
			return err
		}
		demerits := c.Demerits + 1
		patch := map[string]any{"demerits": demerits, "updated_at": now}
		autoNote := ""
		if demerits >= 3 {
			patch["status"] = "banned"
			patch["published"] = false
			autoNote = "累计 3 次，自动封号"
		} else if demerits == 2 {
			patch["status"] = "suspended"
			patch["published"] = false
			autoNote = "累计 2 次，自动下架"
		}
		if err := tx.Model(&model.Coach{}).Where("id = ?", c.ID).Updates(patch).Error; err != nil {
			return err
		}
		if autoNote != "" {
			return s.logAdmin(tx, adminOpenID, "auto_"+patch["status"].(string), "coach", c.ID, autoNote)
		}
		return nil
	})
}

type AdminStats struct {
	CoachesTotal         int64            `json:"coaches_total"`
	CoachesPublished     int64            `json:"coaches_published"`
	CoachesPendingVerify int64            `json:"coaches_pending_verify"`
	BookingsTotal        int64            `json:"bookings_total"`
	BookingsToday        int64            `json:"bookings_today"`
	BookingsByStatus     map[string]int64 `json:"bookings_by_status"`
	GMVTotal             int64            `json:"gmv_total"` // 分
	GMVMonth             int64            `json:"gmv_month"`
	ComplaintsOpen       int64            `json:"complaints_open"`
}

func (s *Service) AdminStats() (*AdminStats, error) {
	db := s.store.DB()
	st := &AdminStats{BookingsByStatus: map[string]int64{}}
	today := time.Now().Format("2006-01-02")
	month := time.Now().Format("2006-01")

	_ = db.Model(&model.Coach{}).Count(&st.CoachesTotal).Error
	_ = db.Model(&model.Coach{}).Where("published = ?", true).Count(&st.CoachesPublished).Error
	_ = db.Model(&model.Coach{}).Where("published = ? AND verified = ?", true, false).Count(&st.CoachesPendingVerify).Error
	_ = db.Model(&model.Booking{}).Count(&st.BookingsTotal).Error
	_ = db.Model(&model.Booking{}).Where("play_date = ?", today).Count(&st.BookingsToday).Error

	var rows []struct {
		Status string
		Count  int64
	}
	_ = db.Model(&model.Booking{}).Select("status, COUNT(*) as count").Group("status").Scan(&rows).Error
	for _, r := range rows {
		st.BookingsByStatus[r.Status] = r.Count
	}

	_ = db.Model(&model.Booking{}).Where("pay_status = 'paid'").
		Select("COALESCE(SUM(price),0)").Scan(&st.GMVTotal).Error
	_ = db.Model(&model.Booking{}).Where("pay_status = 'paid' AND to_char(play_date,'YYYY-MM') = ?", month).
		Select("COALESCE(SUM(price),0)").Scan(&st.GMVMonth).Error
	_ = db.Model(&model.Complaint{}).Where("status = 'open'").Count(&st.ComplaintsOpen).Error

	return st, nil
}
