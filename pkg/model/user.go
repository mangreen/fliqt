package model

import (
	"errors"
	"time"

	"github.com/charmbracelet/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID         string    `gorm:"primary_key;size:64; " json:"id" example:"james_bond"`
	Name       string    `gorm:"not null;size:64;" json:"name" example:"James Bond"`
	Password   string    `gorm:"not null;size:100;" json:"-"`
	Role       string    `gorm:"not null;size:20;" json:"role" example:"employee"`
	Department string    `gorm:"not null;size:20;" json:"department" example:"AI"`
	ManagerID  string    `gorm:"size:64; " json:"manager_id" example:"miss_m"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"upated_at"`
}

// For gorm
func (*User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if len(u.ID) == 0 {
		log.Error("[model] user.ID is invalid")
		return errors.New("[model] user.ID is invalid")
	}

	if len(u.Password) < 8 {
		log.Error("[model] user.Password is invalid")
		return errors.New("[model] user.Password is invalid")
	}

	if len(u.Role) == 0 {
		log.Error("[model] user.Role is invalid")
		return errors.New("[model] user.Role is invalid")
	}

	if len(u.Department) == 0 {
		log.Error("[model] user.Department is invalid")
		return errors.New("[model] user.Department is invalid")
	}

	hashpw, err := HashPassword(u.Password)
	if err != nil {
		log.Error("[model] hash password failed")
		return err
	}
	u.Password = hashpw

	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if u.ID == "admin" {
		log.Error("[model] admin user not allowed to delete")
		return errors.New("[model] admin not allowed to delete")
	}

	return
}
