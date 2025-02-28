resource "kubernetes_namespace_v1" "chaos_shortener" {
  metadata {
    name   = "chaos-shortener"
    labels = local.default_ns_labels
  }
}

resource "kubernetes_namespace_v1" "operators" {
  metadata {
    name = "operators"
  }
}

resource "helm_release" "cnpg" {
  name       = "cloudnative-pg"
  repository = "https://cloudnative-pg.github.io/charts"
  chart      = "cloudnative-pg"
  namespace  = kubernetes_namespace_v1.operators.metadata[0].name
}

resource "kubectl_manifest" "cshort_db" {
  yaml_body = <<EOF
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: cshort-db
  namespace: ${kubernetes_namespace_v1.chaos_shortener.metadata[0].name}
spec:
  instances: 1
  bootstrap:
    initdb:
      database: cshort
      owner: cshort
      secret:
        name: cshort-db-credentials
  storage:
    size: 1Gi
EOF

  depends_on = [helm_release.cnpg]
}

resource "kubectl_manifest" "cshort_pooler" {
  yaml_body = <<EOF
apiVersion: postgresql.cnpg.io/v1
kind: Pooler
metadata:
  name: cshort-pooler
  namespace: ${kubernetes_namespace_v1.chaos_shortener.metadata[0].name}
spec:
  cluster:
    name: cshort-db
  instances: 1
  type: rw
  pgbouncer:
    poolMode: session
    parameters:
      max_client_conn: "1000"
      default_pool_size: "10"
EOF

  depends_on = [helm_release.cnpg]
}

resource "helm_release" "cshort" {
  name          = "chaos-shortener"
  chart         = "../charts/chaos-shortener"
  namespace     = kubernetes_namespace_v1.chaos_shortener.metadata[0].name
  recreate_pods = true
  values = [jsonencode({
    chaosShortener = {
      database = {
        host           = "cshort-pooler"
        database       = "cshort"
        existingSecret = "cshort-db-credentials"
      }
    }
  })]
}
