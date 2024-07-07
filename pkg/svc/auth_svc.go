package svc

import (
	"context"
	"errors"
	"fliqt/pkg/model"
	"fliqt/pkg/repo"

	"github.com/charmbracelet/log"
)

type (
	AuthService interface {
		Validate(ctx context.Context, userID string, password string) (*model.User, error)
	}

	authSvc struct {
		userRepo repo.UserRepository
	}
)

func NewAuthService(userRepo repo.UserRepository) AuthService {
	return &authSvc{
		userRepo: userRepo,
	}
}

func (svc *authSvc) Validate(ctx context.Context, userID string, password string) (*model.User, error) {
	usr, err := svc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !model.CheckPasswordHash(password, usr.Password) {
		log.Error("[svc] validate password failed")
		return nil, errors.New("[svc] validate password failed")
	}

	return usr, nil
}
