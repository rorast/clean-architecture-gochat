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
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3308)/newgochat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
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
