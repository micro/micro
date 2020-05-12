locals {
  debug_name = "go.micro.debug"
  debug_labels = merge(
    local.common_labels,
    {
      "name" = local.debug_name
    }
  )
  debug_annotations = merge(
    local.common_annotations,
    {
      "name" = local.debug_name
    }
  )
  debug_env = merge(
    local.common_env_vars,
    {
      "MICRO_DEBUG_LOG"    = "service"
      "MICRO_DEBUG_WINDOW" = "600"
      "MICRO_AUTH"         = "jwt"
    }
  )
}

module "debug_cert" {
  source = "./cert"

  ca_cert_pem        = tls_self_signed_cert.platform_ca_cert.cert_pem
  ca_private_key_pem = tls_private_key.platform_ca_key.private_key_pem
  private_key_alg    = var.private_key_alg

  subject = local.debug_name
}

resource "kubernetes_secret" "debug_cert" {
  metadata {
    name        = "${replace(local.debug_name, ".", "-")}-cert"
    namespace   = var.platform_namespace
    labels      = local.debug_labels
    annotations = local.debug_annotations
  }
  data = {
    "cert.pem" = module.debug_cert.cert_pem
    "key.pem"  = module.debug_cert.key_pem
  }
  type = "Opaque"
}

resource "kubernetes_deployment" "debug" {
  metadata {
    name        = replace(local.debug_name, ".", "-")
    namespace   = var.platform_namespace
    labels      = local.debug_labels
    annotations = local.debug_annotations
  }
  spec {
    replicas = 1
    selector {
      match_labels = local.debug_labels
    }
    template {
      metadata {
        labels = local.debug_labels
      }
      spec {
        container {
          name = replace(local.debug_name, ".", "-")
          dynamic "env" {
            for_each = local.debug_env
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
          args              = ["debug"]
          image             = var.micro_image
          image_pull_policy = var.image_pull_policy
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
            secret_name  = kubernetes_secret.debug_cert.metadata[0].name
          }
        }
        automount_service_account_token = true
      }
    }
  }
}
