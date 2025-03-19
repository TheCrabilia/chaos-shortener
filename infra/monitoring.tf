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
      config = {
        route = {
          receiver = "telegram"
          routes   = []
        }
        receivers = [{
          name = "telegram"
          telegram_configs = [{
            bot_token = var.telegram_api_key
            chat_id   = var.telegram_chat_id
            message   = <<-EOF
              <b>Alert: {{ .CommonLabels.alertname }}</b>
              <br><i>{{ .CommonAnnotations.summary }}</i>
              <br>
              <b>Status:</b> {{ .Status }}
              <br>
              <b>Instance(s):</b> {{ range .Alerts }}{{ .Labels.instance }} {{ end }}
              <br>
              <b>Starts At:</b> {{ .CommonLabels.startsAt }}
              <br>
              <b>Description:</b>
              <br>{{ .CommonAnnotations.description }}
              <br>
              {{ if gt (len .Alerts.Firing) 0 }}
                <i>Total firing alerts: {{ len .Alerts.Firing }}</i>
              {{ end }}
            EOF
          }]
        }]
      }
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
