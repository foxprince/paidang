package model

import "time"

// 金额全部用"分"存整数
type Coach struct {
	ID           int64     `gorm:"primaryKey"`
	OpenID       string    `gorm:"uniqueIndex;size:64;not null"`
	Nickname     string    `gorm:"size:64;not null"`
	AvatarURL    string    `gorm:"size:512"`
	Phone        string    `gorm:"size:20"` // 展示时脱敏
	Title        string    `gorm:"size:64"`
	Bio          string    `gorm:"type:text"`
	HomeCourt    string    `gorm:"size:128"`
	VideoURL     string    `gorm:"size:512"` // 30秒对打视频，发布必填
	PricePerHour int       `gorm:"not null;default:0"`
	PayQRURL     string    `gorm:"size:512"` // 微信收款码
	Schedule     string    `gorm:"type:jsonb;not null;default:'{}'"`
	Verified     bool      `gorm:"not null;default:false"`
	Published    bool      `gorm:"not null;default:false"`
	Demerits     int       `gorm:"not null;default:0"`
	Status       string    `gorm:"size:16;not null;default:'active'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Student struct {
	ID        int64  `gorm:"primaryKey"`
	CoachID   int64  `gorm:"index;not null"`
	Coach     Coach  `gorm:"foreignKey:CoachID"`
	Name      string `gorm:"size:64;not null"`
	Phone     string `gorm:"size:20"`
	Wechat    string `gorm:"size:64"`
	Level     string `gorm:"size:32"`
	Remark    string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Booking struct {
	ID           int64  `gorm:"primaryKey"`
	CoachID      int64  `gorm:"index;not null"`
	Coach        Coach  `gorm:"foreignKey:CoachID"`
	StudentID    *int64 `gorm:"index"`
	Student      *Student
	StudentName  string `gorm:"size:64;not null"`
	StudentPhone string `gorm:"size:20;not null"`
	StudentOpenID *string `gorm:"size:64;index"` // 下单人 openid，用于订阅消息
	PlayDate     string `gorm:"type:date;not null"` // YYYY-MM-DD
	StartTime    string `gorm:"type:time;not null"`
	EndTime      string `gorm:"type:time;not null"`
	ServiceType  string `gorm:"size:16;not null;default:'single'"` // v2 预留
	Price        int    `gorm:"not null"`
	Status       string `gorm:"size:16;not null;default:'pending'"` // pending/confirmed/completed/cancelled
	PayStatus    string `gorm:"size:16;not null;default:'unpaid'"`   // unpaid/paid/reconciling
	NoShowBy     string `gorm:"size:16"`                             // coach/student
	CancelReason string `gorm:"size:256"`
	IdemKey      string `gorm:"uniqueIndex;size:64;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type LessonNote struct {
	ID        int64 `gorm:"primaryKey"`
	BookingID int64 `gorm:"not null"`
	CoachID   int64 `gorm:"not null"`
	StudentID int64 `gorm:"index;not null"`
	KeyPoints string `gorm:"type:text;not null"` // 本节要点
	NextPlan  string `gorm:"type:text"`          // 下节计划
	CreatedAt time.Time
}

type Review struct {
	ID        int64 `gorm:"primaryKey"`
	BookingID int64 `gorm:"uniqueIndex;not null"` // 仅已完成订单可评价
	CoachID   int64 `gorm:"index;not null"`
	Rating    int16 `gorm:"not null"`
	Comment   string `gorm:"type:text"`
	CreatedAt time.Time
}

type Complaint struct {
	ID        int64  `gorm:"primaryKey"`
	BookingID int64  `gorm:"not null"`
	Evidence  string `gorm:"type:jsonb;not null;default:'[]'"`
	Detail    string `gorm:"type:text"`
	Status    string `gorm:"size:16;not null;default:'open'"` // open/upheld/rejected
	HandledAt *time.Time
	CreatedAt time.Time
}
