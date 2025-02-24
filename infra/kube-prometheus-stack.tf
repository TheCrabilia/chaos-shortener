resource "helm_release" "kube_prometheus_stack" {
  name             = "kube-prometheus-stack"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  namespace        = "monitoring"
  create_namespace = true
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
