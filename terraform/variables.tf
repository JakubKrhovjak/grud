variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP Region"
  type        = string
  default     = "europe-west1"
}

variable "cluster_name" {
  description = "GKE cluster name"
  type        = string
  default     = "grud-cluster"
}

variable "network_name" {
  description = "VPC network name"
  type        = string
  default     = "grud-network"
}

variable "node_count" {
  description = "Number of nodes per zone"
  type        = number
  default     = 1
}

variable "machine_type" {
  description = "GKE node machine type"
  type        = string
  default     = "e2-small"
}

variable "disk_size_gb" {
  description = "Node disk size in GB"
  type        = number
  default     = 20
}

variable "min_node_count" {
  description = "Minimum nodes for autoscaling"
  type        = number
  default     = 1
}

variable "max_node_count" {
  description = "Maximum nodes for autoscaling"
  type        = number
  default     = 3
}

variable "enable_autopilot" {
  description = "Enable GKE Autopilot mode"
  type        = bool
  default     = false
}
