# This file declares variables declared by micro platform to suppress warnings when deploying both from the same config file

variable "per_service_overrides" {
  default = "UNUSED"
}

variable "external_service_type" {
  default = "UNUSED"
}

variable "platform_namespace" {
  default = "UNUSED"
}

variable "acme_hosts" {
  default = "UNUSED"
}

variable "image_pull_secret" {
  default = "UNUSED"
}

variable "private_key_alg" {
  default = "UNUSED"
}

variable "services" {
  default = "UNUSED"
}

variable "external_services" {
  default = "UNUSED"
}
