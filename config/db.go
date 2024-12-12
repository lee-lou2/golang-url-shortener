package config

import (
	"gorm.io/gorm/logger"
	"log"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dbReadInstance  *gorm.DB
	dbReadOnce      sync.Once
	dbWriteInstance *gorm.DB
	dbWriteOnce     sync.Once
)

// GetReadDB 읽기 데이터베이스 인스턴스
func GetReadDB() *gorm.DB {
	dbReadOnce.Do(func() {
		var db *gorm.DB
		var err error
		dbFileName := "sqlite.db"
		db, err = gorm.Open(sqlite.Open(dbFileName), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("failed to get generic database object: %v", err)
		}

		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(time.Hour)

		dbReadInstance = db
	})
	return dbReadInstance
}

// GetWriteDB 쓰기 데이터베이스 인스턴스
func GetWriteDB() *gorm.DB {
	dbWriteOnce.Do(func() {
		var db *gorm.DB
		var err error
		dbFileName := "sqlite.db"
		db, err = gorm.Open(sqlite.Open(dbFileName), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("failed to get generic database object: %v", err)
		}

		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
		sqlDB.SetConnMaxLifetime(time.Hour)

		dbWriteInstance = db
	})
	return dbWriteInstance
}
