package main

import (
	"clean-architecture-gochat/infrastructure/mysql"
	"clean-architecture-gochat/internal/config"
	"clean-architecture-gochat/internal/domain/entities"
	"fmt"
	"log"
)

func main() {
	// **修正: 載入設定檔*
	config.LoadConfig("internal/config") // 確保 `Config` 變數已初始化

	// 取得 DB 連線
	db := mysql.Connect()
	if db == nil {
		log.Fatal("Failed to initialize database connection.")
	}

	// 遷移表結構
	fmt.Println("🔄 開始數據庫遷移...")
	err := db.AutoMigrate(
		&entities.User{},
		&entities.Message{},
		&entities.Group{},
	)
	if err != nil {
		log.Fatalf("❌ 數據庫遷移失敗: %v", err)
	}

	fmt.Println("✅ 數據庫遷移成功！")
}
