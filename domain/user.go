package domain

import (
	"context"
	"time"
)

type Users struct {
	ID        uint      `gorm:"primarykey"`
	Username  string    `json:"username" form:"username" binding:"required"`
	Password  string    `json:"password" form:"password" binding:"required"`
	Firstname string    `json:"firstname" form:"firstname" binding:"required"`
	Lastname  string    `json:"lastname" form:"lastname" binding:"required"`
	Phone     uint      `json:"phone" form:"phone" binding:"required"`
	Email     string    `json:"email" form:"email" binding:"required"`
	Birthday  time.Time `json:"birthday" gorm:"type:date"`
	Status    string    `json:"status"`
	CreatedAt time.Time `gorm:"type:timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`
}

// UserRepository Represent the User's repository contract
type UserRepository interface {
	Fetch(ctx context.Context) ([]Users, error)
	GetByUsername(ctx context.Context, username string) (Users, error)
	GetByUsrPhoneEmail(ctx context.Context, user Users) (Users, error)
	GetById(ctx context.Context, id int) (Users, error)
	Insert(ctx context.Context, usr Users) (int, error)
	Update(ctx context.Context, usr Users) error
	Delete(ctx context.Context, id int) error
}

type UserService interface {
	GetAllUser(ctx context.Context) ([]Users, error)
	GetUser(ctx context.Context, username string) (Users, error)
	CreateUser(ctx context.Context, usr Users) error
	UpdateUser(ctx context.Context, usr Users) error
	DeleteUser(ctx context.Context, id int) error
	ActivateUser(ctx context.Context, username string, code int) error
	RequestActivationCode(ctx context.Context, username string) (ActivationCodes, error)
}
