package repository

import (
	"context"
	"github.com/sirupsen/logrus"
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
		err error
	)

	err = u.db.WithContext(ctx).Find(&res).Error
	if err != nil {
		logrus.Errorf("User - Repository|err when fetch all users, err:%v", err)
		return nil, err
	}

	return res, nil
}

func (u *userRepository) GetByUsername(ctx context.Context, username string) (domain.Users, error) {
	var (
		res domain.Users
		err error
	)

	err = u.db.WithContext(ctx).First(&res, "username = ?", username).Error
	if err != nil {
		logrus.Errorf("User - Repository|err when get user by username, err:%v", err)
		return domain.Users{}, err
	}

	return res, nil
}
