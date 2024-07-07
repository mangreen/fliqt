package svc

import (
	"context"
	"fliqt/pkg/model"
	"fliqt/pkg/repo"
)

type (
	UserService interface {
		List(ctx context.Context, page int, size int) ([]model.User, error)
		FindByID(ctx context.Context, userID string) (*model.User, error)
		Create(ctx context.Context, userRequest model.User) (*model.User, error)
		DeleteByID(ctx context.Context, userID string) error
		Update(ctx context.Context, userID string, userRequest model.User) (*model.User, error)
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

func (svc *userSvc) List(ctx context.Context, page int, size int) ([]model.User, error) {
	users, err := svc.userRepo.List(ctx, page, size)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (svc *userSvc) FindByID(ctx context.Context, userID string) (*model.User, error) {
	user, err := svc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *userSvc) Create(ctx context.Context, userRequest model.User) (*model.User, error) {
	err := svc.userRepo.Create(ctx, userRequest)
	if err != nil {
		return nil, err
	}

	user, err := svc.userRepo.FindByID(ctx, userRequest.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *userSvc) DeleteByID(ctx context.Context, userID string) error {
	err := svc.userRepo.DeleteByID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (svc *userSvc) Update(ctx context.Context, userID string, userRequest model.User) (*model.User, error) {
	err := svc.userRepo.Update(ctx, userID, userRequest)
	if err != nil {
		return nil, err
	}

	user, err := svc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
