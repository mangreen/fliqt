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
		Create(ctx context.Context, usr *model.User) error
		DeleteByID(ctx context.Context, userID string) error
		Update(ctx context.Context, userID string, user *model.User) error
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

	usr := &model.User{}
	err := repo.readDB.Where("id = ?", userID).First(usr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}

		log.Errorf("[repo] get user ID %v failed: %v", userID, err)
		return nil, err
	}
	return usr, nil
}

func (repo *userRepo) Create(ctx context.Context, usr *model.User) error {

	err := repo.readDB.Create(usr).Error
	if err != nil {
		log.Errorf("[repo] create user:%v failed: %v", usr, err)
		return err
	}
	return nil
}

func (repo *userRepo) DeleteByID(ctx context.Context, userID string) error {

	err := repo.readDB.Delete(&model.User{
		ID: userID,
	}).Error
	if err != nil {
		log.Errorf("[repo] delete user ID %v failed: %v", userID, err)
		return err
	}
	return nil
}

func (repo *userRepo) Update(ctx context.Context, userID string, userUpdate *model.User) error {

	err := repo.readDB.Model(&model.User{
		ID: userID,
	}).Updates(userUpdate).Error
	if err != nil {
		log.Errorf("[repo] update user:%v failed: %v", userUpdate, err)
		return err
	}
	return nil
}
