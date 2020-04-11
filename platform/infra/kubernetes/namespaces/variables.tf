variable "resource_namespace" {
  description = "Namespace to deploy shared resources"
  type        = string
}

variable "network_namespace" {
  description = "Namespace to deploy network"
  type        = string
}

variable "control_namespace" {
  description = "Namespace to deploy control place"
  type        = string
}
