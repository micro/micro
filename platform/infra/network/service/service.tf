locals {
  common_labels = {
    "service-name"  = var.service_name
  }
  common_annotations = {
    "name"    = "go.micro.${var.service_name}"
    "version" = "latest"
    "source"  = "github.com/micro/micro"
    "owner"   = "micro"
    "group"   = "micro"
  }
  common_env_vars = {
    "MICRO_LOG_LEVEL"         = "debug"
    "MICRO_BROKER"            = "nats"
    "MICRO_BROKER_ADDRESS"    = "nats-cluster.${var.resource_namespace}.svc"
    "MICRO_REGISTRY"          = "etcd"
    "MICRO_REGISTRY_ADDRESS"  = "etcd-cluster.${var.resource_namespace}.svc"
    "MICRO_REGISTER_TTL"      = "60"
    "MICRO_REGISTER_INTERVAL" = "30"
  }
}

resource "kubernetes_service" "network_service" {
  count = var.create_k8s_service ? 1 : 0
  metadata {
    name      = "micro-${var.service_name}"
    namespace = var.network_namespace
    labels    = merge(local.common_labels, var.extra_labels)
  }
  spec {
    type = var.service_type
    port {
      name        = "${var.service_name}-port"
      port        = var.service_port
      protocol    = var.service_protocol
      target_port = "${var.service_name}-port"
    }
    selector = merge(local.common_labels, var.extra_labels)
  }
  lifecycle {
    ignore_changes = [metadata.0.annotations]
  }
}

resource "kubernetes_ingress" "network_ingress" {
  count = var.create_k8s_ingress ? 1 : 0
  metadata {
    name      = "micro-${var.service_name}"
    namespace = var.network_namespace
    labels    = merge(local.common_labels, var.extra_labels)
    annotations = {
      // We only expose services that manage their own certificates
      "haproxy.org/ssl-passthrough" = "true"
    }
  }
  spec {
    tls {
      hosts = var.domain_names
    }
    dynamic "rule" {
      for_each = var.domain_names
      content {
        host = rule.value
        http {
          path {
            backend {
              service_name = kubernetes_service.network_service.0.metadata.0.name
              service_port = var.service_port
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_deployment" "network_deployment" {
  metadata {
    name      = "micro-${var.service_name}"
    namespace = var.network_namespace
    labels    = merge(local.common_labels, var.extra_labels)
  }
  spec {
    replicas = var.service_replicas
    selector {
      match_labels = merge(local.common_labels, var.extra_labels)
    }
    template {
      metadata {
        labels = merge(local.common_labels, var.extra_labels)
      }
      spec {
        container {
          name              = var.service_name
          args              = split("-", var.service_name)
          image             = var.micro_image
          image_pull_policy = var.image_pull_policy
          dynamic "port" {
            for_each = var.create_k8s_service ? [1] : []
            content {
              name           = "${var.service_name}-port"
              container_port = var.service_port
              protocol       = var.service_protocol
            }
          }
          dynamic "env" {
            for_each = merge(local.common_env_vars, var.extra_env_vars)
            content {
              name  = env.key
              value = env.value
            }
          }
        }
      }
    }
  }
}
