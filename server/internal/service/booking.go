package service

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/foxprince/paidang/server/internal/model"
)

// 预约状态机：待确认 → 已确认 → 已完成；待确认/已确认 → 已取消
var bookingTransitions = map[string][]string{
	"pending":   {"confirmed", "cancelled"},
	"confirmed": {"completed", "cancelled"},
}

func canTransit(from, to string) bool {
	for _, t := range bookingTransitions[from] {
		if t == to {
			return true
		}
	}
	return false
}

// 可约时段模板：{"weekly": {"1": [["09:00","10:00"]], ...}}，键为周几（0=周日）
type weeklySchedule struct {
	Weekly map[string][][2]string `json:"weekly"`
}

func parseSchedule(raw string) weeklySchedule {
	var ws weeklySchedule
	_ = json.Unmarshal([]byte(raw), &ws)
	if ws.Weekly == nil {
		ws.Weekly = map[string][][2]string{}
	}
	return ws
}

type Slot struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// AvailableSlots 某天可约时段 = 模板时段 − 已被占用的
func (s *Service) AvailableSlots(coachID int64, date string) ([]Slot, error) {
	var c model.Coach
	if err := s.store.DB().First(&c, coachID).Error; err != nil {
		return nil, err
	}
	if !c.Published {
		return nil, errors.New("该陪练主页未发布")
	}
	ws := parseSchedule(c.Schedule)
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("日期格式错误")
	}
	tmpl := ws.Weekly[string(rune('0'+int(t.Weekday())))]
	if len(tmpl) == 0 {
		return []Slot{}, nil
	}

	var booked []model.Booking
	_ = s.store.DB().
		Where("coach_id = ? AND play_date = ? AND status IN ('pending','confirmed')", coachID, date).
		Find(&booked).Error

	free := []Slot{}
	for _, sl := range tmpl {
		occupied := false
		for _, b := range booked {
			if sl[0] < b.EndTime && sl[1] > b.StartTime {
				occupied = true
				break
			}
		}
		if !occupied {
			free = append(free, Slot{Start: sl[0], End: sl[1]})
		}
	}
	return free, nil
}

// CreateBookingInput 学员下单
type CreateBookingInput struct {
	CoachID       int64
	Date          string
	Start, End    string
	Name, Phone   string
	StudentOpenID *string
	IdemKey       string
	Price         int
}

