locals {
  ingress_haproxy_labels = {
    "app.kubernetes.io/name"    = "ingress-haproxy"
    "app.kubernetes.io/part-of" = "ingress-haproxy"
  }
}

resource "kubernetes_service_account" "haproxy" {
  metadata {
    name      = "haproxy-ingress"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels    = local.ingress_haproxy_labels
  }
}

resource "random_pet" "haproxy_ingress_cluster_role" {
  prefix    = "haproxy-ingress"
  separator = "-"
  length    = 2
}

resource "kubernetes_cluster_role" "haproxy" {
  metadata {
    name   = random_pet.haproxy_ingress_cluster_role.id
    labels = local.ingress_haproxy_labels
  }
  rule {
    api_groups = [""]
    resources = [
      "configmaps",
      "endpoints",
      "nodes",
      "pods",
      "services",
      "namespaces",
      "events",
      "serviceaccounts",
    ]
    verbs = ["get", "list", "watch"]
  }
  rule {
    api_groups = [""]
    resources  = ["secrets"]
    verbs      = ["get", "list", "watch", "create", "patch", "update"]
  }
  rule {
    api_groups = ["extensions", "networking.k8s.io"]
    resources  = ["ingresses", "ingresses/status"]
    verbs      = ["get", "list", "watch", "update"]
  }
}

resource "kubernetes_cluster_role_binding" "haproxy" {
  metadata {
    name   = random_pet.haproxy_ingress_cluster_role.id
    labels = local.ingress_haproxy_labels
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.haproxy.metadata.0.name
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.haproxy.metadata.0.name
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
  }
}

resource "kubernetes_config_map" "haproxy" {
  metadata {
    name      = "haproxy-ingress"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels    = local.ingress_haproxy_labels
  }
}

resource "kubernetes_deployment" "haproxy_default_backend" {
  metadata {
    name      = "haproxy-default-backend"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels    = merge(local.ingress_haproxy_labels, { "app.kubernetes.io/name" = "default-backend" })
  }
  spec {
    replicas = 1
    selector {
      match_labels = merge(local.ingress_haproxy_labels, { "app.kubernetes.io/name" = "default-backend" })
    }
    template {
      metadata {
        labels = merge(local.ingress_haproxy_labels, { "app.kubernetes.io/name" = "default-backend" })
      }
      spec {
        container {
          name = "default-backend"
          image = "k8s.gcr.io/defaultbackend-amd64:1.5"
          port {
            name = "http"
            container_port = 8080
            protocol = "TCP"
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "haproxy_default_backend" {
  metadata {
    name = "haproxy-default-backend"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels = merge(local.ingress_haproxy_labels, { "app.kubernetes.io/name" = "default-backend" })
  }
  spec {
    selector = merge(local.ingress_haproxy_labels, { "app.kubernetes.io/name" = "default-backend" })
    port {
      name = "http"
      port = 8080
      protocol = "TCP"
      target_port = "http"
    }
  }
}

resource "kubernetes_deployment" "haproxy_ingress" {
  metadata {
    name = "haproxy-ingress"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels = local.ingress_haproxy_labels
  }
  spec {
    selector {
      match_labels = local.ingress_haproxy_labels
    }
    template {
      metadata {
        labels = local.ingress_haproxy_labels
      }
      spec {
        automount_service_account_token = true
        service_account_name = kubernetes_service_account.haproxy.metadata.0.name
        container {
          name = "haproxy"
          image = "haproxytech/kubernetes-ingress"
          args = [
            "--configmap=${kubernetes_config_map.haproxy.id}",
            "--default-backend-service=${data.terraform_remote_state.namespaces.outputs.resource_namespace}/${kubernetes_service.haproxy_default_backend.metadata.0.name}",
            "--publish-service=${kubernetes_service.haproxy_ingress.id}",
          ]
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
          liveness_probe {
            http_get {
              path = "/healthz"
              port = 1042
            }
          }
          port {
            name = "http"
            container_port = 80
          }
          port {
            name = "https"
            container_port = 443
          }
          port {
            name = "stat"
            container_port = 1024
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "haproxy_ingress" {
  metadata {
    name = "haproxy-ingress"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels = local.ingress_haproxy_labels
  }
  spec {
    selector = local.ingress_haproxy_labels
    type = "LoadBalancer"
    port {
      name = "http"
      protocol = "TCP"
      port = 80
      target_port = "http"
    }
    port {
      name = "https"
      protocol = "TCP"
      port = 443
      target_port = "https"
    }
  }
  lifecycle {
    ignore_changes = [metadata.0.annotations]
  }
}
