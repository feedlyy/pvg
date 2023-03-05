package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"pvg/domain"
	"pvg/helper"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(usrRepo domain.UserRepository) domain.UserService {
	return &userService{
		userRepo: usrRepo,
	}
}

func (u *userService) GetAllUser(ctx context.Context) ([]domain.Users, error) {
	return u.userRepo.Fetch(ctx)
}

func (u *userService) GetUser(ctx context.Context, username string) (domain.Users, error) {
	return u.userRepo.GetByUsername(ctx, username)
}

func (u *userService) CreateUser(ctx context.Context, usr domain.Users) error {
	var (
		err    error
		dataDb domain.Users
	)

	// check whether inserted user alr have same credentials like phone, username, email
	dataDb, err = u.userRepo.GetByUsrPhoneEmail(ctx, usr)
	if err != nil {
		return err
	}

	if dataDb != (domain.Users{}) {
		err = errors.New(helper.DataExists)
		logrus.Errorf("User - Service|err when create user, err:%v", err)
		return err
	}

	return u.userRepo.Insert(ctx, usr)
}

func (u *userService) UpdateUser(ctx context.Context, usr domain.Users) error {
	var err error

	_, err = u.userRepo.GetById(ctx, int(usr.ID))
	if err != nil {
		return err
	}

	return u.userRepo.Update(ctx, usr)
}
