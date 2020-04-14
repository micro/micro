# TODO: Pod Security Policy https://github.com/terraform-providers/terraform-provider-kubernetes/pull/624

locals {
  netdata_labels = { "app" = "netdata" }
}

resource "random_uuid" "netdata_stream_id" {}

resource "kubernetes_config_map" "netdata_master" {
  metadata {
    name      = "netdata-conf-master"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels    = local.netdata_labels
  }
  data = {
    "health"  = <<-EOCONF
      SEND_EMAIL="NO"
      SEND_SLACK="YES"
      SLACK_WEBHOOK_URL=""
      DEFAULT_RECIPIENT_SLACK=""
      role_recipients_slack[sysadmin]="$${DEFAULT_RECIPIENT_SLACK}"
      role_recipients_slack[domainadmin]="$${DEFAULT_RECIPIENT_SLACK}"
      role_recipients_slack[dba]="$${DEFAULT_RECIPIENT_SLACK}"
      role_recipients_slack[webmaster]="$${DEFAULT_RECIPIENT_SLACK}"
      role_recipients_slack[proxyadmin]="$${DEFAULT_RECIPIENT_SLACK}"
      role_recipients_slack[sitemgr]="$${DEFAULT_RECIPIENT_SLACK}"
      EOCONF
    "netdata" = <<-EONETDATA
      [global]
        memory mode = save
        bind to = 0.0.0.0:19999
      [plugins]
        cgroups = no
        tc = no
        enable running new plugins = no
        check for new plugins every = 72000
        python.d = no
        charts.d = no
        go.d = no
        node.d = no
        apps = no
        proc = no
        idlejitter = no
        diskspace = no
        micro.d = yes
      [plugin:micro.d]
        command options = --registry=etcd --registry_address=etcd-cluster.${data.terraform_remote_state.namespaces.outputs.resource_namespace}.svc
      EONETDATA
    "stream"  = <<-EOSTREAM
      [${upper(random_uuid.netdata_stream_id.result)}]
        enabled = yes
        history = 3600
        default memory mode = save
        health enabled by default = auto
        allow from = *
      EOSTREAM
  }
}

resource "kubernetes_config_map" "netdata_worker" {
  metadata {
    name      = "netdata-conf-worker"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels    = local.netdata_labels
  }
  data = {
    "netdata"   = <<-EONETDATA
      [global]
        memory mode = none
      [health]
        enabled = no
      [plugins]
        micro.d = no
      EONETDATA
    "stream"    = <<-EOSTREAM
      [stream]
        enabled = yes
        destination = netdata:19999
        api key = ${upper(random_uuid.netdata_stream_id.result)}
        timeout seconds = 60
        buffer size bytes = 1048576
        reconnect delay seconds = 5
        initial clock resync iterations = 60
      EOSTREAM
    "coredns"   = <<-EOCOREDNS
      update_every: 1
      autodetection_retry: 0
      jobs:
        - url: http://127.0.0.1:9153/metrics
      EOCOREDNS
    "kubelet"   = <<-EOKUBELET
      update_every: 1
      autodetection_retry: 0
      jobs:
        - url: http://127.0.0.1:10255/metrics
      EOKUBELET
    "kubeproxy" = <<-EOKUBEPROXY
      update_every: 1
      autodetection_retry: 0
      jobs:
        - url: http://127.0.0.1:10249/metrics
      EOKUBEPROXY
  }
}

resource "kubernetes_service_account" "netdata" {
  metadata {
    name      = "netdata"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels    = local.netdata_labels
  }
}

resource "random_pet" "netdata_cluster_role" {
  prefix    = "netdata"
  separator = "-"
  length    = 2
}

resource "kubernetes_cluster_role" "netdata" {
  metadata {
    name   = random_pet.netdata_cluster_role.id
    labels = local.netdata_labels
  }
  rule {
    api_groups = [""]
    resources  = ["services", "events", "endpoints", "pods", "nodes", "componentstatuses", "nodes/proxy"]
    verbs      = ["get", "list", "watch"]
  }
  rule {
    api_groups = [""]
    resources  = ["resourcequotas"]
    verbs      = ["get", "list"]
  }
  rule {
    api_groups = ["extensions"]
    resources  = ["ingresses"]
    verbs      = ["get", "list", "watch"]
  }
  rule {
    non_resource_urls = ["/version", "/healthz", "/metrics"]
    verbs             = ["get"]
  }
  rule {
    api_groups = [""]
    resources  = ["nodes/metrics", "nodes/spec"]
    verbs      = ["get"]
  }
}

resource "kubernetes_cluster_role_binding" "netdata" {
  metadata {
    name   = random_pet.netdata_cluster_role.id
    labels = local.netdata_labels
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = random_pet.netdata_cluster_role.id
  }
  subject {
    kind      = "ServiceAccount"
    name      = "netdata"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
  }
}

resource "kubernetes_service" "netdata" {
  metadata {
    name      = "netdata"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels    = merge(local.netdata_labels, { "role" = "master" })
  }
  spec {
    type = "ClusterIP"
    port {
      port        = 19999
      target_port = "http"
      protocol    = "TCP"
      name        = "http"
    }
    selector = merge(local.netdata_labels, { "role" = "master" })
  }
}

