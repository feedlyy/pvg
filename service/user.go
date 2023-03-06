package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"pvg/domain"
	"pvg/helper"
	"time"
)

type userService struct {
	userRepo domain.UserRepository
	kafka    domain.KafkaProducer
	acRepo   domain.ActivationCodeRepository
}

func NewUserService(usrRepo domain.UserRepository, k domain.KafkaProducer,
	ac domain.ActivationCodeRepository) domain.UserService {
	return &userService{
		userRepo: usrRepo,
		kafka:    k,
		acRepo:   ac,
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
	usr.ID = uint(insertedId)

	err = u.kafka.SendMessage(helper.EmailTopic, usr)
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

func (u *userService) ActivateUser(ctx context.Context, username string, code int) error {
	var (
		err  error
		ac   domain.ActivationCodes
		user domain.Users
		now  = time.Now()
		diff time.Duration
	)

	// get detail user
	user, err = u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return err
	}

	// get activation codes data
	ac, err = u.acRepo.GetByUserId(ctx, int(user.ID))
	if err != nil {
		return err
	}
	diff = now.Sub(ac.ExpiresAt)

	switch {
	case code != int(ac.Code):
		err = errors.New(helper.NoValidCode)
		logrus.Errorf("User - Service|err when validate user activation, err:%v", err)
		return err
	case int(diff.Hours())+7 > 1: // +7 for convert utc to gmt
		err = errors.New(helper.NoValid)
		logrus.Errorf("User - Service|err when validate user activation codes, err:%v", err)
		return err
	case user.Status == helper.Active:
		err = errors.New(helper.DataActive)
		logrus.Errorf("User - Service|err when validate user activation codes, err:%v", err)
		return err
	default:
	}

	// set user to active if valid
	user.Status = helper.Active
	err = u.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *userService) RequestActivationCode(ctx context.Context, username string) (domain.ActivationCodes, error) {
	var (
		err            error
		ac             domain.ActivationCodes
		user           domain.Users
		now            = time.Now()
		diff           time.Duration
		activationCode = helper.GenerateActivationCodes()
		loc            *time.Location
		expiredAt      time.Time
	)
	loc, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		logrus.Errorf("User - Service|Err when get location %v", err)
		return domain.ActivationCodes{}, err
	}
	expiredAt = time.Now().In(loc)

	// get detail user
	user, err = u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return domain.ActivationCodes{}, err
	}

	// get activation codes data
	ac, err = u.acRepo.GetByUserId(ctx, int(user.ID))
	if err != nil {
		return domain.ActivationCodes{}, err
	}
	diff = now.Sub(ac.ExpiresAt)

	switch {
	case user.Status == helper.Active:
		err = errors.New(helper.DataActive)
		logrus.Errorf("User - Service|err when validate user activation codes, err:%v", err)
		return domain.ActivationCodes{}, err
	case int(diff.Hours())+7 < 2: // +7 for convert utc to gmt
		err = errors.New(helper.StillValid)
		logrus.Errorf("User - Service|err when validate user activation codes, err:%v", err)
		return ac, err
	default:
	}

	ac = domain.ActivationCodes{} // assign nil
	ac = domain.ActivationCodes{
		UserID:    user.ID,
		Code:      uint(activationCode),
		ExpiresAt: expiredAt.Add(time.Hour),
	}
	err = u.acRepo.Insert(ctx, ac)
	if err != nil {
		return domain.ActivationCodes{}, err
	}

	return ac, nil
}
