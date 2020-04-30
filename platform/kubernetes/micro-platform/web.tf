locals {
  web_name = "go.micro.web"
  web_port = 443
  web_labels = merge(
    local.common_labels,
    {
      "name" = local.web_name
    }
  )
  web_annotations = merge(
    local.common_annotations,
    {
      "name" = local.web_name
    }
  )
  web_env = merge(
    local.common_env_vars,
    {
      "MICRO_AUTH"           = "service",
      "MICRO_AUTH_LOGIN_URL" = "https://account.m3o.dev"
      "MICRO_ENABLE_STATS"   = "true",
      "MICRO_ENABLE_ACME"    = "true",
      "MICRO_ACME_PROVIDER"  = "certmagic",
      "MICRO_ACME_HOSTS"     = var.web_acme_hosts
      "MICRO_STORE"          = "service"
      "MICRO_STORE_ADDRESS"  = ""
      "MICRO_STORE_DATABASE" = "micro"
      "MICRO_STORE_TABLE"    = "micro"
      "MICRO_WEB_NAMESPACE"  = "domain"
      "MICRO_WEB_RESOLVER"   = "subdomain"
    }
  )
}

module "web_cert" {
  source = "./cert"

  ca_cert_pem        = tls_self_signed_cert.platform_ca_cert.cert_pem
  ca_private_key_pem = tls_private_key.platform_ca_key.private_key_pem
  private_key_alg    = var.private_key_alg

  subject = local.web_name
}

resource "kubernetes_secret" "web_cert" {
  metadata {
    name        = "${replace(local.web_name, ".", "-")}-cert"
    namespace   = var.platform_namespace
    labels      = local.web_labels
    annotations = local.web_annotations
  }
  data = {
    "cert.pem" = module.web_cert.cert_pem
    "key.pem"  = module.web_cert.key_pem
  }
  type = "Opaque"
}

resource "kubernetes_deployment" "web" {
  metadata {
    name        = replace(local.web_name, ".", "-")
    namespace   = var.platform_namespace
    labels      = local.web_labels
    annotations = merge(local.common_annotations, local.web_annotations)
  }
  spec {
    replicas = 1
    selector {
      match_labels = local.web_labels
    }
    strategy {
      rolling_update {
        max_surge       = "0"
        max_unavailable = "1"
      }
    }
    template {
      metadata {
        labels = local.web_labels
      }
      spec {
        container {
          name = replace(local.web_name, ".", "-")
          dynamic "env" {
            for_each = local.web_env
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
          env {
            name = "MICRO_AUTH_PUBLIC_KEY"
            value_from {
              secret_key_ref {
                key  = "public"
                name = kubernetes_secret.micro_keypair.metadata[0].name
              }
            }
          }
          args              = ["web"]
          image             = var.micro_image
          image_pull_policy = var.image_pull_policy
          port {
            container_port = local.web_port
            name           = replace(local.web_name, ".", "-")
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
            secret_name  = kubernetes_secret.web_cert.metadata[0].name
          }
        }
        automount_service_account_token = true
      }
    }
  }
}

resource "kubernetes_service" "web" {
  metadata {
    name        = replace(local.web_name, ".", "-")
    namespace   = var.platform_namespace
    labels      = local.web_labels
    annotations = merge(local.common_annotations, local.web_annotations)
  }
  spec {
    port {
      name        = "https"
      port        = local.web_port
      target_port = local.web_port
    }
    selector = local.web_labels
    type     = "LoadBalancer"
  }
}
