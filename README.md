# GoChat - 基於 Clean Architecture 的即時通訊應用

## 專案概述
GoChat 是一個基於 Clean Architecture 設計的即時通訊應用，使用 Go 語言開發，並部署在 Google Kubernetes Engine (GKE) 上。

### 主要功能
- 即時通訊
- 群組聊天
- 檔案分享
- 用戶管理

### 技術棧
- 後端：Go (Gin Framework)
- 資料庫：MariaDB
- 快取：Redis
- 容器化：Docker
- 編排：Kubernetes
- 雲端平台：Google Cloud Platform (GCP)

## 快速開始

### 前置需求
- Go 1.21 或以上
- Docker
- kubectl
- Google Cloud SDK
- 存取 GCP 專案的權限

### 本地開發
1. 克隆專案
```bash
git clone https://github.com/your-username/gochat.git
cd gochat
```

2. 安裝依賴
```bash
go mod download
```

3. 設定環境變數
```bash
cp .env.example .env
# 編輯 .env 檔案設定必要的環境變數
```

4. 運行應用
```bash
go run main.go
```

### 部署到 GKE
詳細的部署指南請參考 [DEPLOYMENT.md](./docs/DEPLOYMENT.md)

## 專案結構
```
.
├── cmd/                    # 應用程式入口點
├── internal/              # 內部套件
│   ├── domain/           # 領域模型
│   ├── infrastructure/   # 基礎設施層
│   ├── interface/        # 介面層
│   └── usecase/         # 用例層
├── pkg/                  # 可重用的套件
├── sql/                  # 資料庫遷移和種子
├── web/                  # 前端資源
└── k8s/                  # Kubernetes 配置
```

## 文檔
- [部署指南](./docs/DEPLOYMENT.md)
- [Kubernetes 配置說明](./docs/K8S.md)
- [常用指令手冊](./docs/COMMANDS.md)

## 貢獻指南
1. Fork 專案
2. 建立功能分支
3. 提交變更
4. 發起 Pull Request

## 授權
MIT License

## 功能特點

- 即時聊天
- 群組聊天
- 檔案傳輸
- 語音訊息
- 表情符號支援
- 用戶管理
- 權限控制

## 系統需求

- Go 1.21 或更高版本
- Docker 和 Docker Compose
- MariaDB 10.11
- Redis 6.2

## 開發指南

### 專案結構

```
.
├── cmd/                # 應用入口
├── internal/          # 內部包
│   ├── domain/       # 領域模型
│   ├── repository/   # 資料存取層
│   ├── usecase/      # 業務邏輯層
│   └── delivery/     # 交付層
├── pkg/              # 公共包
├── web/              # Web 相關資源
│   ├── asset/       # 靜態資源
│   └── views/       # 視圖模板
└── scripts/         # 腳本文件
```

### 本地開發

1. 使用 Air 進行熱重載
```bash
# 安裝 Air
go install github.com/cosmtrek/air@latest

# 運行 Air
air
```

2. 資料庫遷移
```bash
# 執行遷移
go run cmd/migrate/migrate.go
```

### 測試

```bash
# 運行所有測試
go test ./...

# 運行特定包的測試
go test ./internal/...
```

## Docker 環境說明

### 開發環境

- `Dockerfile.dev`: 開發環境的 Dockerfile
- `docker-compose.yml`: 開發環境的 Docker Compose 配置
- 支援熱重載
- 掛載本地代碼目錄

### 生產環境

- `Dockerfile`: 生產環境的 Dockerfile
- `docker-compose.prod.yml`: 生產環境的 Docker Compose 配置
- 多階段構建
- 最小化鏡像大小
- 環境變數配置

## 環境變數

### 開發環境

```env
DB_HOST=mariadb
DB_PORT=3306
DB_USER=gochat
DB_PASSWORD=gochat123
DB_NAME=gochat
REDIS_HOST=redis
REDIS_PORT=6379
```

### 生產環境

```env
DB_USER=your_secure_user
DB_PASSWORD=your_secure_password
DB_NAME=your_database_name
MYSQL_ROOT_PASSWORD=your_secure_root_password
REDIS_PASSWORD=your_secure_redis_password
APP_ENV=production
APP_PORT=8080
APP_SECRET=your_secure_app_secret
```

## 安全性考慮

- 使用環境變數管理敏感信息
- 生產環境使用安全的密碼
- 限制端口訪問
- 定期更新依賴包

## 監控和維護

- 使用 Docker 日誌進行監控
- 定期備份資料庫
- 監控系統資源使用情況

## 聯繫方式

- 專案維護者：[您的名字]
- 電子郵件：[您的郵箱]
- 專案連結：[GitHub 專案地址]

## Kubernetes 部署指南

### 前置需求

1. 安裝必要工具
```bash
# 安裝 gcloud CLI (Windows PowerShell)
(New-Object Net.WebClient).DownloadFile("https://dl.google.com/dl/cloudsdk/channels/rapid/GoogleCloudSDKInstaller.exe", "$env:TEMP\GoogleCloudSDKInstaller.exe")
Start-Process -FilePath "$env:TEMP\GoogleCloudSDKInstaller.exe" -ArgumentList "/S"

# 安裝 kubectl
# Windows 用戶可以通過 gcloud 安裝
gcloud components install kubectl

# 安裝 kustomize
# Windows 用戶可以使用 Chocolatey
choco install kustomize
```

