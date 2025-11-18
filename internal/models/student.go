package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Student struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Age       int            `json:"age"`
	Major     string         `json:"major"`
	GPA       float64        `json:"gpa"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (s *Student) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
