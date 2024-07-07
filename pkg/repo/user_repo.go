package repo

import (
	"context"
	"fliqt/pkg/common/mysql"
	"fliqt/pkg/model"

	"github.com/charmbracelet/log"
	"gorm.io/gorm"
)

type (
	UserRepository interface {
		List(ctx context.Context, page int, size int) ([]model.User, error)
		FindByID(ctx context.Context, userID string) (*model.User, error)
		Create(ctx context.Context, userRequest model.User) error
		DeleteByID(ctx context.Context, userID string) error
		Update(ctx context.Context, userID string, userRequest model.User) error
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

func (repo *userRepo) List(ctx context.Context, page int, size int) ([]model.User, error) {

	users := []model.User{}
	err := repo.readDB.Scopes(mysql.Paginate(page, size)).Find(&users).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}

		log.Errorf("[repo] list user failed: %v", err)
		return nil, err
	}
	return users, nil
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

func (repo *userRepo) Create(ctx context.Context, userRequest model.User) error {

	err := repo.readDB.Create(&userRequest).Error
	if err != nil {
		log.Errorf("[repo] create user:%v failed: %v", userRequest, err)
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

func (repo *userRepo) Update(ctx context.Context, userID string, userRequest model.User) error {

	err := repo.readDB.Model(&model.User{
		ID: userID,
	}).Updates(&userRequest).Error
	if err != nil {
		log.Errorf("[repo] update user:%v failed: %v", userRequest, err)
		return err
	}
	return nil
}
