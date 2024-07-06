package model

import (
	"errors"
	"time"

	"github.com/charmbracelet/log"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"primary_key;" json:"id" example:"10"`
	Name      string `gorm:"not null;unique" json:"name"`
	Password  string `gorm:"not null;" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// For gorm
func (*User) TableName() string {
	return "users"
}

func (mdl *User) BeforeCreate(tx *gorm.DB) error {
	uuid, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 6)
	if err != nil {
		log.Error("[model] gonanoid failed")
		return err
	}

	mdl.ID = "u" + time.Now().UTC().Format("20060102150405") + uuid

	if len(mdl.Name) <= 2 {
		return errors.New("[model] user.Name length is less than 2")
	}

	if len(mdl.Password) < 8 {
		log.Error("[model] user.Password is invalid")
		return errors.New("[model] user.Password is invalid")
	}

	hashpw, err := HashPassword(mdl.Password)
	if err != nil {
		log.Error("[model] hash password failed")
		return err
	}
	mdl.Password = hashpw

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
