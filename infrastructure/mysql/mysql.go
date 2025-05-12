package mysql

import (
	"clean-architecture-gochat/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

var (
	db   *gorm.DB
	once sync.Once
)

func Connect() *gorm.DB {
	once.Do(func() {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		)

		var err error
		// **修正: 載入設定檔*
		config.LoadConfig("internal/config") // 確保 `Config` 變數已初始化
		db, err = gorm.Open(mysql.Open(config.Config.MySQL.DNS), &gorm.Config{Logger: newLogger})
		if err != nil {
			log.Fatalf("MySQL connection failed: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("MySQL SQL DB init failed: %v", err)
		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Minute * 5)

		log.Println("MySQL connected successfully.")
	})
	return db
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database not initialized")
	}
	return db
}
