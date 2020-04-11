locals {
  athens_labels = { "app" = "athens-proxy" }
}

resource "kubernetes_persistent_volume_claim" "athens" {
  metadata {
    name      = "athens-storage"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels    = local.athens_labels
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        "storage" = var.athens_storage
      }
    }
  }
}

resource "kubernetes_service" "athens" {
  metadata {
    name      = "athens-proxy"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels    = local.athens_labels
  }
  spec {
    type     = "ClusterIP"
    selector = local.athens_labels
    port {
      name        = "http"
      port        = 80
      target_port = "http"
      protocol    = "TCP"
    }
  }
}

resource "kubernetes_deployment" "athens" {
  metadata {
    name      = "athens-proxy"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels    = local.athens_labels
  }
  spec {
    selector {
      match_labels = local.athens_labels
    }
    template {
      metadata {
        labels = local.athens_labels
      }
      spec {
        container {
          name              = "athens-proxy"
          image             = var.athens_image
          image_pull_policy = var.image_pull_policy

          liveness_probe {
            failure_threshold = 3
            http_get {
              path = "/healthz"
              port = "http"
            }
          }

          env {
            name  = "ATHENS_GOGET_WORKERS"
            value = "4"
          }

          env {
            name  = "ATHENS_STORAGE_TYPE"
            value = "disk"
          }

          env {
            name  = "ATHENS_DISK_STORAGE_ROOT"
            value = "/var/lib/athens"
          }

          port {
            name           = "http"
            container_port = 3000
          }

          volume_mount {
            name       = "storage-volume"
            mount_path = "/var/lib/athens"
          }
        }
        volume {
          name = "storage-volume"
          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.athens.metadata.0.name
          }
        }
      }
    }
  }
}
