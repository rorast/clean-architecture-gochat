variable "project_id" {
  description = "GCP 專案 ID"
  type        = string
}

variable "region" {
  description = "GCP 區域"
  type        = string
  default     = "asia-east1"
}

variable "zone" {
  description = "GCP 可用區"
  type        = string
  default     = "asia-east1-a"
}

variable "cluster_name" {
  description = "GKE 集群名稱"
  type        = string
  default     = "gochat-cluster"
}

variable "node_count" {
  description = "節點數量"
  type        = number
  default     = 1
}

variable "machine_type" {
  description = "節點機器類型"
  type        = string
  default     = "e2-medium"
}

variable "disk_size_gb" {
  description = "節點磁碟大小（GB）"
  type        = number
  default     = 20
} 