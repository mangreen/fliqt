package svc

import (
	"context"
	"errors"
	"fliqt/pkg/http/api"
	"fliqt/pkg/model"
	"fliqt/pkg/repo"
	"time"

	"github.com/charmbracelet/log"
	"github.com/dgrijalva/jwt-go"
)

type (
	AuthService interface {
		Validate(ctx context.Context, userName string, password string) (*api.LoginResponse, error)
	}

	authSvc struct {
		userRepo repo.UserRepository
	}

	JwtClaims struct {
		*model.User
		jwt.StandardClaims
	}
)

const FLIQT_CONST = "FliQt"

func NewAuthService(userRepo repo.UserRepository) AuthService {
	return &authSvc{
		userRepo: userRepo,
	}
}

func (svc *authSvc) Validate(ctx context.Context, userName string, password string) (*api.LoginResponse, error) {
	usr, err := svc.userRepo.FindByName(ctx, userName)
	if err != nil {
		return nil, err
	}

	if !model.CheckPasswordHash(password, usr.Password) {
		log.Error("[svc] validate password failed")
		return nil, errors.New("[svc] validate password failed")
	}

	token, err := GenToken(usr)
	if err != nil {
		log.Errorf("[svc] generate token failed: %v", err)
		return nil, err
	}
	res := &api.LoginResponse{
		User:  usr,
		Token: token,
	}

	return res, nil
}

func GenToken(user *model.User) (string, error) {
	claims := JwtClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Issuer:    FLIQT_CONST,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(FLIQT_CONST))
}
