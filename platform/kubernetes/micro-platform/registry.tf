locals {
  registry_name = "go.micro.registry"
  registry_port = 8000
  registry_labels = merge(
    local.common_labels,
    {
      "name" = local.registry_name
    }
  )
  registry_annotations = merge(
    local.common_annotations,
    {
      "name" = local.registry_name
    }
  )
  registry_env = merge(
    local.common_env_vars,
    {
    }
  )
}

module "registry_cert" {
  source = "./cert"

  ca_cert_pem        = tls_self_signed_cert.platform_ca_cert.cert_pem
  ca_private_key_pem = tls_private_key.platform_ca_key.private_key_pem
  private_key_alg    = var.private_key_alg

  subject = local.registry_name
}

resource "kubernetes_secret" "registry_cert" {
  metadata {
    name        = "${replace(local.registry_name, ".", "-")}-cert"
    namespace   = var.platform_namespace
    labels      = local.registry_labels
    annotations = local.registry_annotations
  }
  data = {
    "cert.pem" = module.registry_cert.cert_pem
    "key.pem"  = module.registry_cert.key_pem
  }
  type = "Opaque"
}

resource "kubernetes_deployment" "registry" {
  metadata {
    name        = replace(local.registry_name, ".", "-")
    namespace   = var.platform_namespace
    labels      = local.registry_labels
    annotations = local.registry_annotations
  }
  spec {
    replicas = 1
    selector {
      match_labels = local.registry_labels
    }
    template {
      metadata {
        labels = local.registry_labels
      }
      spec {
        container {
          name = replace(local.registry_name, ".", "-")
          dynamic "env" {
            for_each = local.registry_env
            content {
              name  = env.key
              value = env.value
            }
          }
          args              = ["registry"]
          image             = var.micro_image
          image_pull_policy = var.image_pull_policy
          port {
            container_port = local.registry_port
            name           = "auth-port"
          }
          volume_mount {
            mount_path = "/etc/micro/certs"
            name       = "certs"
          }
          volume_mount {
            mount_path = "/etc/micro/ca"
            name       = "platform-ca"
          }
        }
        volume {
          name = "platform-ca"
          secret {
            secret_name  = kubernetes_secret.platform_ca.metadata[0].name
            default_mode = "0600"
            items {
              key  = "ca.pem"
              path = "ca.pem"
            }
          }
        }
        volume {
          name = "certs"
          secret {
            default_mode = "0600"
            secret_name  = kubernetes_secret.registry_cert.metadata[0].name
          }
        }
        automount_service_account_token = true
      }
    }
  }
}

resource "kubernetes_service" "registry" {
  metadata {
    name        = replace(local.registry_name, ".", "-")
    namespace   = var.platform_namespace
    labels      = local.registry_labels
    annotations = local.registry_annotations
  }
  spec {
    port {
      port        = local.registry_port
      target_port = local.registry_port
    }
    selector = local.registry_labels
  }
}
