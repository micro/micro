locals {
  cockroachdb_labels = { "app" = "cockroachdb" }
}

resource "kubernetes_service" "cockroachdb_public" {
  metadata {
    name      = "cockroachdb-public"
    namespace = kubernetes_namespace.resource_namespace.id
  }
  spec {
    port {
      name        = "grpc"
      port        = 26257
      target_port = "grpc"
    }
    port {
      name        = "http"
      port        = 8080
      target_port = "http"
    }
    selector = local.cockroachdb_labels
  }
}

resource "kubernetes_service" "cockroachdb" {
  metadata {
    name      = "cockroachdb"
    namespace = kubernetes_namespace.resource_namespace.id

    labels = local.cockroachdb_labels
    annotations = {
      "service.alpha.kubernetes.io/tolerate-unready-endpoints" = "true"
    }
  }
  spec {
    port {
      name        = "grpc"
      port        = 26257
      target_port = "grpc"
    }
    port {
      name        = "http"
      port        = 8080
      target_port = "http"
    }
    publish_not_ready_addresses = true
    cluster_ip                  = "None"
    selector                    = local.cockroachdb_labels
  }
}

resource "kubernetes_pod_disruption_budget" "cockroachdb" {
  metadata {
    name      = "cockroachdb"
    namespace = kubernetes_namespace.resource_namespace.id
    labels    = local.cockroachdb_labels
  }
  spec {
    selector {
      match_labels = local.cockroachdb_labels
    }
    max_unavailable = 1
  }
}

resource "kubernetes_stateful_set" "cockroachdb" {
  metadata {
    name      = "cockroachdb"
    namespace = kubernetes_namespace.resource_namespace.id
  }
  spec {
    service_name = kubernetes_service.cockroachdb.metadata.0.name
    replicas     = var.cockroach_replicas
    selector {
      match_labels = local.cockroachdb_labels
    }
    template {
      metadata {
        labels = local.cockroachdb_labels
      }
      spec {
        affinity {
          pod_anti_affinity {
            preferred_during_scheduling_ignored_during_execution {
              weight = 100
              pod_affinity_term {
                topology_key = "kubernetes.io/hostname"
                label_selector {
                  dynamic "match_expressions" {
                    for_each = local.cockroachdb_labels
                    content {
                      key      = match_expressions.key
                      operator = "In"
                      values   = [match_expressions.value]
                    }
                  }
                }
              }
            }
          }
        }
        container {
          name              = "cockroachdb"
          image             = var.cockroachdb_image
          image_pull_policy = var.image_pull_policy

          port {
            container_port = 26257
            name           = "grpc"
          }
          port {
            container_port = 8080
            name           = "http"
          }

          liveness_probe {
            http_get {
              path = "/health"
              port = "http"
            }
            initial_delay_seconds = 30
            period_seconds        = 5
          }
          readiness_probe {
            http_get {
              path = "/health?ready=1"
              port = "http"
            }
            initial_delay_seconds = 10
            period_seconds        = 5
            failure_threshold     = 2
          }

          volume_mount {
            name       = "datadir"
            mount_path = "/cockroach/cockroach-data"
          }

          env {
            name  = "COCKROACH_CHANNEL"
            value = "kubernetes-insecure"
          }

          command = [
            "/bin/bash",
            "-ecx",
            "exec /cockroach/cockroach start --logtostderr --insecure --advertise-host $(hostname -f) --http-addr 0.0.0.0 --join cockroachdb-0.cockroachdb,cockroachdb-1.cockroachdb,cockroachdb-2.cockroachdb --cache 25% --max-sql-memory 25%"
          ]
        }

        termination_grace_period_seconds = 60

        volume {
          name = "datadir"
          persistent_volume_claim {
            claim_name = "datadir"
          }
        }
      }
    }

    update_strategy {
      type = "RollingUpdate"
    }
    pod_management_policy = "Parallel"

    volume_claim_template {
      metadata {
        name = "datadir"
      }
      spec {
        access_modes = ["ReadWriteOnce"]
        resources {
          requests = {
            "storage" = var.cockroachdb_storage
          }
        }
      }
    }
  }
}

resource "kubernetes_job" "cockroachdb_init" {
  metadata {
    name      = "cockroachdb-init"
    namespace = kubernetes_namespace.resource_namespace.id
  }
  spec {
    selector {
      match_labels = {}
    }
    template {
      metadata {}
      spec {
        active_deadline_seconds = 300
        restart_policy          = "OnFailure"
        init_container {
          name  = "wait-for-cluster"
          image = "alpine:3.9"
          command = [
            "/bin/sh",
            "-c",
            "until wget -O - http://${kubernetes_service.cockroachdb.metadata.0.name}-${kubernetes_stateful_set.cockroachdb.spec.0.replicas - 1}.${kubernetes_service.cockroachdb.metadata.0.name}:${kubernetes_service.cockroachdb.spec.0.port.1.port}/health; do sleep 5; done; echo Waiting 20 seconds; sleep 20"
          ]
        }
        container {
          name              = "cluster-init"
          image             = var.cockroachdb_image
          image_pull_policy = "IfNotPresent"
          command = [
            "/cockroach/cockroach",
            "init",
            "--insecure",
            "--host=${kubernetes_service.cockroachdb.metadata.0.name}-0.${kubernetes_service.cockroachdb.metadata.0.name}"
          ]
        }
      }
    }
  }
}
