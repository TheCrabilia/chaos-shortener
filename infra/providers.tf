terraform {
  required_version = ">= 1.9"
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "3.2.3"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.17.0"
    }
    http = {
      source  = "hashicorp/http"
      version = "3.4.5"
    }
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = "1.19.0"
    }
  }
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
