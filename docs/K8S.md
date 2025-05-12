# Kubernetes 配置說明

本文檔詳細說明 GoChat 專案中的 Kubernetes 配置文件。

## 配置文件結構

```
k8s/
└── base/
    ├── deployment.yaml    # 應用程式部署配置
    ├── service.yaml       # 服務配置
    └── migration-job.yaml # 資料庫遷移任務配置
```

## 配置文件說明

### 1. deployment.yaml

主要應用程式的部署配置。

#### 重要配置項：
- `replicas`: Pod 副本數量
- `image`: 應用程式容器映像
- `ports`: 容器端口配置
- `env`: 環境變數配置
- `resources`: 資源限制
- `livenessProbe`: 存活探針
- `readinessProbe`: 就緒探針

#### 需要替換的項目：
- `image`: 將 `marine-embassy-455211-m1` 替換為你的 GCP 專案 ID

### 2. service.yaml

負載平衡器服務配置。

#### 重要配置項：
- `type`: 服務類型（LoadBalancer）
- `ports`: 端口映射配置
- `selector`: Pod 選擇器

#### 需要替換的項目：
- 無需替換

### 3. migration-job.yaml

資料庫遷移任務配置。

#### 重要配置項：
- `image`: 遷移工具容器映像
- `env`: 環境變數配置
- `restartPolicy`: 重啟策略

#### 需要替換的項目：
- `image`: 將 `marine-embassy-455211-m1` 替換為你的 GCP 專案 ID

## 常用指令

### 1. 部署相關

```bash
# 部署所有資源
kubectl apply -f k8s/base/

# 部署特定資源
kubectl apply -f k8s/base/deployment.yaml
kubectl apply -f k8s/base/service.yaml
kubectl apply -f k8s/base/migration-job.yaml

# 刪除所有資源
kubectl delete -f k8s/base/

# 刪除特定資源
kubectl delete -f k8s/base/deployment.yaml
kubectl delete -f k8s/base/service.yaml
kubectl delete -f k8s/base/migration-job.yaml
```

### 2. 查看狀態

```bash
# 查看部署狀態
kubectl get deployments
kubectl get pods
kubectl get services

# 查看詳細資訊
kubectl describe deployment gochat
kubectl describe pod -l app=gochat
kubectl describe service gochat

# 查看日誌
kubectl logs -l app=gochat
```

### 3. 更新部署

```bash
# 重新部署
kubectl rollout restart deployment gochat

# 查看部署歷史
kubectl rollout history deployment gochat

# 回滾到特定版本
kubectl rollout undo deployment gochat --to-revision=1
```

### 4. 擴縮容

```bash
# 擴展副本數
kubectl scale deployment gochat --replicas=5

# 自動擴縮容
kubectl autoscale deployment gochat --min=2 --max=5 --cpu-percent=80
```

## 注意事項

1. 部署前請確保：
   - 已正確設定 GCP 專案 ID
   - 已建立必要的 Secrets
   - 已推送 Docker 映像到 GCR

2. 更新部署時：
   - 先更新 Docker 映像
   - 再執行 `kubectl rollout restart`

3. 資料庫遷移：
   - 建議在應用程式更新前執行
   - 確保遷移成功後再更新應用程式

4. 監控：
   - 定期檢查 Pod 狀態
   - 監控資源使用情況
   - 查看應用程式日誌 