resource "kubernetes_daemonset" "netdata-worker" {
  metadata {
    name      = "netdata-worker"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels    = merge(local.netdata_labels, { "role" = "worker" })
  }
  spec {
    selector {
      match_labels = merge(local.netdata_labels, { "role" = "worker" })
    }
    template {
      metadata {
        annotations = {
          "container.apparmor.security.beta.kubernetes.io/netdata" = "unconfined"
        }
        labels = merge(local.netdata_labels, { "role" = "worker" })
      }
      spec {
        service_account_name = kubernetes_service_account.netdata.metadata.0.name
        restart_policy       = "Always"
        host_pid             = true
        host_ipc             = true
        host_network         = true
        container {
          name              = "netdata"
          image             = var.netdata_image
          image_pull_policy = var.image_pull_policy
          env {
            name = "MY_POD_NAME"
            value_from {
              field_ref {
                field_path = "metadata.name"
              }
            }
          }
          env {
            name = "MY_POD_NAMESPACE"
            value_from {
              field_ref {
                field_path = "metadata.namespace"
              }
            }
          }
          lifecycle {
            pre_stop {
              exec {
                command = ["/bin/sh", "-c", "killall netdata; while killall -0 netdata; do sleep 1; done"]
              }
            }
          }
          port {
            name           = "http"
            container_port = 19999
            host_port      = 19999
            protocol       = "TCP"
          }
          liveness_probe {
            http_get {
              path = "/api/v1/info"
              port = "http"
            }
            timeout_seconds   = 1
            period_seconds    = 30
            success_threshold = 1
            failure_threshold = 3
          }
          readiness_probe {
            http_get {
              path = "/api/v1/info"
              port = "http"
            }
            timeout_seconds   = 1
            period_seconds    = 30
            success_threshold = 1
            failure_threshold = 3
          }
          volume_mount {
            name       = "proc"
            read_only  = true
            mount_path = "/host/proc"
          }
          volume_mount {
            name       = "run"
            mount_path = "/var/run/docker.sock"
          }
          volume_mount {
            name       = "sys"
            mount_path = "/host/sys"
          }
          volume_mount {
            name       = "config"
            mount_path = "/etc/netdata/go.d/coredns.conf"
            sub_path   = "coredns"
          }
          volume_mount {
            name       = "config"
            mount_path = "/etc/netdata/go.d/k8s_kubelet.conf"
            sub_path   = "kubelet"
          }
          volume_mount {
            name       = "config"
            mount_path = "/etc/netdata/go.d/k8s_kubeproxy.conf"
            sub_path   = "kubeproxy"
          }
          volume_mount {
            name       = "config"
            mount_path = "/etc/netdata/netdata.conf"
            sub_path   = "netdata"
          }
          volume_mount {
            name       = "config"
            mount_path = "/etc/netdata/stream.conf"
            sub_path   = "stream"
          }
          security_context {
            capabilities {
              add = ["SYS_PTRACE", "SYS_ADMIN"]
            }
          }
        }
        toleration {
          effect   = "NoSchedule"
          operator = "Exists"
        }
        volume {
          name = "proc"
          host_path {
            path = "/proc"
          }
        }
        volume {
          name = "run"
          host_path {
            path = "/var/run/docker.sock"
          }
        }
        volume {
          name = "sys"
          host_path {
            path = "/sys"
          }
        }
        volume {
          name = "config"
          config_map {
            name = kubernetes_config_map.netdata_worker.metadata.0.name
          }
        }
        dns_policy = "ClusterFirstWithHostNet"
      }
    }
  }
  depends_on = [kubernetes_config_map.netdata_worker]
}

resource "kubernetes_stateful_set" "netdata_master" {
  metadata {
    name      = "netdata-master"
    namespace = data.terraform_remote_state.namespaces.outputs.resource_namespace
    labels    = merge(local.netdata_labels, { "role" = "master" })
  }
  spec {
    service_name = "netdata"
    replicas     = 1
    selector {
      match_labels = merge(local.netdata_labels, { "role" = "master" })
    }
    template {
      metadata {
        labels = merge(local.netdata_labels, { "role" = "master" })
      }
      spec {
        security_context {
          fs_group = 201
        }
        service_account_name = kubernetes_service_account.netdata.metadata.0.name
        container {
          name              = "netdata"
          image             = var.netdata_image
          image_pull_policy = var.image_pull_policy
          env {
            name = "MY_POD_NAME"
            value_from {
              field_ref {
                field_path = "metadata.name"
              }
            }
          }
          env {
            name = "MY_POD_NAMESPACE"
            value_from {
              field_ref {
                field_path = "metadata.namespace"
              }
            }
          }
          lifecycle {
            pre_stop {
              exec {
                command = ["/bin/sh", "-c", "killall netdata; while killall -0 netdata; do sleep 1; done"]
              }
            }
          }
          port {
            name           = "http"
            container_port = 19999
            protocol       = "TCP"
          }
          liveness_probe {
            http_get {
              path = "/api/v1/info"
              port = "http"
            }
            timeout_seconds   = 1
            period_seconds    = 30
            success_threshold = 1
            failure_threshold = 3
          }
          readiness_probe {
            http_get {
              path = "/api/v1/info"
              port = "http"
            }
            timeout_seconds   = 1
            period_seconds    = 30
            success_threshold = 1
            failure_threshold = 3
          }
          volume_mount {
            name       = "config"
            mount_path = "/etc/netdata/health_alarm_notify.conf"
            sub_path   = "health"
          }
          volume_mount {
            name       = "config"
            mount_path = "/etc/netdata/netdata.conf"
            sub_path   = "netdata"
          }
          volume_mount {
            name       = "config"
            mount_path = "/etc/netdata/stream.conf"
            sub_path   = "stream"
          }
        }
        volume {
          name = "config"
          config_map {
            name = kubernetes_config_map.netdata_master.metadata.0.name
          }
        }
      }
    }
    update_strategy {
      type = "RollingUpdate"
    }
  }
  depends_on = [kubernetes_config_map.netdata_master]
}
