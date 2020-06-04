locals {
  runtime_name = "go.micro.runtime"
  runtime_labels = merge(
    local.common_labels,
    {
      "name" = local.runtime_name
    }
  )
  runtime_annotations = merge(
    local.common_annotations,
    {
      "name" = local.runtime_name
    }
  )
  runtime_env = merge(
    local.common_env_vars,
    {
      "MICRO_RUNTIME"         = "kubernetes"
      "MICRO_RUNTIME_PROFILE" = "platform"
      "MICRO_AUTO_UPDATE"     = "true"
      "MICRO_AUTH"            = "jwt"
    }
  )
}

module "runtime_cert" {
  source = "./cert"

  ca_cert_pem        = tls_self_signed_cert.platform_ca_cert.cert_pem
  ca_private_key_pem = tls_private_key.platform_ca_key.private_key_pem
  private_key_alg    = var.private_key_alg

  subject = local.runtime_name
}

resource "kubernetes_secret" "runtime_cert" {
  metadata {
    name        = "${replace(local.runtime_name, ".", "-")}-cert"
    namespace   = kubernetes_namespace.platform.id
    labels      = local.runtime_labels
    annotations = local.runtime_annotations
  }
  data = {
    "cert.pem" = module.runtime_cert.cert_pem
    "key.pem"  = module.runtime_cert.key_pem
  }
  type = "Opaque"
}

resource "kubernetes_deployment" "runtime" {
  metadata {
    name        = replace(local.runtime_name, ".", "-")
    namespace   = kubernetes_namespace.platform.id
    labels      = local.runtime_labels
    annotations = local.runtime_annotations
  }
  spec {
    replicas = 1
    selector {
      match_labels = local.runtime_labels
    }
    template {
      metadata {
        labels = local.runtime_labels
      }
      spec {
        container {
          name = replace(local.runtime_name, ".", "-")
          dynamic "env" {
            for_each = local.runtime_env
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
          args              = ["runtime"]
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
            secret_name  = kubernetes_secret.runtime_cert.metadata[0].name
          }
        }
        automount_service_account_token = true
        service_account_name            = kubernetes_service_account.runtime.metadata[0].name
      }
    }
  }
}

resource "kubernetes_service_account" "runtime" {
  metadata {
    name        = replace(local.runtime_name, ".", "-")
    namespace   = kubernetes_namespace.platform.id
    labels      = local.runtime_labels
    annotations = local.runtime_annotations
  }
}

resource "random_id" "runtime" {
  byte_length = 3
}

resource "kubernetes_cluster_role" "runtime" {
  metadata {
    name        = "${replace(local.runtime_name, ".", "-")}-${random_id.runtime.hex}"
    labels      = local.runtime_labels
    annotations = local.runtime_annotations
  }
  rule {
    api_groups = [""]
    resources = [
      "pods",
      "services",
      "namespaces",
    ]
    verbs = [
      "create",
      "update",
      "delete",
      "list",
      "patch",
      "watch",
    ]
  }
  rule {
    api_groups = ["apps"]
    resources  = ["deployments"]
    verbs = [
      "create",
      "update",
      "delete",
      "list",
      "patch",
      "watch",
    ]
  }
  rule {
    api_groups = [""]
    resources = [
      "secrets",
      "pods",
      "pods/logs",
    ]
    verbs = [
      "get",
      "list",
      "watch",
    ]
  }
}

resource "kubernetes_cluster_role_binding" "runtime" {
  metadata {
    name        = "${replace(local.runtime_name, ".", "-")}-${random_id.runtime.hex}"
    labels      = local.runtime_labels
    annotations = local.runtime_annotations
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.runtime.metadata[0].name
    namespace = kubernetes_namespace.platform.id
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.runtime.metadata[0].name
  }
}