func (s *Service) CreateBooking(in CreateBookingInput) (*model.Booking, error) {
	// 幂等：同样的 idem_key 直接返回已有订单
	var existing model.Booking
	if err := s.store.DB().Where("idem_key = ?", in.IdemKey).First(&existing).Error; err == nil {
		return &existing, nil
	}

	if in.Start >= in.End {
		return nil, errors.New("时间不合法")
	}
	if in.Date < time.Now().Format("2006-01-02") {
		return nil, errors.New("不能预约过去的时间")
	}

	var c model.Coach
	if err := s.store.DB().First(&c, in.CoachID).Error; err != nil {
		return nil, errors.New("陪练不存在")
	}
	if !c.Published {
		return nil, errors.New("该陪练主页未发布")
	}

	// 时段必须落在模板内，且未被占用
	free, err := s.AvailableSlots(in.CoachID, in.Date)
	if err != nil {
		return nil, err
	}
	ok := false
	for _, sl := range free {
		if sl.Start == in.Start && sl.End == in.End {
			ok = true
			break
		}
	}
	if !ok {
		return nil, errors.New("该时段不可约")
	}

	// 老学员自动关联（按手机号）
	var studentID *int64
	var st model.Student
	if err := s.store.DB().Where("coach_id = ? AND phone = ?", in.CoachID, in.Phone).
		First(&st).Error; err == nil {
		studentID = &st.ID
	} else {
		st = model.Student{CoachID: in.CoachID, Name: in.Name, Phone: in.Phone}
		if err := s.store.DB().Create(&st).Error; err == nil {
			studentID = &st.ID
		}
	}

	b := model.Booking{
		CoachID: in.CoachID, StudentID: studentID,
		StudentName: in.Name, StudentPhone: in.Phone, StudentOpenID: in.StudentOpenID,
		PlayDate: in.Date, StartTime: in.Start, EndTime: in.End,
		Price: in.Price, Status: "pending", PayStatus: "unpaid", IdemKey: in.IdemKey,
	}
	if err := s.store.DB().Create(&b).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

// ListBookings 预约列表
func (s *Service) ListBookings(coachID int64, status, date string, page, pageSize int) ([]model.Booking, int64, error) {
	q := s.store.DB().Where("coach_id = ?", coachID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if date != "" {
		q = q.Where("play_date = ?", date)
	}
	var total int64
	_ = q.Model(&model.Booking{}).Count(&total).Error
	var list []model.Booking
	err := q.Order("play_date DESC, start_time DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// Transition 状态流转（确认/拒绝/完成/取消）
func (s *Service) Transition(id, coachID int64, to, reason, noShowBy string) (*model.Booking, error) {
	var b model.Booking
	if err := s.store.DB().Where("id = ? AND coach_id = ?", id, coachID).First(&b).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}
	if !canTransit(b.Status, to) {
		return nil, errors.New("当前状态不允许这个操作")
	}
	patch := map[string]any{"status": to, "updated_at": time.Now()}
	if to == "cancelled" {
		patch["cancel_reason"] = reason
		patch["no_show_by"] = noShowBy // coach / student
	}
	if err := s.store.DB().Model(&b).Updates(patch).Error; err != nil {
		return nil, err
	}
	b.Status = to
	return &b, nil
}

// MarkPaid 确认收款 → 流水入账
func (s *Service) MarkPaid(id, coachID int64) (*model.Booking, error) {
	var b model.Booking
	if err := s.store.DB().Where("id = ? AND coach_id = ?", id, coachID).First(&b).Error; err != nil {
		return nil, errors.New("订单不存在")
	}
	if b.Status != "confirmed" && b.Status != "completed" {
		return nil, errors.New("只有已确认的订单才能确认收款")
	}
	if err := s.store.DB().Model(&b).
		Updates(map[string]any{"pay_status": "paid", "updated_at": time.Now()}).Error; err != nil {
		return nil, err
	}
	b.PayStatus = "paid"
	return &b, nil
}

type StudentIncome struct {
	Name  string `json:"name"`
	Total int64  `json:"total"`
	Count int64  `json:"count"`
}

type IncomeSummary struct {
	Total     int64           `json:"total"` // 分
	Orders    int64           `json:"orders"`
	Hours     float64         `json:"hours"`
	ByStudent []StudentIncome `json:"by_student"`
}

// IncomeSummary 收入看板：只统计已确认收款的订单
func (s *Service) IncomeSummary(coachID int64, month string) (*IncomeSummary, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	var list []model.Booking
	if err := s.store.DB().
		Where("coach_id = ? AND pay_status = 'paid' AND to_char(play_date,'YYYY-MM') = ?", coachID, month).
		Find(&list).Error; err != nil {
		return nil, err
	}
	sum := &IncomeSummary{}
	byStudent := map[string]*StudentIncome{}
	for _, b := range list {
		sum.Total += int64(b.Price)
		sum.Orders++
		st, _ := time.Parse("15:04:05", b.StartTime)
		et, _ := time.Parse("15:04:05", b.EndTime)
		if et.After(st) {
			sum.Hours += et.Sub(st).Hours()
		}
		si, ok := byStudent[b.StudentName]
		if !ok {
			si = &StudentIncome{Name: b.StudentName}
			byStudent[b.StudentName] = si
		}
		si.Total += int64(b.Price)
		si.Count++
	}
	for _, si := range byStudent {
		sum.ByStudent = append(sum.ByStudent, *si)
	}
	return sum, nil
}
