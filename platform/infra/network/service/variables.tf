variable "service_name" {
  description = "Service name"
  type        = string
}

variable "service_replicas" {
  description = "Number of replicas"
  type        = number
  default     = 1
}

variable "create_k8s_service" {
  description = "Whether to create a Kubernetes service"
  type        = bool
  default     = true
}

variable "create_k8s_ingress" {
  description = "Whether to create a Kubernetes ingress"
  type        = bool
  default     = false
}

variable "service_type" {
  description = "Kubernetes service type"
  type        = string
  default     = "ClusterIP"
}

variable "service_port" {
  description = "Service port"
  type        = number
  default     = -1
}

variable "service_protocol" {
  description = "Service protocol"
  type        = string
  default     = "TCP"
}

variable "extra_annotations" {
  description = "Additonal Kubernetes annotations"
  type        = map(string)
  default     = {}
}

variable "extra_labels" {
  description = "Additional Kubernetes labels"
  type        = map(string)
  default     = {}
}

variable "extra_env_vars" {
  description = "Additonal environment variables"
  type        = map(string)
  default     = {}
}

variable "network_namespace" {
  description = "Network components namespace"
  type        = string
  default     = "network"
}

variable "resource_namespace" {
  description = "Shared Infrastructure components namespace"
  type        = string
  default     = "resource"
}

variable "micro_image" {
  description = "Micro Docker image"
  type        = string
  default     = "micro/micro:latest"
}

variable "image_pull_policy" {
  description = "Kubernetes image pull policy"
  type        = string
  default     = "Always"
}

variable "domain_names" {
  description = "List of valid domain names for network services"
  type        = list(string)
  default = []
}
