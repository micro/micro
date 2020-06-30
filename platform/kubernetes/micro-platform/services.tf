// Which services to start
variable "services" {
  description = "Core services to deploy on which port"
  type        = map(number)
  default = {
    "config"   = 8080,
    "auth"     = 8010,
    "network"  = 8085,
    "runtime"  = 8088,
    "registry" = 8000,
    "broker"   = 8001,
    "store"    = 8002,
    "router"   = 8084,
    "debug"    = 8080,
    "proxy"    = 443,
    "api"      = 443,
    "web"      = 443,
  }
}

variable "external_services" {
  description = "services to expose externally"
  type        = list(string)
  default = [
    "api",
    "proxy",
    "web",
  ]
}

variable "external_service_type" {
  type        = string
  description = "LoadBalancer for Cloud, NodePort for kind"
  default     = "NodePort"
}

variable "per_service_overrides" {
  type        = map(map(string))
  description = "Per service env var overrides. Merged with local.core_config below"
  default = {
    "auth" = { "MICRO_AUTH" = "jwt" },
  }
}

locals {
  // Default: all config talks to the infrastructure directly
  core_config = {
    "MICRO_LOG_LEVEL"         = "debug",
    "MICRO_AUTH"              = "jwt",
    "MICRO_AUTH_PRIVATE_KEY"  = base64encode(tls_private_key.platform_jwt_key.private_key_pem)
    "MICRO_AUTH_PUBLIC_KEY"   = base64encode(tls_private_key.platform_jwt_key.public_key_pem)
    "MICRO_BROKER"            = "nats",
    "MICRO_BROKER_ADDRESS"    = "nats-cluster.${var.resource_namespace}.svc",
    "MICRO_REGISTRY"          = "etcd",
    "MICRO_REGISTRY_ADDRESS"  = "etcd-cluster.${var.resource_namespace}.svc",
    "MICRO_REGISTER_TTL"      = "60",
    "MICRO_REGISTER_INTERVAL" = "30",
    "MICRO_STORE"             = "cockroach",
    "MICRO_STORE_ADDRESS"     = "postgres://root@cockroachdb-public.${var.resource_namespace}.svc:26257/?sslmode=disable",
  }
}

resource "kubernetes_service_account" "unprivileged" {
  metadata {
    name      = "micro-unprivileged"
    namespace = kubernetes_namespace.platform.metadata[0].name
  }
  automount_service_account_token = false

  dynamic "image_pull_secret" {
    for_each = length(var.image_pull_secret) > 0 ? { "a" = "a" } : {}
    content {
      name = kubernetes_secret.image_pull_secret[0].metadata[0].name
    }
  }
}

resource "kubernetes_deployment" "core_services" {
  for_each = var.services
  metadata {
    name      = replace(each.key, " ", "-")
    namespace = kubernetes_namespace.platform.metadata[0].name
    labels = {
      "micro" = replace(each.key, " ", "-")
    }
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        "micro" = replace(each.key, " ", "-")
      }
    }
    template {
      metadata {
        name      = replace(each.key, " ", "-")
        namespace = kubernetes_namespace.platform.metadata[0].name
        labels = {
          "micro" = replace(each.key, " ", "-")
        }
      }
      spec {
        container {
          name              = each.key
          args              = split(" ", each.key)
          image             = var.micro_image
          image_pull_policy = var.image_pull_policy
          port {
            container_port = each.value
            name           = "${replace(each.key, " ", "-")}-port"
          }
          dynamic "env" {
            for_each = merge(local.core_config, lookup(var.per_service_overrides, each.key, {}))
            content {
              name  = env.key
              value = env.value
            }
          }
          volume_mount {
            mount_path = "/etc/micro/certs"
            name       = "cert-bundle"
          }
        }
        service_account_name            = kubernetes_service_account.unprivileged.metadata[0].name
        automount_service_account_token = false
        volume {
          name = "cert-bundle"
          secret {
            default_mode = "0600"
            secret_name  = kubernetes_secret.cert_bundles[each.key].metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "external_facing_services" {
  count = length(var.external_services)

  metadata {
    name      = replace(var.external_services[count.index], " ", "-")
    namespace = kubernetes_namespace.platform.id
    labels    = { "micro" = replace(var.external_services[count.index], " ", "-") }
  }
  spec {
    port {
      name        = "https"
      port        = 443
      target_port = var.services[var.external_services[count.index]]
    }
    selector = {
      "micro" = replace(var.external_services[count.index], " ", "-")
    }
    type = var.external_service_type
  }
}