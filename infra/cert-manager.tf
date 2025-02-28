resource "kubernetes_namespace_v1" "cert_manager" {
  metadata {
    name   = "cert-manager"
    labels = local.default_ns_labels
  }
}

resource "helm_release" "certmanager" {
  name       = "cert-manager"
  repository = "https://charts.jetstack.io"
  chart      = "cert-manager"
  namespace  = kubernetes_namespace_v1.cert_manager.metadata[0].name
  values = [jsonencode({
    crds = {
      enabled = true
    }
  })]
}
