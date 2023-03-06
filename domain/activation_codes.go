package domain

import (
	"context"
	"time"
)

type ActivationCodes struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint
	Code      uint      `form:"code" binding:"required"`
	CreatedAt time.Time `gorm:"type:timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`
	ExpiresAt time.Time `gorm:"type:timestamp"`
}

type ActivationCodeRepository interface {
	GetByUserId(ctx context.Context, id int) (ActivationCodes, error)
	Insert(ctx context.Context, ac ActivationCodes) error
}
