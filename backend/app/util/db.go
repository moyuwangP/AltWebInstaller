package util

import (
	db2 "AltWebServer/app/model/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB = nil
)

func init() {
	d, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	db = d

	err = db.AutoMigrate(
		&db2.Package{},
		&db2.Installation{},
		&db2.Setting{},
	)
}

func DB() *gorm.DB {
	return db
}
