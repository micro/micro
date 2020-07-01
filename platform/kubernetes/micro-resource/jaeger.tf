locals {
  jaeger_labels = { "app" = "jaeger" }
}

resource "kubernetes_deployment" "jaeger" {
  metadata {
    name      = "jaeger-tracing"
    namespace = kubernetes_namespace.resource_namespace.id
    labels    = local.jaeger_labels
  }
  spec {
    replicas = 1
    strategy {
      type = "Recreate"
    }
    selector {
      match_labels = local.jaeger_labels
    }
    template {
      metadata {
        labels = local.jaeger_labels
      }
      spec {
        container {
          name = "jaeger"
          env {
            name  = "COLLECTOR_ZIPKIN_HTTP_PORT"
            value = "9411"
          }
          image             = var.jaeger_image
          image_pull_policy = var.image_pull_policy
          port {
            name           = "agent-zip-thrft"
            container_port = 5775
            protocol       = "UDP"
          }
          port {
            name           = "agent-compact"
            container_port = 6831
            protocol       = "UDP"
          }
          port {
            name           = "agent-binary"
            container_port = 6832
            protocol       = "UDP"
          }
          port {
            name           = "agent-configs"
            container_port = 5778
            protocol       = "TCP"
          }
          port {
            name           = "query-http"
            container_port = 16686
            protocol       = "TCP"
          }
          port {
            name           = "coll-zipkin"
            container_port = 9411
            protocol       = "TCP"
          }
          port {
            name           = "health"
            container_port = 14269
            protocol       = "TCP"
          }
          port {
            name           = "coll-tchan"
            container_port = 14267
            protocol       = "TCP"
          }
          port {
            name           = "collector-http"
            container_port = 14268
            protocol       = "TCP"
          }
          readiness_probe {
            http_get {
              path = "/"
              port = "health"
            }
            initial_delay_seconds = 5
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "jaeger_query" {
  metadata {
    name      = "jaeger-tracing-query"
    namespace = kubernetes_namespace.resource_namespace.id
    labels    = local.jaeger_labels
  }
  spec {
    port {
      name        = "query-http"
      port        = 80
      protocol    = "TCP"
      target_port = "query-http"
    }
    type     = "ClusterIP"
    selector = local.jaeger_labels
  }
}

resource "kubernetes_service" "jaeger_collector" {
  metadata {
    name      = "jaeger-tracing-collector"
    namespace = kubernetes_namespace.resource_namespace.id
  }
  spec {
    selector = local.jaeger_labels
    port {
      name        = "tchannel"
      port        = 14267
      protocol    = "TCP"
      target_port = "coll-tchan"
    }
    port {
      name        = "http"
      port        = 14268
      protocol    = "TCP"
      target_port = "collector-http"
    }
    port {
      name        = "zipkin"
      port        = 9411
      protocol    = "TCP"
      target_port = "coll-zipkin"
    }
  }
}

resource "kubernetes_service" "jaeger_agent" {
  metadata {
    name      = "jaeger-tracing-agent"
    namespace = kubernetes_namespace.resource_namespace.id
    labels    = local.jaeger_labels
  }
  spec {
    selector = local.jaeger_labels
    port {
      name        = "zipkin-thrift"
      port        = 5775
      protocol    = "UDP"
      target_port = "agent-zip-thrft"
    }
    port {
      name        = "compact"
      port        = 6831
      protocol    = "UDP"
      target_port = "agent-compact"
    }
    port {
      name        = "binary"
      port        = 6832
      protocol    = "UDP"
      target_port = "agent-binary"
    }
    port {
      name        = "configs"
      port        = 5778
      protocol    = "TCP"
      target_port = "agent-configs"
    }

  }
}

resource "kubernetes_service" "zipkin" {
  metadata {
    name      = "zipkin"
    namespace = kubernetes_namespace.resource_namespace.id
    labels    = local.jaeger_labels
  }
  spec {
    selector = local.jaeger_labels
    port {
      name        = "coll-zipkin"
      port        = 9411
      protocol    = "TCP"
      target_port = "coll-zipkin"
    }
    cluster_ip = "None"
  }
}
