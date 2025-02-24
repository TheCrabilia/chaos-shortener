data "http" "subscription_manifest" {
  url = var.subscription_url
}

resource "kubectl_manifest" "apply" {
  yaml_body = data.http.subscription_manifest.response_body
}
