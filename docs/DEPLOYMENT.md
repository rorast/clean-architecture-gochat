# GoChat 部署指南

本文檔詳細說明如何將 GoChat 部署到 Google Kubernetes Engine (GKE)。

## 前置需求

1. 安裝必要工具
```bash
# 安裝 Google Cloud SDK
# 安裝 kubectl
# 安裝 Docker
```

2. 設定 GCP 專案
```bash
# 設定 GCP 專案 ID
export PROJECT_ID=your-project-id

# 設定 GCP 區域
export REGION=asia-east1
```

## 自動化部署 (GitHub Actions)

### 1. 設定 GitHub Secrets

在 GitHub 專案設定中添加以下 Secrets：

- `GCP_SA_KEY`: GCP 服務帳號金鑰（JSON 格式）

### 2. 工作流程說明

當代碼推送到 `main` 分支時，GitHub Actions 會自動執行以下步驟：

1. 構建並推送 Docker 映像
   - 應用程式映像
   - 資料庫遷移映像

2. 執行資料庫遷移
   - 刪除舊的遷移任務
   - 部署新的遷移任務
   - 等待遷移完成

3. 部署應用程式
   - 更新部署配置
   - 部署新版本
   - 驗證部署狀態

### 3. 手動觸發部署

如果需要手動觸發部署，可以：

1. 在 GitHub 專案頁面中：
   - 點擊 "Actions" 標籤
   - 選擇 "Deploy to GKE" 工作流程
   - 點擊 "Run workflow"

2. 使用 GitHub CLI：
```bash
gh workflow run deploy.yml
```

## 部署步驟

### 1. 準備 Docker Images

#### 1.1 設定 Docker 認證
```bash
# 登入 GCP Container Registry
gcloud auth configure-docker
```

#### 1.2 建立應用程式 Image
```bash
# 建立應用程式 Image
docker build -t gcr.io/${PROJECT_ID}/gochat:latest .

# 推送 Image 到 GCR
docker push gcr.io/${PROJECT_ID}/gochat:latest
```

#### 1.3 建立資料庫遷移 Image
```bash
# 建立遷移工具 Image
docker build -t gcr.io/${PROJECT_ID}/gochat-migration:latest -f sql/Dockerfile .

# 推送 Image 到 GCR
docker push gcr.io/${PROJECT_ID}/gochat-migration:latest
```

### 2. 設定 Kubernetes 配置

#### 2.1 更新配置中的專案 ID
需要修改以下檔案中的專案 ID：
- `k8s/base/deployment.yaml`
- `k8s/base/migration-job.yaml`

將 `marine-embassy-455211-m1` 替換為你的專案 ID。

#### 2.2 建立 Kubernetes Secrets
```bash
# 建立資料庫認證 Secret
kubectl create secret generic gochat-secrets \
  --from-literal=DB_USER=gochat \
  --from-literal=DB_PASSWORD=gochat123
```

### 3. 部署應用

#### 3.1 部署資料庫遷移
```bash
# 部署遷移 Job
kubectl apply -f k8s/base/migration-job.yaml

# 檢查遷移狀態
kubectl get jobs
kubectl logs -l job-name=db-migration
```

#### 3.2 部署應用程式
```bash
# 部署應用
kubectl apply -f k8s/base/deployment.yaml

# 部署服務
kubectl apply -f k8s/base/service.yaml
```

### 4. 驗證部署

#### 4.1 檢查 Pod 狀態
```bash
kubectl get pods -l app=gochat
```

#### 4.2 檢查服務狀態
```bash
kubectl get svc gochat
```

#### 4.3 獲取外部 IP
```bash
EXTERNAL_IP=$(kubectl get svc gochat -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo "應用程式可通過 http://${EXTERNAL_IP} 訪問"
```

## 更新部署

### 1. 更新應用程式
```bash
# 建立新的 Image
docker build -t gcr.io/${PROJECT_ID}/gochat:latest .

# 推送 Image
docker push gcr.io/${PROJECT_ID}/gochat:latest

# 重新部署
kubectl rollout restart deployment gochat
```

### 2. 更新資料庫結構
```bash
# 建立新的遷移 Image
docker build -t gcr.io/${PROJECT_ID}/gochat-migration:latest -f sql/Dockerfile .

# 推送 Image
docker push gcr.io/${PROJECT_ID}/gochat-migration:latest

# 執行遷移
kubectl delete job db-migration
kubectl apply -f k8s/base/migration-job.yaml
```

## 故障排除

### 1. 檢查 Pod 日誌
```bash
kubectl logs -l app=gochat
```

### 2. 檢查 Pod 狀態
```bash
kubectl describe pod -l app=gochat
```

### 3. 檢查服務狀態
```bash
kubectl describe svc gochat
```

## 清理資源

### 1. 刪除部署
```bash
kubectl delete -f k8s/base/
```

### 2. 刪除 Secrets
```bash
kubectl delete secret gochat-secrets
```

### 3. 刪除 Images
```bash
# 刪除 GCR 中的 Images
gcloud container images delete gcr.io/${PROJECT_ID}/gochat:latest
gcloud container images delete gcr.io/${PROJECT_ID}/gochat-migration:latest
``` 