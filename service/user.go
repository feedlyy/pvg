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
	kafka    domain.KafkaProducer
}

func NewUserService(usrRepo domain.UserRepository, k domain.KafkaProducer) domain.UserService {
	return &userService{
		userRepo: usrRepo,
		kafka:    k,
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
		err        error
		dataDb     domain.Users
		insertedId int
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

	insertedId, err = u.userRepo.Insert(ctx, usr)
	if err != nil {
		return err
	}

	err = u.kafka.SendMessage(helper.EmailTopic, insertedId)
	if err != nil {
		logrus.Errorf("User - Service|err send email user with kafka, err:%v", err)
		return err
	}

	return nil
}

func (u *userService) UpdateUser(ctx context.Context, usr domain.Users) error {
	var err error

	_, err = u.userRepo.GetById(ctx, int(usr.ID))
	if err != nil {
		return err
	}

	return u.userRepo.Update(ctx, usr)
}

func (u *userService) DeleteUser(ctx context.Context, id int) error {
	var err error

	_, err = u.userRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	return u.userRepo.Delete(ctx, id)
}
