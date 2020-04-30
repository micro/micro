locals {
  auth_name = "go.micro.auth"
  auth_port = 8000
  auth_labels = merge(
    local.common_labels,
    {
      "name" = local.auth_name
    }
  )
  auth_annotations = merge(
    local.common_annotations,
    {
      "name" = local.auth_name
    }
  )
  auth_env = merge(
    local.common_env_vars,
    {
    }
  )
}

module "auth_cert" {
  source = "./cert"

  ca_cert_pem        = tls_self_signed_cert.platform_ca_cert.cert_pem
  ca_private_key_pem = tls_private_key.platform_ca_key.private_key_pem
  private_key_alg    = var.private_key_alg

  subject = local.auth_name
}

resource "kubernetes_secret" "auth_cert" {
  metadata {
    name        = "${replace(local.auth_name, ".", "-")}-cert"
    namespace   = var.platform_namespace
    labels      = local.auth_labels
    annotations = local.auth_annotations
  }
  data = {
    "cert.pem" = module.auth_cert.cert_pem
    "key.pem"  = module.auth_cert.key_pem
  }
  type = "Opaque"
}

resource "kubernetes_deployment" "auth" {
  metadata {
    name        = replace(local.auth_name, ".", "-")
    namespace   = var.platform_namespace
    labels      = local.auth_labels
    annotations = merge(local.common_annotations, local.auth_annotations)
  }
  spec {
    replicas = 1
    selector {
      match_labels = local.auth_labels
    }
    template {
      metadata {
        labels = local.auth_labels
      }
      spec {
        container {
          name = replace(local.auth_name, ".", "-")
          dynamic "env" {
            for_each = local.auth_env
            content {
              name  = env.key
              value = env.value
            }
          }
          env {
            name = "MICRO_AUTH_PUBLIC_KEY"
            value_from {
              secret_key_ref {
                key  = "public"
                name = kubernetes_secret.micro_keypair.metadata[0].name
              }
            }
          }
          env {
            name = "MICRO_AUTH_PRIVATE_KEY"
            value_from {
              secret_key_ref {
                key  = "private"
                name = kubernetes_secret.micro_keypair.metadata[0].name
              }
            }
          }
          args              = ["auth"]
          image             = var.micro_image
          image_pull_policy = var.image_pull_policy
          port {
            container_port = local.auth_port
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
            secret_name  = kubernetes_secret.auth_cert.metadata[0].name
          }
        }
        automount_service_account_token = true
      }
    }
  }
}

resource "kubernetes_service" "auth" {
  metadata {
    name        = replace(local.auth_name, ".", "-")
    namespace   = var.platform_namespace
    labels      = local.auth_labels
    annotations = merge(local.common_annotations, local.auth_annotations)
  }
  spec {
    port {
      port        = local.auth_port
      target_port = local.auth_port
    }
    selector = local.auth_labels
  }
}
