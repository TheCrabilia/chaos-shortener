resource "kubectl_manifest" "cshort_db" {
  yaml_body = <<EOF
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: cshort-db
  namespace: chaos-shortener
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

  depends_on = [module.cloudnative_pg_manifest]
}

resource "kubectl_manifest" "cshort_pooler" {
  yaml_body = <<EOF
apiVersion: postgresql.cnpg.io/v1
kind: Pooler
metadata:
  name: cshort-pooler
  namespace: chaos-shortener
spec:
  cluster:
    name: cshort-db
  instances: 3
  type: rw
  pgbouncer:
    poolMode: session
    parameters:
      max_client_conn: "1000"
      default_pool_size: "10"
EOF

  depends_on = [module.cloudnative_pg_manifest]
}

resource "helm_release" "cshort" {
  name             = "chaos-shortener"
  chart            = "../charts/chaos-shortener"
  namespace        = "chaos-shortener"
  create_namespace = true
  recreate_pods    = true
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
