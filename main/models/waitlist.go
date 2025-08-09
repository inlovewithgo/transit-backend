package models

import (
	"time"

	"gorm.io/gorm"
)

type Waitlist struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	Email           string         `json:"email" gorm:"unique;not null;index"`
	ResendContactID string         `json:"resend_contact_id,omitempty" gorm:"column:resend_contact_id"`
	Status          string         `json:"status" gorm:"default:pending"`
	IPAddress       string         `json:"ip_address,omitempty" gorm:"column:ip_address"`
	UserAgent       string         `json:"user_agent,omitempty" gorm:"column:user_agent"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type WaitlistRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type WaitlistResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Email   string `json:"email,omitempty"`
}