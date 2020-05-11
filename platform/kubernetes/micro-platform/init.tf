locals {
  init_name = "go.micro.init"
  init_labels = merge(
    local.common_labels,
    {
      "name" = local.init_name
    }
  )
  init_annotations = merge(
    local.common_annotations,
    {
      "name" = local.init_name
    }
  )
  init_env = merge(
    {
      "MICRO_NAMESPACE" = "micro"
      "MICRO_RUNTIME"   = "kubernetes"
      "MICRO_LOG_LEVEL" = "debug"
    }
  )
}

module "init_cert" {
  source = "./cert"

  ca_cert_pem        = tls_self_signed_cert.platform_ca_cert.cert_pem
  ca_private_key_pem = tls_private_key.platform_ca_key.private_key_pem
  private_key_alg    = var.private_key_alg

  subject = local.init_name
}

resource "kubernetes_secret" "init_cert" {
  metadata {
    name        = "${replace(local.init_name, ".", "-")}-cert"
    namespace   = var.platform_namespace
    labels      = local.init_labels
    annotations = local.init_annotations
  }
  data = {
    "cert.pem" = module.init_cert.cert_pem
    "key.pem"  = module.init_cert.key_pem
  }
  type = "Opaque"
}

resource "kubernetes_deployment" "init" {
  metadata {
    name        = replace(local.init_name, ".", "-")
    namespace   = var.platform_namespace
    labels      = local.init_labels
    annotations = local.init_annotations
  }
  spec {
    replicas = 1
    selector {
      match_labels = local.init_labels
    }
    template {
      metadata {
        labels = local.init_labels
      }
      spec {
        container {
          name = replace(local.init_name, ".", "-")
          dynamic "env" {
            for_each = local.init_env
            content {
              name  = env.key
              value = env.value
            }
          }
          args              = ["init"]
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
            secret_name  = kubernetes_secret.init_cert.metadata[0].name
          }
        }
        automount_service_account_token = true
        service_account_name            = kubernetes_service_account.init.metadata[0].name
      }
    }
  }
}

resource "kubernetes_service_account" "init" {
  metadata {
    name        = replace(local.init_name, ".", "-")
    namespace   = var.platform_namespace
    labels      = local.init_labels
    annotations = local.init_annotations
  }
}

resource "random_id" "init" {
  byte_length = 3
}

resource "kubernetes_cluster_role" "init" {
  metadata {
    name        = "${replace(local.init_name, ".", "-")}-${random_id.init.hex}"
    labels      = local.init_labels
    annotations = local.init_annotations
  }
  rule {
    api_groups = [""]
    resources = [
      "pods"
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
}

resource "kubernetes_role_binding" "init" {
  metadata {
    name        = "${replace(local.init_name, ".", "-")}-${random_id.init.hex}"
    namespace   = var.platform_namespace
    labels      = local.init_labels
    annotations = local.init_annotations
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.init.metadata[0].name
  }
  subject {
    kind = "ServiceAccount"
    name = kubernetes_service_account.init.metadata[0].name
  }
}