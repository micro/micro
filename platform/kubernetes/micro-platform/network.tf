locals {
  network_name         = "go.micro.network"
  network_port         = 8085
  network_service_port = 8080
  network_labels = merge(
    local.common_labels,
    {
      "name" = local.network_name
    }
  )
  network_annotations = merge(
    local.common_annotations,
    {
      "name" = local.network_name
    }
  )
  network_env = merge(
    local.common_env_vars,
    {
      "MICRO_SERVER_ADDRESS"             = "0.0.0.0:8080"
      "MICRO_NETWORK_TOKEN"              = "micro.mu"
      "MICRO_NETWORK_ADVERTISE_STRATEGY" = "best"
    }
  )
}

module "network_cert" {
  source = "./cert"

  ca_cert_pem        = tls_self_signed_cert.platform_ca_cert.cert_pem
  ca_private_key_pem = tls_private_key.platform_ca_key.private_key_pem
  private_key_alg    = var.private_key_alg

  subject = local.network_name
}

resource "kubernetes_secret" "network_cert" {
  metadata {
    name        = "${replace(local.network_name, ".", "-")}-cert"
    namespace   = var.platform_namespace
    labels      = local.network_labels
    annotations = local.network_annotations
  }
  data = {
    "cert.pem" = module.network_cert.cert_pem
    "key.pem"  = module.network_cert.key_pem
  }
  type = "Opaque"
}

resource "kubernetes_deployment" "network" {
  metadata {
    name        = replace(local.network_name, ".", "-")
    namespace   = var.platform_namespace
    labels      = local.network_labels
    annotations = merge(local.common_annotations, local.network_annotations)
  }
  spec {
    replicas = 1
    selector {
      match_labels = local.network_labels
    }
    template {
      metadata {
        labels = local.network_labels
      }
      spec {
        container {
          name = replace(local.network_name, ".", "-")
          dynamic "env" {
            for_each = local.network_env
            content {
              name  = env.key
              value = env.value
            }
          }
          args              = ["network"]
          image             = var.micro_image
          image_pull_policy = var.image_pull_policy
          port {
            container_port = local.network_port
            name           = "network-port"
            protocol       = "UDP"
          }
          port {
            container_port = local.network_service_port
            name           = "service-port"
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
            secret_name  = kubernetes_secret.network_cert.metadata[0].name
          }
        }
        automount_service_account_token = true
      }
    }
  }
}

resource "kubernetes_service" "network" {
  metadata {
    name        = replace(local.network_name, ".", "-")
    namespace   = var.platform_namespace
    labels      = local.network_labels
    annotations = merge(local.common_annotations, local.network_annotations)
  }
  spec {
    port {
      port        = local.network_port
      target_port = local.network_port
      name        = "network-port"
      protocol    = "UDP"
    }
    port {
      port        = local.network_service_port
      target_port = local.network_service_port
      name        = "service-port"
    }
    selector = local.network_labels
  }
}
