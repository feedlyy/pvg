package domain

import (
	"context"
	"time"
)

type Users struct {
	ID        uint      `gorm:"primarykey"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Phone     uint      `json:"phone"`
	Email     string    `json:"email"`
	Birthday  time.Time `json:"birthday"`
	Status    string    `json:"status"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRepository Represent the User's repository contract
type UserRepository interface {
	Fetch(ctx context.Context) ([]Users, error)
}

type UserService interface {
	GetAllUser(ctx context.Context) ([]Users, error)
}
