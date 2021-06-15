package model

import (
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now()
		},
	})
	if err != nil {
		zap.S().Errorw("failed to open db", "error", err)
	}
	DB.AutoMigrate(&Client{})
}
