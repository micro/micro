variable "region" {
  description = "AWS Region"
  default     = "eu-west-2"
}

variable "k8s_version" {
  description = "Major+minor Kubernetes version (e.g. 1.15)"
  default     = "1.14"
}

variable "node_count" {
  description = "Number of nodes in the default node pool"
  default     = 3
}

variable "node_flavor" {
  description = "Acceptable vCPU values for nodes"
  default     = "t2.small"
}
