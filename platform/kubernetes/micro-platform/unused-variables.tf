# This file declares variables declared by the shared resources to suppress warnings when deploying both from the same config file
variable "etcd_image" {
  default = "UNUSED"
}

variable "cockroachdb_image" {
  default = "UNUSED"
}

variable "nats_image" {
  default = "UNUSED"
}

variable "jaeger_image" {
  default = "UNUSED"
}

variable "etcd_replicas" {
  default = "UNUSED"
}

variable "cockroach_replicas" {
  default = "UNUSED"
}

variable "nats_replicas" {
  default = "UNUSED"
}

variable "cockroachdb_storage" {
  default = "UNUSED"
}
