package model

import "time"

type CaptchaLog struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Scene       string    `gorm:"size:100;not null;index" json:"scene"`
	IP          string    `gorm:"column:ip;size:50;not null" json:"ip"`
	UserAgent   string    `gorm:"column:user_agent;size:500" json:"user_agent"`
	Success     bool      `gorm:"not null" json:"success"`
	RequestID   string    `gorm:"column:request_id;size:100" json:"request_id"`
	VerifyParam string    `gorm:"column:verify_param;type:text" json:"verify_param"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (CaptchaLog) TableName() string {
	return "captcha_logs"
}
