/**
* @Auth:ShenZ
* @Description: 資料庫遷移工具
* @CreateDate:2022/06/15 10:57:44
 */
package main

import (
	"clean-architecture-gochat/internal/domain/entities"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 從環境變數獲取資料庫連接資訊
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "mariadb"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "rootpassword"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "newgochat"
	}

	// 構建資料庫連接字串
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// 刪除現有的表
	err = db.Migrator().DropTable(
		&entities.Group{},
		&entities.GroupMember{},
		&entities.User{},
		&entities.Message{},
		&entities.Contact{},
	)
	if err != nil {
		log.Printf("Warning: failed to drop tables: %v", err)
	}

	// 重新創建表
	err = db.AutoMigrate(
		&entities.User{},
		&entities.Message{},
		&entities.Contact{},
		&entities.Group{},
		&entities.GroupMember{},
	)
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	// 創建頭像目錄
	avatarDir := "web/asset/avatars"
	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		log.Printf("Warning: failed to create avatar directory: %v", err)
	}

	log.Println("Database migration completed successfully")
}
