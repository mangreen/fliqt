package svc

import (
	"context"
	"fliqt/pkg/model"
	"fliqt/pkg/repo"
)

type (
	UserService interface {
		FindByID(ctx context.Context, userID string) (*model.User, error)
	}

	userSvc struct {
		userRepo repo.UserRepository
	}
)

func NewUserService(userRepo repo.UserRepository) UserService {
	return &userSvc{
		userRepo: userRepo,
	}
}

func (svc *userSvc) FindByID(ctx context.Context, userID string) (*model.User, error) {
	usr, err := svc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
