package service

import (
	"context"
	"pvg/domain"
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
