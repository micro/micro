resource "kubernetes_deployment" "network" {
  metadata {
    namespace   = data.terraform_remote_state.namespaces.outputs.control_namespace
    name        = "micro-network"
    labels      = merge(local.common_labels, { "name" = "micro-network" })
    annotations = merge(local.common_annotations, { "name" = "go.micro.network" })
  }
  spec {
    replicas = var.replicas
    selector {
      match_labels = merge(local.common_labels, { "name" = "micro-network" })
    }
    template {
      metadata {
        labels = merge(local.common_labels, { "name" = "micro-network" })
      }
      spec {
        container {
          name = "micro-network"

          command = [
            "/micro",
            "network",
          ]

          image             = var.micro_image
          image_pull_policy = var.image_pull_policy

          dynamic "env" {
            for_each = local.common_env_vars
            content {
              name  = env.key
              value = env.value
            }
          }
          env {
            name  = "MICRO_SERVER_ADDRESS"
            value = "0.0.0.0:8080"
          }
          env {
            name  = "MICRO_NETWORK_NODES"
            value = "network.${var.domain_name}"
          }

          port {
            container_port = 8080
            name           = "service-port"
          }
          port {
            container_port = 8085
            name           = "network-port"
          }
        }
        affinity {
          pod_anti_affinity {
            preferred_during_scheduling_ignored_during_execution {
              weight = 100
              pod_affinity_term {
                label_selector {
                  match_expressions {
                    key      = "name"
                    operator = "In"
                    values   = ["micro-network"]
                  }
                }
                topology_key = "kubernetes.io/hostname"
              }
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "network" {
  metadata {
    namespace = data.terraform_remote_state.namespaces.outputs.control_namespace
    name      = "micro-network"
    labels    = merge(local.common_labels, { "name" = "micro-network" })
  }
  spec {
    type     = "NodePort"
    selector = merge(local.common_labels, { "name" = "micro-network" })
    port {
      name      = "network"
      port      = 8085
      node_port = 30038
      protocol  = "UDP"
    }
  }
}
