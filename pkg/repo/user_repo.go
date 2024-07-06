package repo

import (
	"context"
	"fliqt/pkg/model"

	"github.com/charmbracelet/log"
	"gorm.io/gorm"
)

type (
	UserRepository interface {
		FindByID(ctx context.Context, userID string) (*model.User, error)
		FindByName(ctx context.Context, userName string) (*model.User, error)
	}

	userRepo struct {
		readDB  *gorm.DB
		writeDB *gorm.DB
	}
)

func NewUserRepository(readDB *gorm.DB, writeDB *gorm.DB) UserRepository {
	return &userRepo{
		readDB:  readDB,
		writeDB: writeDB,
	}
}

func (repo *userRepo) FindByID(ctx context.Context, userID string) (*model.User, error) {

	user := &model.User{}
	err := repo.readDB.Where("id = ?", userID).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}

		log.Errorf("[repo] get user ID %v failed: %v", userID, err)
		return nil, err
	}
	return user, nil
}

func (repo *userRepo) FindByName(ctx context.Context, userName string) (*model.User, error) {

	user := &model.User{}
	err := repo.readDB.Where("name = ?", userName).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}

		log.Errorf("[repo] get user by name %v failed: %v", userName, err)
		return nil, err
	}
	return user, nil
}
