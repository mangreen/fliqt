package repo

import (
	"context"
	"database/sql"
	"fliqt/pkg/model"
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

func TestUserList(t *testing.T) {
	sqlDB, gormDB, mock := initMockDB(t)
	defer sqlDB.Close()

	mockUser := model.User{
		ID:         "kevin_chen",
		Name:       "Kevin Chen",
		Password:   "pw123456",
		Role:       "employee",
		Department: "AI",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "password", "role", "department", "created_at", "updated_at"}).
		AddRow(mockUser.ID, mockUser.Name, mockUser.Password, mockUser.Role, mockUser.Department, mockUser.CreatedAt, mockUser.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	userRepo := NewUserRepository(gormDB, gormDB)

	if _, err := userRepo.List(context.TODO(), 0, 0); err != nil {
		t.Errorf("error while running: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestUserFindByID(t *testing.T) {
	sqlDB, gormDB, mock := initMockDB(t)
	defer sqlDB.Close()

	mockUser := model.User{
		ID:         "kevin_chen",
		Name:       "Kevin Chen",
		Password:   "pw123456",
		Role:       "employee",
		Department: "AI",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "password", "role", "department", "created_at", "updated_at"}).
		AddRow(mockUser.ID, mockUser.Name, mockUser.Password, mockUser.Role, mockUser.Department, mockUser.CreatedAt, mockUser.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	userRepo := NewUserRepository(gormDB, gormDB)

	if _, err := userRepo.FindByID(context.TODO(), mockUser.ID); err != nil {
		t.Errorf("error while running: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestUserCreate(t *testing.T) {
	sqlDB, gormDB, mock := initMockDB(t)
	defer sqlDB.Close()

	mockUser := model.User{
		ID:         "kevin_chen",
		Name:       "Kevin Chen",
		Password:   "pw123456",
		Role:       "employee",
		Department: "AI",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WithArgs(mockUser.ID, mockUser.Name, sqlmock.AnyArg(), mockUser.Role, mockUser.Department, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	userRepo := NewUserRepository(gormDB, gormDB)

	if err := userRepo.Create(context.TODO(), mockUser); err != nil {
		t.Errorf("error while running: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestUserUpdate(t *testing.T) {
	sqlDB, gormDB, mock := initMockDB(t)
	defer sqlDB.Close()

	mockUserID := "kevin_chen"
	mockUser := model.User{
		Name:       "Kevin Chen",
		Role:       "employee",
		Department: "AI",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `name`=?,`role`=?,`department`=?,`updated_at`=? WHERE `id` = ?")).
		WithArgs(mockUser.Name, mockUser.Role, mockUser.Department, sqlmock.AnyArg(), mockUserID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	userRepo := NewUserRepository(gormDB, gormDB)

	if err := userRepo.Update(context.TODO(), mockUserID, mockUser); err != nil {
		t.Errorf("error while running: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestUserDeleteByID(t *testing.T) {
	sqlDB, gormDB, mock := initMockDB(t)
	defer sqlDB.Close()

	mockUserID := "kevin_chen"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `users`")).
		WithArgs(mockUserID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	userRepo := NewUserRepository(gormDB, gormDB)

	if err := userRepo.DeleteByID(context.TODO(), mockUserID); err != nil {
		t.Errorf("error while running: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}
