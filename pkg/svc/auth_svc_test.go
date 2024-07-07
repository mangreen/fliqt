package svc

import (
	"context"
	"database/sql"
	"fliqt/pkg/model"
	"fliqt/pkg/repo"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initMockDB(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	return sqlDB, gormDB, mock
}

func TestAuthValidate(t *testing.T) {
	sqlDB, gormDB, mock := initMockDB(t)
	defer sqlDB.Close()

	mockUser := model.User{
		ID:         "kevin_chen",
		Password:   "pw123456",
		Role:       "employee",
		Department: "AI",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	hashPw, _ := model.HashPassword(mockUser.Password)

	rows := sqlmock.NewRows([]string{"id", "name", "password", "role", "department", "created_at", "updated_at"}).
		AddRow(mockUser.ID, mockUser.Name, hashPw, mockUser.Role, mockUser.Department, mockUser.CreatedAt, mockUser.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	userRepo := repo.NewUserRepository(gormDB, gormDB)
	authSvc := NewAuthService(userRepo)

	if _, err := authSvc.Validate(context.TODO(), mockUser.ID, mockUser.Password); err != nil {
		t.Errorf("error while running: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}
