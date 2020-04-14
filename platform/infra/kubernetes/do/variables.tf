variable "name" {
  description = "Cluster Name"
  default     = "micro"
}

variable "region" {
  description = "DigitalOcean region"
  default     = "lon1"
}

variable "k8s_version" {
  description = "Major+minor Kubernetes version (e.g. 1.16)"
  default     = "1.16"
}

variable "node_count" {
  description = "Number of nodes in the default node pool"
  default     = 3
}

variable "node_cpu" {
  description = "Acceptable vCPU values for nodes"
  default     = [2]
}

variable "node_memory" {
  description = "Acceptable memory values for nodes (MiB)"
  default     = [4096]
}
