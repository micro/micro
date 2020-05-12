locals {
  bot_name = "go.micro.bot"
  bot_labels = merge(
    local.common_labels,
    {
      "name" = local.bot_name
    }
  )
  bot_annotations = merge(
    local.common_annotations,
    {
      "name" = local.bot_name
    }
  )
  bot_env = merge(
    local.common_env_vars,
    {
      "MICRO_AUTH" = "SERVICE"
    }
  )
}

module "bot_cert" {
  source = "./cert"

  ca_cert_pem        = tls_self_signed_cert.platform_ca_cert.cert_pem
  ca_private_key_pem = tls_private_key.platform_ca_key.private_key_pem
  private_key_alg    = var.private_key_alg

  subject = local.bot_name
}

resource "kubernetes_secret" "bot_cert" {
  metadata {
    name        = "${replace(local.bot_name, ".", "-")}-cert"
    namespace   = var.platform_namespace
    labels      = local.bot_labels
    annotations = local.bot_annotations
  }
  data = {
    "cert.pem" = module.bot_cert.cert_pem
    "key.pem"  = module.bot_cert.key_pem
  }
  type = "Opaque"
}

resource "kubernetes_deployment" "bot" {
  metadata {
    name        = replace(local.bot_name, ".", "-")
    namespace   = var.platform_namespace
    labels      = local.bot_labels
    annotations = local.bot_annotations
  }
  spec {
    replicas = 1
    selector {
      match_labels = local.bot_labels
    }
    template {
      metadata {
        labels = local.bot_labels
      }
      spec {
        container {
          name = replace(local.bot_name, ".", "-")
          dynamic "env" {
            for_each = local.bot_env
            content {
              name  = env.key
              value = env.value
            }
          }
          env {
            name = "MICRO_SLACK_TOKEN"
            value_from {
              secret_key_ref {
                key  = "token"
                name = kubernetes_secret.slack_token.metadata[0].name
              }
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
          args              = ["bot", "--inputs=slack"]
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
            secret_name  = kubernetes_secret.bot_cert.metadata[0].name
          }
        }
        automount_service_account_token = true
      }
    }
  }
}
