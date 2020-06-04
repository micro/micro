variable "ca_private_key_pem" {
  type        = string
  description = "CA Private Key in PEM format"
}

variable "ca_cert_pem" {
  type        = string
  description = "CA Certificate in PEM format"
}

variable "private_key_alg" {
  type        = string
  description = "TLS Private Key alg"
}

variable "subject" {
  type        = string
  description = "Certificate Subject"
}

variable "organization" {
  type        = string
  description = "Certificate Organization"
  default     = "micro"
}

variable "ou" {
  type        = string
  description = "Certificate organizational unit"
  default     = "platform"
}

variable "allowed_uses" {
  type        = list(string)
  description = "Key allowed uses"
  default = [
    "digital_signature",
    "key_encipherment",
    "client_auth",
    "server_auth"
  ]
}

variable "is_ca_cert" {
  type        = bool
  description = "is CA certificate"
  default     = false
}
