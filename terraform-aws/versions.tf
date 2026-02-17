terraform {
  required_version = ">= 1.14.0"

  # For testing, use local state. Switch to S3 backend for production.
  # backend "s3" {
  #   bucket = "grud-terraform-state"
  #   key    = "eks/terraform.tfstate"
  #   region = "eu-central-1"
  # }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.25"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.12"
    }
  }
}

provider "aws" {
  region = var.region
}

data "aws_eks_cluster" "cluster" {
  count = var.skip_kubernetes_provider ? 0 : 1
  name  = module.eks.cluster_name
}

data "aws_eks_cluster_auth" "cluster" {
  count = var.skip_kubernetes_provider ? 0 : 1
  name  = module.eks.cluster_name
}

provider "kubernetes" {
  host                   = var.skip_kubernetes_provider ? "https://localhost" : module.eks.cluster_endpoint
  token                  = var.skip_kubernetes_provider ? "" : data.aws_eks_cluster_auth.cluster[0].token
  cluster_ca_certificate = var.skip_kubernetes_provider ? "" : base64decode(module.eks.cluster_certificate_authority_data)
}

provider "helm" {
  kubernetes {
    host                   = var.skip_kubernetes_provider ? "https://localhost" : module.eks.cluster_endpoint
    token                  = var.skip_kubernetes_provider ? "" : data.aws_eks_cluster_auth.cluster[0].token
    cluster_ca_certificate = var.skip_kubernetes_provider ? "" : base64decode(module.eks.cluster_certificate_authority_data)
  }
}
