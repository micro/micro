resource "tls_private_key" "cert" {
  algorithm   = var.private_key_alg
  rsa_bits    = var.private_key_alg == "RSA" ? 4096 : null
  ecdsa_curve = var.private_key_alg == "ECDSA" ? "P384" : null
}

resource "tls_cert_request" "cert" {
  key_algorithm   = var.private_key_alg
  private_key_pem = tls_private_key.cert.private_key_pem

  subject {
    common_name         = var.subject
    organization        = var.organization
    organizational_unit = var.ou
  }

}

resource "tls_locally_signed_cert" "cert" {
  cert_request_pem      = tls_cert_request.cert.cert_request_pem
  ca_key_algorithm      = var.private_key_alg
  ca_private_key_pem    = var.ca_private_key_pem
  ca_cert_pem           = var.ca_cert_pem
  validity_period_hours = 87600
  allowed_uses          = var.allowed_uses
  is_ca_certificate     = var.is_ca_cert
}

output "key_pem" {
  value = tls_private_key.cert.private_key_pem
}

output "cert_pem" {
  value = tls_locally_signed_cert.cert.cert_pem
}
