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

func (u *userRepository) GetByUsrPhoneEmail(ctx context.Context, user domain.Users) (domain.Users, error) {
	var (
		res domain.Users
		err error
	)

	err = u.db.WithContext(ctx).Raw("SELECT username, email, phone FROM users WHERE username = ? OR "+
		"email = ? OR phone = ?", user.Username, user.Email, user.Phone).Scan(&res).Error
	if err != nil {
		logrus.Errorf("User - Repository|err when get user by username, phone, or email, err:%v", err)
		return domain.Users{}, err
	}

	return res, nil
}

func (u *userRepository) Insert(ctx context.Context, usr domain.Users) (int, error) {
	var err error

	if err = u.db.WithContext(ctx).Create(&usr).Error; err != nil {
		logrus.Errorf("User - Repository|err when store user, err:%v", err)
		return 0, err
	}

	return int(usr.ID), nil
}

func (u *userRepository) Update(ctx context.Context, usr domain.Users) error {
	var err error

	if err = u.db.WithContext(ctx).Model(&usr).Updates(domain.Users{
		Password:  usr.Password,
		Firstname: usr.Firstname,
		Lastname:  usr.Lastname,
		Phone:     usr.Phone,
		Birthday:  usr.Birthday,
		Status:    usr.Status,
	}).Error; err != nil {
		logrus.Errorf("User - Repository|err when update user, err:%v", err)
		return err
	}

	return nil
}

func (u *userRepository) GetById(ctx context.Context, id int) (domain.Users, error) {
	var (
		res domain.Users
		err error
	)

	err = u.db.WithContext(ctx).First(&res, "id = ?", id).Error
	if err != nil {
		logrus.Errorf("User - Repository|err when get user by id, err:%v", err)
		return domain.Users{}, err
	}

	return res, nil
}

func (u *userRepository) Delete(ctx context.Context, id int) error {
	var err error

	err = u.db.WithContext(ctx).Delete(domain.Users{}, id).Error
	if err != nil {
		logrus.Errorf("User - Repository|err when delete user, err:%v", err)
		return err
	}

	return nil
}
