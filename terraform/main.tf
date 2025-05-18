terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
  zone    = var.zone
}

# GKE 集群
resource "google_container_cluster" "gochat_cluster" {
  name     = "gochat-cluster"
  location = var.zone

  # 使用較小的初始配置
  initial_node_count = 1
  
  # 移除默認節點池
  remove_default_node_pool = true

  network    = google_compute_network.vpc.name
  subnetwork = google_compute_subnetwork.subnet.name

  # 啟用 Workload Identity
  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  # 使用次要 IP 範圍
  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }

  # 啟用 VPC 原生
  networking_mode = "VPC_NATIVE"

  # 維護窗口
  maintenance_policy {
    recurring_window {
      start_time = "2024-01-01T00:00:00Z"
      end_time   = "2024-12-31T23:59:59Z"
      recurrence = "FREQ=WEEKLY;BYDAY=SA,SU"
    }
  }
}

# 節點池
resource "google_container_node_pool" "primary_nodes" {
  name       = "gochat-node-pool"
  location   = var.zone
  cluster    = google_container_cluster.gochat_cluster.name
  
  # 設置初始節點數
  initial_node_count = 1

  # 自動擴展配置
  autoscaling {
    min_node_count = 1
    max_node_count = 2
  }

  # 管理配置
  management {
    auto_repair  = true
    auto_upgrade = true
  }

  # 升級設置
  upgrade_settings {
    max_surge       = 1
    max_unavailable = 0
  }

  node_config {
    # 基本配置
    machine_type = "e2-small"
    disk_size_gb = 20
    disk_type    = "pd-standard"
    image_type   = "COS_CONTAINERD"

    # 標籤和標記
    labels = {
      environment = "production"
      node_pool   = "primary"
    }
    tags = ["gochat-node"]

    # Workload Identity 配置
    workload_metadata_config {
      mode = "GKE_METADATA"
    }

    # OAuth 範圍
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]

    # 元數據
    metadata = {
      disable-legacy-endpoints = "true"
    }
  }

  lifecycle {
    ignore_changes = [
      initial_node_count,
      node_count
    ]
  }
}

# VPC 網絡
resource "google_compute_network" "vpc" {
  name                    = "gochat-vpc"
  auto_create_subnetworks = false
}

# 子網
resource "google_compute_subnetwork" "subnet" {
  name          = "gochat-subnet"
  ip_cidr_range = "10.2.0.0/16"
  region        = "asia-east1"
  network       = google_compute_network.vpc.name

  secondary_ip_range {
    range_name    = "pods"
    ip_cidr_range = "10.10.0.0/16"
  }

  secondary_ip_range {
    range_name    = "services"
    ip_cidr_range = "10.20.0.0/16"
  }
}

# 防火牆規則
resource "google_compute_firewall" "allow_internal" {
  name    = "allow-internal"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  source_ranges = ["10.2.0.0/16", "10.10.0.0/16", "10.20.0.0/16"]
}

resource "google_compute_firewall" "allow_ssh" {
  name    = "allow-ssh"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]
} 

# 根據配置文件，我們可以看到應用程式的部署架構如下：
# 基礎設施（由 Terraform 管理）：
# VPC 網絡：gochat-vpc
# 子網：gochat-subnet
# GKE 集群：gochat-cluster
# 節點池：gochat-node-pool（e2-small 機型，1-2 個節點）
# Kubernetes 資源（由 Kustomize 管理）：
# 應用程式部署（gochat）：
# 容器映像：gcr.io/marine-embassy-455211-m1/gochat:latest
# 端口：8099
# 資源限制：CPU 100m-200m，記憶體 128Mi-256Mi
# MariaDB 部署：
# 容器映像：mariadb:10.6
# 端口：3306
# 持久化：使用 emptyDir（臨時存儲）
# Redis 部署：
# 容器映像：redis:6.2
# 端口：6379
# 持久化：使用 emptyDir（臨時存儲）
# 配置和密鑰：
# ConfigMap：gochat-config（包含資料庫和 Redis 連接資訊）
# Secret：gochat-secrets（包含資料庫認證資訊）
# 服務：
# gochat：LoadBalancer 類型，暴露端口 80
# mariadb：ClusterIP 類型，內部訪問
# redis：ClusterIP 類型，內部訪問
# 等待 Terraform 完成 GKE 集群的部署後，我們需要：
# 獲取集群認證
# 創建必要的 Kubernetes 資源
# 部署應用程式和依賴服務