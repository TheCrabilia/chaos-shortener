terraform {
  required_version = ">= 1.9"
  required_providers {
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
