// CA for mTLS
resource "tls_private_key" "platform_ca_key" {
  algorithm   = var.private_key_alg
  rsa_bits    = var.private_key_alg == "RSA" ? 4096 : null
  ecdsa_curve = var.private_key_alg == "ECDSA" ? "P384" : null
}

resource "tls_self_signed_cert" "platform_ca_cert" {
  key_algorithm   = var.private_key_alg
  private_key_pem = tls_private_key.platform_ca_key.private_key_pem

  subject {
    common_name  = "Micro Platform"
    organization = "Micro"
  }

  validity_period_hours = 876000

  allowed_uses = [
    "cert_signing",
    "crl_signing",
    "client_auth",
    "server_auth",
  ]
  is_ca_certificate = true
}

resource "tls_private_key" "service_keys" {
  for_each = var.services

  algorithm   = var.private_key_alg
  rsa_bits    = var.private_key_alg == "RSA" ? 4096 : null
  ecdsa_curve = var.private_key_alg == "ECDSA" ? "P521" : null
}

resource "tls_cert_request" "service_csrs" {
  for_each = var.services

  key_algorithm   = var.private_key_alg
  private_key_pem = tls_private_key.service_keys[each.key].private_key_pem

  subject {
    common_name         = replace(each.key, " ", "-")
    organization        = "micro"
    organizational_unit = "platform"
  }
}

resource "tls_locally_signed_cert" "service_certs" {
  for_each = var.services

  cert_request_pem      = tls_cert_request.service_csrs[each.key].cert_request_pem
  ca_key_algorithm      = var.private_key_alg
  ca_private_key_pem    = tls_private_key.platform_ca_key.private_key_pem
  ca_cert_pem           = tls_self_signed_cert.platform_ca_cert.cert_pem
  validity_period_hours = 87600
  allowed_uses = [
    "digital_signature",
    "key_encipherment",
    "client_auth",
    "server_auth"
  ]
  is_ca_certificate = false
}

resource "kubernetes_secret" "cert_bundles" {
  for_each = var.services

  type = "Opaque"
  metadata {
    name      = "${replace(each.key, " ", "-")}-certs"
    namespace = kubernetes_namespace.platform.metadata[0].name
    labels = {
      "micro" = replace(each.key, " ", "-")
    }
  }
  data = {
    "ca.pem"                                  = tls_self_signed_cert.platform_ca_cert.cert_pem,
    "${replace(each.key, " ", "-")}-cert.pem" = tls_locally_signed_cert.service_certs[each.key].cert_pem,
    "${replace(each.key, " ", "-")}-key.pem"  = tls_private_key.service_keys[each.key].private_key_pem,
  }
}
