variable "platform_namespace" {
  type        = string
  description = "Namespace containing the micro platform"
}

variable "resource_namespace" {
  type        = string
  description = "Namespace containing shared resources"
}

variable "micro_image" {
  type        = string
  description = "Micro docker image"
  default     = "micro/micro:latest"
}

variable "image_pull_policy" {
  type        = string
  description = "Kubernetes Image pull policy"
  default     = "Always"
}

variable "private_key_alg" {
  type        = string
  description = "Private Key algorithm for platform CA"
  default     = "ECDSA"
}

variable "acme_hosts" {
  type        = string
  description = "Comma-separated ACME hosts for micro api"
  default     = "*.m3o.dev,m3o.dev"
}

variable "cf_api_token" {
  type        = string
  description = "Cloudflare API Token"
  default     = ""
}

variable "micro_slack_token" {
  type        = string
  description = "Micro slack token"
  default     = ""
}

// Specify image_pull_secret as a json string:
// Usually thes contain auth, or user/password/email
# {
#   "auths": {
#     "REGISTRY_SERVER": {
#       "username": "DOCKER_LOGIN_USERNAME",
#       "password": "DOCKER_LOGIN_PASSWORD",
#       "email": "DUMMY_DOCKER_EMAIL",
#       "auth": "AUTH_TOKEN"
#     }
#   }
# }

variable "image_pull_secret" {
  type        = string
  description = "Image pull secret credentials"
  default     = ""
}
