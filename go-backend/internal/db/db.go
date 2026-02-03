package db

import (
	"task-system-go/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg config.Config) error {
	db, err := gorm.Open(mysql.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}
