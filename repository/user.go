package repository

import (
	"context"
	"gorm.io/gorm"
	"pvg/domain"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Fetch(ctx context.Context) ([]domain.Users, error) {
	var (
		res []domain.Users
	)
	u.db.WithContext(ctx).Find(&res)

	return res, nil
}
