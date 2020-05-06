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

variable "netdata_image" {
  description = "Micro customised netdata image"
  default     = "micro/netdata:latest"
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

variable "athens_image" {
  description = "Athens Go Module Proxy image"
  default     = "gomods/athens:v0.7.2"
}

variable "athens_storage" {
  description = "Athens Go Mpdule Proxy Kubernetes storage request"
  default     = "10Gi"
}

variable "nginx_ingress_image" {
  description = "nginx ingress controller image"
  default     = "quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.30.0"
}

variable "in_aws" {
  description = "Are you deploying into an AWS Snowflake env?"
  type        = bool
  default     = false
}
