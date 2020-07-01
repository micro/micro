variable "resource_namespace" {
  description = "Namespace name to create"
  type        = string
}

variable "image_pull_policy" {
  description = "Kubernetes image pull policy for control plane deployments"
  default     = "Always"
}

variable "micro_image" {
  description = "Micro docker image"
  default     = "micro/micro"
}

variable "etcd_image" {
  description = "etcd docker image"
  default     = "gcr.io/etcd-development/etcd:v3.3.18"
}

variable "nats_image" {
  description = "nats-io docker image"
  default     = "nats:2.1.0-alpine3.10"
}

variable "cockroachdb_image" {
  description = "CockroachDB Image"
  default     = "cockroachdb/cockroach:v19.2.1"
}

variable "cockroachdb_storage" {
  description = "CockroachDB Kubernetes storage request"
  default     = "10Gi"
}

variable "jaeger_image" {
  description = "Jaeger Tracing All in one image"
  default     = "jaegertracing/all-in-one"
}

variable "etcd_replicas" {
  type = number
  description = "number of etcd replicas to deploy"
  default = 3
}

variable "cockroach_replicas" {
  type = number
  description = "number of cockroach replicas to deploy"
  default = 3
}

variable "nats_replicas" {
  type = number
  description = "number of nats replicas to deploy"
  default = 3
}
