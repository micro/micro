variable "replicas" {
  type = number
  description = "Replicas of control plane deployments"
  default     = 1
}

variable "micro_image" {
  description = "Micro docker image"
  default     = "micro/micro"
}

variable "image_pull_policy" {
  type = string
  description = "Kubernetes image pull policy for control plane deployments"
  default     = "Always"
}

variable "domain_name" {
  type = string
  description = "Domain name of the platform (e.g. micro.mu)"
}

variable "slack_token" {
  type = string
  description = "Slack token for micro bot"
}
