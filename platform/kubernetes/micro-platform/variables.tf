variable "platform_namespace" {
  type        = string
  description = "Namespace containing the micro platform"
  default     = "platform"
}

variable "resource_namespace" {
  type        = string
  description = "Namespace containing shared resources"
  default     = "resource"
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

variable "api_acme_hosts" {
  type        = string
  description = "Comma-separated ACME hosts for micro api"
  default     = "*.m3o.dev,m3o.dev"
}

variable "proxy_acme_hosts" {
  type        = string
  description = "Comma-separated ACME hosts for micro proxy"
  default     = "*.m3o.dev,m3o.dev"
}

variable "web_acme_hosts" {
  type        = string
  description = "Comma-separated ACME hosts for micro web"
  default     = "*.m3o.dev,m3o.dev"
}

variable "cf_api_token" {
  type        = string
  description = "Cloudflare API Token"
}

variable "micro_auth_private" {
  type        = string
  description = "base64 encoded RSA private key PEM for JWTs"
}

variable "micro_auth_public" {
  type        = string
  description = "base64 encoded RSA public key PEM for JWTs"
}

variable "micro_slack_token" {
  type        = string
  description = "Micro slack token"
}
