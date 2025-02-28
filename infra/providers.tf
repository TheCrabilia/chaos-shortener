terraform {
  required_version = ">= 1.9"
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.17.0"
    }
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = "1.19.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.36.0"
    }
  }
}

provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = "abracadabra-lab-k8s"
}

provider "helm" {
  kubernetes {
    config_path    = "~/.kube/config"
    config_context = "abracadabra-lab-k8s"
  }
}

provider "kubectl" {
  config_path    = "~/.kube/config"
  config_context = "abracadabra-lab-k8s"
}
