package svc

import (
	"context"
	"fliqt/pkg/model"
	"fliqt/pkg/repo"
)

type (
	UserService interface {
		FindByID(ctx context.Context, userID string) (*model.User, error)
		Create(ctx context.Context, user *model.User) (*model.User, error)
		DeleteByID(ctx context.Context, userID string) error
		Update(ctx context.Context, userID string, userUpdate *model.User) (*model.User, error)
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

func (svc *userSvc) Create(ctx context.Context, user *model.User) (*model.User, error) {
	err := svc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	usr, err := svc.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (svc *userSvc) DeleteByID(ctx context.Context, userID string) error {
	err := svc.userRepo.DeleteByID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (svc *userSvc) Update(ctx context.Context, userID string, userUpdate *model.User) (*model.User, error) {
	err := svc.userRepo.Update(ctx, userID, userUpdate)
	if err != nil {
		return nil, err
	}

	usr, err := svc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
