# GoChat - 基於 Clean Architecture 的即時通訊系統

## 目錄
- [專案概述](#專案概述)
- [功能特性](#功能特性)
- [技術架構](#技術架構)
- [安裝部署](#安裝部署)
- [使用說明](#使用說明)
- [API 文檔](#api-文檔)
- [開發指南](#開發指南)
- [常見問題](#常見問題)
- [版本歷史](#版本歷史)

## 專案概述
GoChat 是一個基於 Clean Architecture 設計的即時通訊系統，提供私聊和群聊功能。專案採用 Go 語言開發，使用 Gin 框架作為 Web 服務器，GORM 作為 ORM 框架，MySQL 作為數據庫，WebSocket 實現即時通訊。

## 功能特性
- 用戶管理
  - 註冊/登錄
  - 個人資料管理
  - 好友管理
- 即時通訊
  - 私聊功能
  - 群聊功能
  - 消息歷史記錄
- 群組管理
  - 創建/刪除群組
  - 群組成員管理
  - 群組設置

## 技術架構
### Clean Architecture 層次
- **Entities (實體層)**
  - 核心商業邏輯
  - 與業務概念直接相關的模型
  - 包含 User、Message、Group 等實體

- **Use Cases (用例層)**
  - 處理應用程式的商業邏輯
  - 包含 chat 和 user 服務
  - 實現具體的業務功能

- **Interface Adapters (介面適配器層)**
  - 將外部接口轉換為內部模型
  - 包含 controllers 和 routers
  - 處理 HTTP 請求和響應

- **Frameworks and Drivers (框架和驅動層)**
  - 具體的框架和技術實現
  - 包含 MySQL、Redis、WebSocket 等基礎設施
  - 處理與外部系統的交互

### 專案結構
```
clean-architecture-gochat/
├── cmd/                    # 應用程序入口
│   └── main.go            # 主程序入口
├── internal/              # 內部包
│   ├── domain/           # 領域層
│   │   ├── entities/     # 實體
│   │   │   ├── user.go
│   │   │   ├── message.go
│   │   │   └── group.go
│   │   └── repositories/ # 倉儲接口
│   │       ├── user_repository.go
│   │       ├── message_repository.go
│   │       └── group_repository.go
│   ├── infrastructure/   # 基礎設施層
│   │   ├── mysql/       # MySQL 實現
│   │   │   └── mysql.go
│   │   ├── redis/       # Redis 實現
│   │   │   └── redis.go
│   │   └── websocket/   # WebSocket 實現
│   │       └── websocket.go
│   └── interface/       # 接口層
│       ├── controllers/ # 控制器
│       │   ├── chat_controller.go
│       │   └── user_controller.go
│       └── routers/     # 路由
│           └── router.go
├── usecases/            # 用例層
│   ├── chat/           # 聊天相關用例
│   │   └── chat_service.go
│   └── user/           # 用戶相關用例
│       └── user_service.go
├── web/                # 前端資源
│   ├── asset/         # 靜態資源
│   └── views/         # 視圖模板
├── sql/               # 數據庫相關
│   └── migrate.go     # 數據庫遷移
└── test/              # 測試相關
    ├── chat_test.go
    ├── group_test.go
    └── user_test.go
```

### 技術棧
- 後端：Go 1.21
- Web 框架：Gin
- ORM：GORM v2
- 數據庫：MySQL 8.0
- 緩存：Redis
- 前端：HTML5, CSS3, JavaScript
- 即時通訊：WebSocket
- API 文檔：Swagger
- 測試框架：testing + testify

## 安裝部署
### 環境要求
- Go 1.21 或更高版本
- MySQL 8.0 或更高版本
- Redis 6.0 或更高版本
- Git

### 依賴安裝
```bash
# 安裝 Gin 框架
go get -u github.com/gin-gonic/gin

# 安裝 GORM v2
go get -u gorm.io/gorm@v1.21.12
go get -u gorm.io/driver/mysql@v1.5.7

# 安裝 Redis 客戶端
go get -u github.com/go-redis/redis/v8

# 安裝 WebSocket
go get -u github.com/gorilla/websocket

# 安裝 Swagger
go get -u github.com/swaggo/swag/cmd/swag@v1.6.5
go get -u github.com/swaggo/gin-swagger@v1.2.0
go get -u github.com/swaggo/files
go get -u github.com/alecthomas/template

# 安裝 Validator
go get -u github.com/go-playground/validator/v10

# 安裝測試相關
go get -u github.com/stretchr/testify

# 更新依賴
go mod tidy
```

### 安裝步驟
1. 克隆專案
```bash
git clone https://github.com/yourusername/clean-architecture-gochat.git
cd clean-architecture-gochat
```

2. 安裝依賴
```bash
go mod download
```

3. 配置數據庫
- 創建數據庫
```sql
CREATE DATABASE newgochat CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```
- 修改數據庫連接配置（在 `infrastructure/mysql/mysql.go` 中）

4. 運行數據庫遷移
```bash
go run sql/migrate.go
```

5. 生成 Swagger 文檔
```bash
swag init -g cmd/main.go
```

6. 啟動服務
```bash
go run cmd/main.go
```

## 使用說明
### 用戶註冊
1. 訪問 `/register` 頁面
2. 填寫註冊信息
3. 提交表單完成註冊

### 登錄系統
1. 訪問首頁
2. 輸入用戶名和密碼
3. 點擊登錄按鈕

### 聊天功能
- 私聊：選擇好友開始聊天
- 群聊：創建或加入群組進行群聊

## API 文檔
API 文檔可通過 Swagger UI 訪問：`http://localhost:8080/swagger/index.html`

### 主要 API 端點
- 用戶相關
  - POST `/user/create` - 創建用戶
  - POST `/user/login` - 用戶登錄
  - GET `/user/list` - 獲取用戶列表

- 聊天相關
  - GET `/chat/ws` - WebSocket 連接
  - POST `/chat/private/send` - 發送私聊消息
  - POST `/chat/group/send` - 發送群聊消息

- 群組相關
  - POST `/chat/group/create` - 創建群組
  - GET `/chat/groups` - 獲取群組列表
  - POST `/chat/group/:id/members` - 添加群組成員

## 開發指南
### 代碼規範
- 遵循 Go 官方代碼規範
- 使用 gofmt 格式化代碼
- 編寫單元測試
- 保持代碼簡潔清晰

### 開發流程
1. 創建功能分支
2. 實現功能
3. 編寫測試
4. 提交代碼
5. 發起 Pull Request

### 測試
#### 安裝測試工具
```bash
# 安裝 testify
go get -u github.com/stretchr/testify
```

#### 運行測試
```bash
# 運行所有測試
go test ./...

# 運行特定測試
go test ./test/...

# 運行測試並顯示詳細信息
go test -v ./...

# 運行測試並生成覆蓋率報告
go test -cover ./...
```

#### 測試示例
```go
func TestUserService(t *testing.T) {
    // 初始化測試環境
    db := setupTestDB(t)
    ctx := context.Background()

    // 初始化服務
    userRepo := repositories.NewUserRepository(db)
    userService := user.NewService(userRepo)

    // 測試用例
    t.Run("CreateUser", func(t *testing.T) {
        user, err := userService.CreateUser(ctx, "testuser", "password")
        assert.NoError(t, err)
        assert.NotNil(t, user)
        assert.Equal(t, "testuser", user.Username)
    })
}
```

## 常見問題
### 數據庫連接失敗
- 檢查數據庫服務是否運行
- 確認連接字符串配置正確
- 驗證數據庫用戶權限

### WebSocket 連接問題
- 檢查瀏覽器控制台錯誤
- 確認服務器是否正常運行
- 驗證網絡連接

### Redis 連接問題
- 檢查 Redis 服務是否運行
- 確認連接配置正確
- 驗證 Redis 版本兼容性

## 版本歷史
### v1.0.0 (2024-03-20)
- 初始版本發布
- 實現基本聊天功能
- 支持私聊和群聊
- 完成用戶管理功能

### v0.1.0 (2024-03-15)
- 專案初始化
- 基礎架構搭建
- 數據庫設計
- 用戶模組實現

Clean Architecture 一般包括以下層次：

Entities (Entities 層): 核心商業邏輯，通常是與業務概念直接相關的模型。
Use Cases (用例層): 處理應用程式的商業邏輯。
Interface Adapters (介面適配器層): 將外部接口轉換為內部模型，並將內部模型轉換為外部接口所需的格式。
Frameworks and Drivers (框架和驅動層): 具體的框架和技術，比如 Gin 路由，WebSocket，資料庫等。

clean-architecture-gochat/
- ├── cmd
- │   └── main.go
- ├── internal
- │   ├── config
- │   │   └── config.go
- │   └── domain
- │       ├── entities
- │       │   ├── user.go
- │       │   ├── message.go
- │       │   └── group.go
- │       └── repositories
- │           ├── user_repository.go
- │           ├── message_repository.go
- │           └── group_repository.go
- ├── infrastructure
- │   ├── mysql
- │   │   └── mysql.go
- │   ├── redis
- │   │   └── redis.go
- │   └── websocket
- │       └── websocket.go
- ├── interfaces
- │   ├── controllers
- │   │   ├── chat_controller.go
- │   │   └── user_controller.go
- │   └── routers
- │       └── router.go
- ├── usecases
- │   ├── chat
- │   │   └── chat_service.go
- │   └── user
- │       └── user_service.go
- └── pkg
- └── utils
- ├── crypto.go
- └── response.go

WebSocket 應該放在哪？
目前 WebSocket (websocket.go) 是放在 infrastructure/websocket，這是 正確的，因為：

1. WebSocket 屬於基礎設施 (Infrastructure)：

- WebSocket 負責網路通訊，不是業務邏輯的一部分，而是通訊機制，因此屬於 infrastructure。
- 就像 mysql 和 redis 一樣，websocket 是技術基礎設施，不是應用邏輯。
2. internal/interface/websocket/ 不是理想位置：

- interface/ 層通常用於 API 入口 (controllers)，但 WebSocket 不是 API，而是一種通訊技術，所以不應放在 interface/。

- 確保所有 WebSocket 相關邏輯都在這裡處理，而 業務邏輯 (聊天訊息) 仍然由 usecases/chat/chat_service.go 負責。

15. 安裝 Validator
    Validator 是一個 Go 語言的驗證庫，支持結構體驗證、自定義驗證規則和多語言錯誤信息。
    安裝 Validator：
    go get -u github.com/go-playground/validator/v10
    go mod tidy

8. 安裝 GORM (有 v1 跟 v2 版本)
   GORM 是一個 Go 語言的 ORM 框架，支持 MySQL、PostgreSQL、SQLite 和 SQL Server 等多種數據庫。
   安裝 GORM v1：
   go get -u github.com/jinzhu/gorm@v1.9.12
   go get -u github.com/jinzhu/gorm/dialects/mysql@v1.9.16
   go mod tidy
   安裝 GORM v2：
   go get -u gorm.io/gorm@v1.21.12
   go get -u gorm.io/driver/mysql@v1.5.7
   go mod tidy

go get -u github.com/go-redis/redis/v8

go get -u github.com/gorilla/websocket

go get -u github.com/gin-gonic/gin

11. 安裝 Swagger (要有main程式，如在cmd下要請 : swag init -g cmd/main.go )
    Swagger 是一個用於設計、構建、文件化和消費 RESTful API 的工具。
    安裝 Swagger：
    go get -u github.com/swaggo/swag/cmd/swag@v1.6.5
    go get -u github.com/swaggo/gin-swagger@v1.2.0
    go get -u github.com/swaggo/files
    go get -u github.com/alecthomas/template
    go mod tidy

12. 安裝測試依賴
    go get github.com/stretchr/testify
    更新 go.mod：go mod tidy
    運行測試：go test -v ./test/...