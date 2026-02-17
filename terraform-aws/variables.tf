variable "region" {
  description = "AWS Region"
  type        = string
  default     = "eu-central-1"
}

variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
  default     = "grud-cluster"
}

variable "cluster_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.31"
}

variable "skip_kubernetes_provider" {
  description = "Skip kubernetes/helm provider configuration (for bootstrap)"
  type        = bool
  default     = false
}

# VPC
variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
  default     = "10.0.0.0/16"
}

# Node pools - mirrors GKE variables
variable "system_instance_type" {
  description = "Instance type for system nodes"
  type        = string
  default     = "t3.small"
}

variable "infra_instance_type" {
  description = "Instance type for infra nodes (monitoring stack)"
  type        = string
  default     = "t3.medium"
}

variable "app_instance_type" {
  description = "Instance type for app nodes"
  type        = string
  default     = "t3.small"
}

variable "app_min_size" {
  description = "Minimum app nodes for autoscaling"
  type        = number
  default     = 1
}

variable "app_max_size" {
  description = "Maximum app nodes for autoscaling"
  type        = number
  default     = 4
}

variable "disk_size_gb" {
  description = "Node disk size in GB"
  type        = number
  default     = 20
}
