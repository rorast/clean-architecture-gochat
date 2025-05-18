# 暫時移除所有輸出 

output "kubernetes_cluster_name" {
  value       = google_container_cluster.gochat_cluster.name
  description = "GKE 集群名稱"
}

output "kubernetes_cluster_host" {
  value       = google_container_cluster.gochat_cluster.endpoint
  description = "GKE 集群主機"
}

output "kubernetes_cluster_ca_certificate" {
  value       = base64decode(google_container_cluster.gochat_cluster.master_auth[0].cluster_ca_certificate)
  description = "GKE 集群 CA 證書"
  sensitive   = true
} 