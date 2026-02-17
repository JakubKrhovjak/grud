# =============================================================================
# EKS Cluster (equivalent to GKE cluster in terraform/gke.tf)
# =============================================================================
# Node Pool Strategy (mirrors GKE):
#   - system-pool: kube-system components (NO TAINT)
#   - infra-pool: monitoring stack (TAINT: workload=infra)
#   - app-pool: application workloads (TAINT: workload=app, autoscaling)
# =============================================================================

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 20.0"

  cluster_name    = var.cluster_name
  cluster_version = var.cluster_version

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  cluster_endpoint_public_access = true

  # IRSA - equivalent to GKE Workload Identity
  enable_irsa = true

  # Cluster addons
  cluster_addons = {
    coredns = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    vpc-cni = {
      most_recent = true
    }
  }

  # ==========================================================================
  # Node Groups - mirrors GKE node pool strategy
  # ==========================================================================
  eks_managed_node_groups = {

    # System pool - kube-system components (NO TAINT)
    # Equivalent to GKE system-pool (e2-medium, 2 nodes, stable)
    system-pool = {
      instance_types = [var.system_instance_type]
      capacity_type  = "ON_DEMAND"

      min_size     = 2
      max_size     = 2
      desired_size = 2

      disk_size = var.disk_size_gb

      labels = {
        "node-type" = "system"
      }
    }

    # Infra pool - monitoring stack (Prometheus, Grafana, Loki, Tempo, NATS)
    # Equivalent to GKE infra-pool (e2-standard-4, 2 nodes, spot)
    infra-pool = {
      instance_types = [var.infra_instance_type]
      capacity_type  = "SPOT"

      min_size     = 2
      max_size     = 2
      desired_size = 2

      disk_size = var.disk_size_gb

      labels = {
        "node-type" = "infra"
      }

      taints = [
        {
          key    = "workload"
          value  = "infra"
          effect = "NO_SCHEDULE"
        }
      ]
    }

    # App pool - application workloads (student-service, project-service)
    # Equivalent to GKE app-pool (e2-medium, autoscaling 1-4, spot)
    app-pool = {
      instance_types = [var.app_instance_type]
      capacity_type  = "SPOT"

      min_size     = var.app_min_size
      max_size     = var.app_max_size
      desired_size = var.app_min_size

      disk_size = var.disk_size_gb

      labels = {
        "node-type" = "app"
      }

      taints = [
        {
          key    = "workload"
          value  = "app"
          effect = "NO_SCHEDULE"
        }
      ]
    }
  }

  tags = {
    Environment = "test"
    Project     = "grud"
  }
}
