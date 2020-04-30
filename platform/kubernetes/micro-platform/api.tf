locals {
  api_name = "go.micro.api"
  api_port = 443
  api_labels = merge(
    local.common_labels,
    {
      "name" = local.api_name
    }
  )
  api_annotations = merge(
    local.common_annotations,
    {
      "name" = local.api_name
    }
  )
  api_env = merge(
    local.common_env_vars,
    {
      "MICRO_API_NAMESPACE"  = "domain",
      "MICRO_ENABLE_STATS"   = "true",
      "MICRO_ENABLE_ACME"    = "true",
      "MICRO_ACME_PROVIDER"  = "certmagic",
      "MICRO_ACME_HOSTS"     = var.api_acme_hosts
      "MICRO_STORE"          = "service"
      "MICRO_STORE_ADDRESS"  = ""
      "MICRO_STORE_DATABASE" = "micro"
      "MICRO_STORE_TABLE"    = "micro"
    }
  )
}

module "api_cert" {
  source = "./cert"

  ca_cert_pem        = tls_self_signed_cert.platform_ca_cert.cert_pem
  ca_private_key_pem = tls_private_key.platform_ca_key.private_key_pem
  private_key_alg    = var.private_key_alg

  subject = local.api_name
}

resource "kubernetes_secret" "api_cert" {
  metadata {
    name        = "${replace(local.api_name, ".", "-")}-cert"
    namespace   = var.platform_namespace
    labels      = local.api_labels
    annotations = local.api_annotations
  }
  data = {
    "cert.pem" = module.api_cert.cert_pem
    "key.pem"  = module.api_cert.key_pem
  }
  type = "Opaque"
}

resource "kubernetes_deployment" "api" {
  metadata {
    name        = replace(local.api_name, ".", "-")
    namespace   = var.platform_namespace
    labels      = local.api_labels
    annotations = merge(local.common_annotations, local.api_annotations)
  }
  spec {
    replicas = 1
    selector {
      match_labels = local.api_labels
    }
    strategy {
      rolling_update {
        max_surge       = "0"
        max_unavailable = "1"
      }
    }
    template {
      metadata {
        labels = local.api_labels
      }
      spec {
        container {
          name = replace(local.api_name, ".", "-")
          dynamic "env" {
            for_each = local.api_env
            content {
              name  = env.key
              value = env.value
            }
          }
          env {
            name = "CF_API_TOKEN"
            value_from {
              secret_key_ref {
                key  = "CF_API_TOKEN"
                name = kubernetes_secret.cloudflare_credentals.metadata[0].name
              }
            }
          }
          args              = ["api"]
          image             = var.micro_image
          image_pull_policy = var.image_pull_policy
          port {
            container_port = local.api_port
            name           = replace(local.api_name, ".", "-")
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
            secret_name  = kubernetes_secret.api_cert.metadata[0].name
          }
        }
        automount_service_account_token = true
      }
    }
  }
}

resource "kubernetes_service" "api" {
  metadata {
    name        = replace(local.api_name, ".", "-")
    namespace   = var.platform_namespace
    labels      = local.api_labels
    annotations = merge(local.common_annotations, local.api_annotations)
  }
  spec {
    port {
      name        = "https"
      port        = local.api_port
      target_port = local.api_port
    }
    selector = local.api_labels
    type     = "LoadBalancer"
  }
}
