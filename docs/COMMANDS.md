# GoChat 常用指令手冊

本文檔收錄了 GoChat 專案中常用的指令，方便快速查詢和使用。

## Docker 相關指令

### 1. 映像管理

```bash
# 建立應用程式映像
docker build -t gcr.io/${PROJECT_ID}/gochat:latest .

# 建立遷移工具映像
docker build -t gcr.io/${PROJECT_ID}/gochat-migration:latest -f sql/Dockerfile .

# 推送映像到 GCR
docker push gcr.io/${PROJECT_ID}/gochat:latest
docker push gcr.io/${PROJECT_ID}/gochat-migration:latest

# 查看本地映像
docker images | grep gochat
```

### 2. 容器管理

```bash
# 運行容器
docker run -d -p 8080:8080 gcr.io/${PROJECT_ID}/gochat:latest

# 查看運行中的容器
docker ps | grep gochat

# 停止容器
docker stop $(docker ps -q --filter ancestor=gcr.io/${PROJECT_ID}/gochat:latest)

# 查看容器日誌
docker logs $(docker ps -q --filter ancestor=gcr.io/${PROJECT_ID}/gochat:latest)
```

## Kubernetes 相關指令

### 1. 部署管理

```bash
# 部署所有資源
kubectl apply -f k8s/base/

# 刪除所有資源
kubectl delete -f k8s/base/

# 重新部署
kubectl rollout restart deployment gochat

# 查看部署狀態
kubectl get deployments
kubectl get pods
kubectl get services
```

### 2. Pod 管理

```bash
# 查看 Pod 狀態
kubectl get pods -l app=gochat

# 查看 Pod 詳細資訊
kubectl describe pod -l app=gochat

# 查看 Pod 日誌
kubectl logs -l app=gochat

# 進入 Pod 容器
kubectl exec -it $(kubectl get pod -l app=gochat -o jsonpath='{.items[0].metadata.name}') -- /bin/sh
```

### 3. 服務管理

```bash
# 查看服務狀態
kubectl get svc gochat

# 查看服務詳細資訊
kubectl describe svc gochat

# 獲取外部 IP
kubectl get svc gochat -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
```

### 4. 資料庫遷移

```bash
# 執行遷移
kubectl apply -f k8s/base/migration-job.yaml

# 查看遷移狀態
kubectl get jobs
kubectl logs -l job-name=db-migration

# 刪除遷移任務
kubectl delete job db-migration
```

### 5. 故障排除

```bash
# 查看 Pod 事件
kubectl get events --sort-by='.lastTimestamp'

# 查看 Pod 日誌
kubectl logs -l app=gochat --tail=100

# 查看 Pod 詳細狀態
kubectl describe pod -l app=gochat
```

## GCP 相關指令

### 1. 專案設定

```bash
# 設定專案 ID
export PROJECT_ID=your-project-id

# 設定區域
export REGION=asia-east1

# 設定叢集
gcloud container clusters get-credentials gochat-cluster --region ${REGION}
```

### 2. 映像管理

```bash
# 列出 GCR 中的映像
gcloud container images list --repository=gcr.io/${PROJECT_ID}

# 刪除映像
gcloud container images delete gcr.io/${PROJECT_ID}/gochat:latest
gcloud container images delete gcr.io/${PROJECT_ID}/gochat-migration:latest
```

## 常用組合指令

### 1. 完整部署流程

```bash
# 設定環境變數
export PROJECT_ID=your-project-id
export REGION=asia-east1

# 建立並推送映像
docker build -t gcr.io/${PROJECT_ID}/gochat:latest .
docker build -t gcr.io/${PROJECT_ID}/gochat-migration:latest -f sql/Dockerfile .
docker push gcr.io/${PROJECT_ID}/gochat:latest
docker push gcr.io/${PROJECT_ID}/gochat-migration:latest

# 部署應用
kubectl apply -f k8s/base/
```

### 2. 更新部署

```bash
# 更新映像
docker build -t gcr.io/${PROJECT_ID}/gochat:latest .
docker push gcr.io/${PROJECT_ID}/gochat:latest

# 重新部署
kubectl rollout restart deployment gochat

# 監控部署狀態
kubectl rollout status deployment gochat
```

### 3. 清理環境

```bash
# 刪除所有資源
kubectl delete -f k8s/base/

# 刪除 Secrets
kubectl delete secret gochat-secrets

# 刪除映像
gcloud container images delete gcr.io/${PROJECT_ID}/gochat:latest
gcloud container images delete gcr.io/${PROJECT_ID}/gochat-migration:latest
``` 