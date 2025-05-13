package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	StatusActive   int = 1  // active
	StatusInactive int = 0  // non-active
	StatusBanned   int = -1 // disable
)

type BaseModel struct {
	Id        uint64         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *BaseModel) GetCreatedTime() string {
	return m.CreatedAt.Format("2006-01-02 15:04:05")
}

func (m *BaseModel) GetUpdatedTime() string {
	return m.UpdatedAt.Format("2006-01-02 15:04:05")
}

func (m *BaseModel) IsDeleted() bool {
	return !m.DeletedAt.Time.IsZero()
}
