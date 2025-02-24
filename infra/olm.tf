resource "null_resource" "olm_install" {
  provisioner "remote-exec" {
    connection {
      host     = var.cp_host
      user     = var.cp_username
      password = var.cp_password
    }

    inline = ["curl -sL https://github.com/operator-framework/operator-lifecycle-manager/releases/download/v0.31.0/install.sh | bash -s v0.31.0"]
  }
}

module "cloudnative_pg_manifest" {
  source           = "./modules/olm-operator"
  subscription_url = "https://operatorhub.io/install/cloudnative-pg.yaml"
}

resource "kubectl_manifest" "cshort_postgresql" {
  yaml_body = <<EOF
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: cshort-postgresql
spec:
  instances: 1
  storage:
    size: 1Gi
  monitoring:
    enablePodMonitor: true
  EOF

  depends_on = [module.cloudnative_pg_manifest]
}
