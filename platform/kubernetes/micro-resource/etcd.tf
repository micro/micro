locals {
  etcd_labels = { "component" = "etcd" }
}

resource "kubernetes_service" "etcd" {
  metadata {
    name      = "etcd"
    namespace = kubernetes_namespace.resource_namespace.id
    annotations = {
      # Deprecated but still around
      "service.alpha.kubernetes.io/tolerate-unready-endpoints" = "true"
    }
  }
  spec {
    port {
      port = 2379
      name = "client"
    }
    port {
      port = 2380
      name = "peer"
    }
    cluster_ip                  = "None"
    selector                    = local.etcd_labels
    publish_not_ready_addresses = true
  }
}

resource "kubernetes_service" "etcd_cluster" {
  metadata {
    name      = "etcd-cluster"
    namespace = kubernetes_namespace.resource_namespace.id
    labels    = local.etcd_labels
  }
  spec {
    selector = local.etcd_labels
    port {
      name        = "client"
      port        = 2379
      target_port = "client"
    }
  }
}

resource "kubernetes_pod_disruption_budget" "etcd" {
  metadata {
    name      = "etcd"
    namespace = kubernetes_namespace.resource_namespace.id
  }
  spec {
    max_unavailable = "1"
    selector {
      match_labels = local.etcd_labels
    }
  }
}

locals {
  etcd_initial_cluster_string = join(",", formatlist("etcd-%s=http://etcd-%s.etcd:2380", ["0", "1", "2"], ["0", "1", "2"]))
}

resource "random_id" "etcd_cluster_token" {
  byte_length = 16
}

resource "kubernetes_stateful_set" "etcd" {
  metadata {
    name      = "etcd"
    namespace = kubernetes_namespace.resource_namespace.id
    labels    = local.etcd_labels
  }
  spec {
    service_name = "etcd"
    replicas     = var.etcd_replicas
    update_strategy {
      type = "RollingUpdate"
    }
    selector {
      match_labels = local.etcd_labels
    }
    template {
      metadata {
        name   = "etcd"
        labels = local.etcd_labels
      }
      spec {
        container {
          name              = "etcd"
          image             = var.etcd_image
          image_pull_policy = var.image_pull_policy
          command = [
            "/bin/sh",
            "-ecx",
            <<-ETCDLAUNCHER
            exec etcd --name $${POD_NAME} \
              --listen-peer-urls=http://$${POD_IP}:2380 \
              --listen-client-urls=http://$${POD_IP}:2379,http://127.0.0.1:2379 \
              --advertise-client-urls=http://$${POD_NAME}.etcd.$${POD_NAMESPACE}.svc:2379 \
              --initial-advertise-peer-urls=http://$${POD_NAME}.etcd.$${POD_NAMESPACE}.svc:2380 \
              --initial-cluster-token=${random_id.etcd_cluster_token.hex} \
              --initial-cluster=${local.etcd_initial_cluster_string} \
              --data-dir=/var/run/etcd/default.etcd
            ETCDLAUNCHER
          ]
          env {
            name = "POD_IP"
            value_from {
              field_ref {
                field_path = "status.podIP"
              }
            }
          }
          env {
            name = "POD_NAME"
            value_from {
              field_ref {
                field_path = "metadata.name"
              }
            }
          }
          env {
            name = "POD_NAMESPACE"
            value_from {
              field_ref {
                field_path = "metadata.namespace"
              }
            }
          }
          env {
            name  = "ETCDCTL_API"
            value = "3"
          }
          port {
            container_port = 2379
            name           = "client"
          }
          port {
            container_port = 2380
            name           = "peer"
          }
          volume_mount {
            name       = "etcd-data"
            mount_path = "/var/run/etcd"
          }
        }
      }
    }
    volume_claim_template {
      metadata {
        name = "etcd-data"
      }
      spec {
        access_modes = ["ReadWriteOnce"]
        resources {
          requests = {
            "storage" = "10Gi"
          }
        }
      }
    }
  }
}
