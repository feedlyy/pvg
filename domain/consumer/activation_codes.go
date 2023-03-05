package consumer

import (
	"context"
	"os"
	"pvg/domain"
	"time"
)

type ActivationCodes struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint
	User      domain.Users `gorm:"references:ID"`
	Code      uint
	CreatedAt time.Time `gorm:"type:timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`
	ExpiresAt time.Time `gorm:"type:timestamp"`
}

// ACRepository Represent the Activation Code's repository contract
type ACRepository interface {
	Insert(ctx context.Context, ac ActivationCodes) error
}

// ACService Represent the Activation Code's service contract
type ACService interface {
	Process(topics []string, signals chan os.Signal)
}
