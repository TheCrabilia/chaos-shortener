resource "kubernetes_namespace_v1" "monitoring" {
  metadata {
    name   = "monitoring"
    labels = local.default_ns_labels
  }
}

resource "helm_release" "kube_prometheus_stack" {
  name       = "kube-prometheus-stack"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-prometheus-stack"
  version    = "69.6.0"
  namespace  = kubernetes_namespace_v1.monitoring.metadata[0].name
  values = [jsonencode({
    alertmanager = {
      ingress = {
        enabled = true
        hosts   = ["alertmanager.abracadabra-lab.net"]
      }
    },
    grafana = {
      ingress = {
        enabled = true
        hosts   = ["grafana.abracadabra-lab.net"]
      }
    },
    prometheus = {
      ingress = {
        enabled = true
        hosts   = ["prometheus.abracadabra-lab.net"]
      }
    }
  })]
}
