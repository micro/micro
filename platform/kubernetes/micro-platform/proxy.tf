locals {
  proxy_name       = "go.micro.proxy"
  proxy_port       = 8081
  proxy_https_port = 443
  proxy_labels = merge(
    local.common_labels,
    {
      "name" = local.proxy_name
    }
  )
  proxy_annotations = merge(
    local.common_annotations,
    {
      "name" = local.proxy_name
    }
  )
  proxy_env = merge(
    local.common_env_vars,
    {
      "MICRO_PROXY_ADDRESS"  = "0.0.0.0:${local.proxy_https_port}",
      "MICRO_SERVER_ADDRESS" = "0.0.0.0:${local.proxy_port}"
      "MICRO_ENABLE_STATS"   = "true",
      "MICRO_ENABLE_ACME"    = "true",
      "MICRO_ACME_PROVIDER"  = "certmagic",
      "MICRO_ACME_HOSTS"     = var.proxy_acme_hosts
      "MICRO_STORE"          = "service"
      "MICRO_STORE_ADDRESS"  = ""
      "MICRO_STORE_DATABASE" = "micro"
      "MICRO_STORE_TABLE"    = "micro"
      "MICRO_AUTH"           = "service"
    }
  )
}

module "proxy_cert" {
  source = "./cert"

  ca_cert_pem        = tls_self_signed_cert.platform_ca_cert.cert_pem
  ca_private_key_pem = tls_private_key.platform_ca_key.private_key_pem
  private_key_alg    = var.private_key_alg

  subject = local.proxy_name
}

resource "kubernetes_secret" "proxy_cert" {
  metadata {
    name        = "${replace(local.proxy_name, ".", "-")}-cert"
    namespace   = kubernetes_namespace.platform.id
    labels      = local.proxy_labels
    annotations = local.proxy_annotations
  }
  data = {
    "cert.pem" = module.proxy_cert.cert_pem
    "key.pem"  = module.proxy_cert.key_pem
  }
  type = "Opaque"
}

resource "kubernetes_deployment" "proxy" {
  metadata {
    name        = replace(local.proxy_name, ".", "-")
    namespace   = kubernetes_namespace.platform.id
    labels      = local.proxy_labels
    annotations = local.proxy_annotations
  }
  spec {
    replicas = 1
    selector {
      match_labels = local.proxy_labels
    }
    strategy {
      rolling_update {
        max_surge       = "0"
        max_unavailable = "1"
      }
    }
    template {
      metadata {
        labels = local.proxy_labels
      }
      spec {
        container {
          name = replace(local.proxy_name, ".", "-")
          dynamic "env" {
            for_each = local.proxy_env
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
          args              = ["proxy"]
          image             = var.micro_image
          image_pull_policy = var.image_pull_policy
          port {
            container_port = local.proxy_port
            name           = "proxy"
          }
          port {
            container_port = local.proxy_https_port
            name           = "https"
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
            secret_name  = kubernetes_secret.proxy_cert.metadata[0].name
          }
        }
        automount_service_account_token = true
      }
    }
  }
}

resource "kubernetes_service" "proxy" {
  metadata {
    name        = replace(local.proxy_name, ".", "-")
    namespace   = kubernetes_namespace.platform.id
    labels      = local.proxy_labels
    annotations = local.proxy_annotations
  }
  spec {
    port {
      name        = "proxy"
      port        = local.proxy_port
      target_port = local.proxy_port
    }
    port {
      name        = "https"
      port        = local.proxy_https_port
      target_port = local.proxy_https_port
    }
    selector = local.proxy_labels
    type     = "LoadBalancer"
  }
}
