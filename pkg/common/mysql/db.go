package mysql

import (
	"errors"
	"fliqt/pkg/model"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabases() (*gorm.DB, *gorm.DB, error) {
	var err error

	var readDB, writeDB *gorm.DB

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.ip"),
		viper.GetString("mysql.dbname"))

	readDB, err = setup(dsn)
	if err != nil {
		return nil, nil, err
	}

	writeDB, err = setup(dsn)
	if err != nil {
		return nil, nil, err
	}

	migration(writeDB)

	return readDB, writeDB, nil
}

func setup(dsn string) (*gorm.DB, error) {
	log.Info("[mysql] database", "dsn", dsn)

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	var conn *gorm.DB

	var err error
	err = backoff.Retry(func() error {
		conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Errorf("[mysql] mysql open failed: %v", err)
			return err
		}

		sqlDB, err := conn.DB()
		if err != nil {
			log.Errorf("[mysql] get *sql.DB failed: %v", err)
			return err
		}
		if err = sqlDB.Ping(); err != nil {
			log.Errorf("[mysql] mysql ping error: %v", err)
			return err
		}
		log.Info("[mysql] database ping success")

		sqlDB.SetMaxIdleConns(150)
		sqlDB.SetMaxOpenConns(300)
		sqlDB.SetConnMaxLifetime(14400 * time.Second)

		return nil
	}, bo)

	if err != nil {
		log.Errorf("[mysql] mysql connect err: %v", err)
		return nil, err
	}

	return conn, nil
}

func migration(conn *gorm.DB) error {
	err := conn.AutoMigrate(&model.User{})
	if err != nil {
		log.Errorf("[mysql] aotumigration failed: %v", err)
		return err
	}

	if conn.Migrator().HasTable(&model.User{}) {
		if err := conn.First(&model.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			admin := model.User{
				ID:         "admin",
				Name:       "Admin",
				Password:   "admin123",
				Role:       "admin",
				Department: "admin",
			}

			if err := conn.Create(&admin).Error; err != nil {
				log.Errorf("[mysql] create admin failed: %v", err)
				return err
			}
		} else {
			log.Info("[mysql] user 'admin' existed")
		}
	}

	return nil
}
