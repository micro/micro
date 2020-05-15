locals {
  store_name = "go.micro.store"
  store_port = 8082
  store_labels = merge(
    local.common_labels,
    {
      "name" = local.store_name
    }
  )
  store_annotations = merge(
    local.common_annotations,
    {
      "name" = local.store_name
    }
  )
  store_env = merge(
    local.common_env_vars,
    {
      "MICRO_SERVER_ADDRESS" = "0.0.0.0:${local.store_port}"
      "MICRO_AUTH"           = "jwt"
    }
  )
}

module "store_cert" {
  source = "./cert"

  ca_cert_pem        = tls_self_signed_cert.platform_ca_cert.cert_pem
  ca_private_key_pem = tls_private_key.platform_ca_key.private_key_pem
  private_key_alg    = var.private_key_alg

  subject = local.store_name
}

resource "kubernetes_secret" "store_cert" {
  metadata {
    name        = "${replace(local.store_name, ".", "-")}-cert"
    namespace   = kubernetes_namespace.platform.id
    labels      = local.store_labels
    annotations = local.store_annotations
  }
  data = {
    "cert.pem" = module.store_cert.cert_pem
    "key.pem"  = module.store_cert.key_pem
  }
  type = "Opaque"
}

resource "kubernetes_deployment" "store" {
  metadata {
    name        = replace(local.store_name, ".", "-")
    namespace   = kubernetes_namespace.platform.id
    labels      = local.store_labels
    annotations = local.store_annotations
  }
  spec {
    replicas = 1
    selector {
      match_labels = local.store_labels
    }
    template {
      metadata {
        labels = local.store_labels
      }
      spec {
        container {
          name = replace(local.store_name, ".", "-")
          dynamic "env" {
            for_each = local.store_env
            content {
              name  = env.key
              value = env.value
            }
          }
          env {
            name = "MICRO_AUTH_PUBLIC_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.micro_keypair.metadata[0].name
                key  = "public"
              }
            }
          }
          env {
            name = "MICRO_AUTH_PRIVATE_KEY"
            value_from {
              secret_key_ref {
                name = kubernetes_secret.micro_keypair.metadata[0].name
                key  = "private"
              }
            }
          }
          args              = ["store"]
          image             = var.micro_image
          image_pull_policy = var.image_pull_policy
          port {
            container_port = local.store_port
            name           = "service"
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
            secret_name  = kubernetes_secret.store_cert.metadata[0].name
          }
        }
        automount_service_account_token = true
      }
    }
  }
}
