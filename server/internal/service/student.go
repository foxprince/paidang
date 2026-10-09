package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/foxprince/paidang/server/internal/model"
)

// ListStudents 学员列表，按最近上课时间倒序太贵，MVP 按创建时间倒序
func (s *Service) ListStudents(coachID int64, page, pageSize int) ([]model.Student, int64, error) {
	q := s.store.DB().Where("coach_id = ?", coachID)
	var total int64
	_ = q.Model(&model.Student{}).Count(&total).Error
	var list []model.Student
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (s *Service) CreateStudent(coachID int64, name, phone, wechat, level string) (*model.Student, error) {
	if name == "" {
		return nil, errors.New("学员姓名必填")
	}
	st := model.Student{CoachID: coachID, Name: name, Phone: phone, Wechat: wechat, Level: level}
	if err := s.store.DB().Create(&st).Error; err != nil {
		return nil, err
	}
	return &st, nil
}

type StudentDetail struct {
	Student model.Student      `json:"student"`
	Notes   []model.LessonNote `json:"notes"` // 按次课备注时间线
}

func (s *Service) GetStudentDetail(coachID, id int64) (*StudentDetail, error) {
	var st model.Student
	if err := s.store.DB().Where("id = ? AND coach_id = ?", id, coachID).First(&st).Error; err != nil {
		return nil, errors.New("学员不存在")
	}
	var notes []model.LessonNote
	_ = s.store.DB().Where("student_id = ?", id).Order("created_at DESC").Find(&notes).Error
	return &StudentDetail{Student: st, Notes: notes}, nil
}

// AddNote 写按次课备注：本节要点 + 下节计划
func (s *Service) AddNote(coachID, studentID int64, bookingID *int64, keyPoints, nextPlan string) (*model.LessonNote, error) {
	if keyPoints == "" {
		return nil, errors.New("本节要点不能为空")
	}
	var st model.Student
	if err := s.store.DB().Where("id = ? AND coach_id = ?", studentID, coachID).First(&st).Error; err != nil {
		return nil, errors.New("学员不存在")
	}
	note := model.LessonNote{
		CoachID: coachID, StudentID: studentID,
		KeyPoints: keyPoints, NextPlan: nextPlan,
	}
	if bookingID != nil {
		note.BookingID = *bookingID
	}
	if err := s.store.DB().Create(&note).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

// CreateReview 学员评价：仅已完成订单可评价，一单一评
func (s *Service) CreateReview(bookingID int64, rating int, comment string) (*model.Review, error) {
	if rating < 1 || rating > 5 {
		return nil, errors.New("评分 1-5")
	}
	var b model.Booking
	if err := s.store.DB().First(&b, bookingID).Error; err != nil {
		return nil, errors.New("订单不存在")
	}
	if b.Status != "completed" {
		return nil, errors.New("课程完成后才能评价")
	}
	rv := model.Review{BookingID: bookingID, CoachID: b.CoachID, Rating: int16(rating), Comment: comment}
	if err := s.store.DB().Create(&rv).Error; err != nil {
		return nil, errors.New("已经评价过了")
	}
	return &rv, nil
}

// CreateComplaint 投诉：付款截图 + 聊天记录
func (s *Service) CreateComplaint(bookingID int64, evidence []string, detail string) (*model.Complaint, error) {
	var b model.Booking
	if err := s.store.DB().First(&b, bookingID).Error; err != nil {
		return nil, errors.New("订单不存在")
	}
	evBytes, _ := json.Marshal(evidence)
	if len(evidence) == 0 {
		evBytes = []byte("[]")
	}
	cp := model.Complaint{BookingID: bookingID, Evidence: string(evBytes), Detail: detail}
	if err := s.store.DB().Create(&cp).Error; err != nil {
		return nil, err
	}
	return &cp, nil
}

// Touch 更新 updated_at 小工具
func (s *Service) touch(table string, id int64) {
	_ = s.store.DB().Table(table).Where("id = ?", id).
		Update("updated_at", time.Now()).Error
}