2. GCP 專案設置
```bash
# 初始化 gcloud
gcloud init

# 設置專案 ID
gcloud config set project marine-embassy-455211-m1

# 啟用必要的 API
gcloud services enable container.googleapis.com compute.googleapis.com
```

3. 啟用 GCP 計費
- 訪問 https://console.cloud.google.com/billing/enable?project=marine-embassy-455211-m1
- 選擇或創建計費帳戶
- 將計費帳戶關聯到專案

4. 創建 GKE 集群
```bash
# 創建集群
gcloud container clusters create gochat-cluster \
    --num-nodes=3 \
    --zone=asia-east1-a \
    --machine-type=e2-medium \
    --disk-size=20GB

# 獲取集群認證
gcloud container clusters get-credentials gochat-cluster --zone asia-east1-a

# 驗證集群連接
kubectl cluster-info
```

### 部署應用

1. 配置 Docker 認證
```bash
# 配置 Docker 使用 GCP 容器註冊表
gcloud auth configure-docker
```

2. 構建和推送 Docker 映像
```bash
# 構建映像
docker build -t gcr.io/marine-embassy-455211-m1/gochat:latest .

# 推送到 GCR
docker push gcr.io/marine-embassy-455211-m1/gochat:latest
```

3. 部署到 Kubernetes
```bash
# 使用 kustomize 部署
kubectl apply -k k8s/overlays/prod

# 查看部署狀態
kubectl get pods
kubectl get services
```

### 測試部署

1. 檢查 Pod 狀態
```bash
# 查看所有 Pod
kubectl get pods

# 查看特定 Pod 的詳細信息
kubectl describe pod <pod-name>

# 查看 Pod 日誌
kubectl logs <pod-name>
```

2. 檢查服務狀態
```bash
# 查看所有服務
kubectl get services

# 查看特定服務的詳細信息
kubectl describe service gochat
```

3. 測試應用功能
```bash
# 獲取服務的外部 IP
kubectl get service gochat

# 使用 curl 測試 API
curl http://<EXTERNAL-IP>/health

# 測試 WebSocket 連接
# 使用瀏覽器訪問 http://<EXTERNAL-IP>
```

4. 監控應用
```bash
# 查看 Pod 資源使用情況
kubectl top pods

# 查看節點資源使用情況
kubectl top nodes

# 查看事件
kubectl get events --sort-by='.lastTimestamp'
```

### 故障排除

1. 常見問題解決
```bash
# 檢查 Pod 是否正常運行
kubectl get pods -o wide

# 檢查 Pod 日誌
kubectl logs <pod-name>

# 檢查 Pod 描述
kubectl describe pod <pod-name>

# 檢查服務配置
kubectl describe service gochat
```

2. 擴展和更新
```bash
# 手動擴展副本數
kubectl scale deployment gochat --replicas=5

# 更新映像
kubectl set image deployment/gochat gochat=gcr.io/marine-embassy-455211-m1/gochat:new-version

# 回滾部署
kubectl rollout undo deployment/gochat --to-revision=2
```

3. 清理資源
```bash
# 刪除部署
kubectl delete -k k8s/overlays/prod

# 刪除集群
gcloud container clusters delete gochat-cluster --zone asia-east1-a
```

### 安全建議

1. 配置網絡策略
```bash
# 創建網絡策略
kubectl apply -f k8s/network-policy.yaml
```

2. 配置資源限制
```bash
# 查看資源使用情況
kubectl top pods
kubectl top nodes
```

3. 定期更新
```bash
# 更新集群
gcloud container clusters upgrade gochat-cluster --zone asia-east1-a

# 更新節點
gcloud container clusters upgrade gochat-cluster --zone asia-east1-a --node-pool=default-pool
```

## 資料庫配置

### MariaDB 配置

本專案使用 MariaDB 10.11 作為主要資料庫。在 Kubernetes 環境中，MariaDB 的配置如下：

1. 資料庫憑證
```yaml
# 使用 Kubernetes Secrets 存儲敏感信息
- MYSQL_ROOT_PASSWORD
- DB_USER
- DB_PASSWORD
```

2. 資料庫連接配置
```yaml
# 開發環境
DB_HOST=mariadb
DB_PORT=3306
DB_NAME=gochat

# 生產環境
DB_HOST=mariadb
DB_PORT=3306
DB_NAME=gochat
```

3. 資料持久化
- 使用 Kubernetes PersistentVolume 存儲資料
- 資料存儲在 `/var/lib/mysql` 目錄

4. 資料庫備份
```bash
# 備份資料庫
kubectl exec -it <mariadb-pod-name> -- mysqldump -u root -p gochat > backup.sql

# 還原資料庫
kubectl exec -i <mariadb-pod-name> -- mysql -u root -p gochat < backup.sql
```

5. 資料庫監控
- 使用 Kubernetes 的資源監控
- 監控資料庫連接數
- 監控查詢性能
- 監控磁碟使用情況