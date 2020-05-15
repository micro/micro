locals {
  common_labels = {
    "micro" = "runtime"
  }
  common_annotations = {
    "source"  = "github.com/micro/micro"
    "owner"   = "micro"
    "group"   = "micro"
    "version" = "latest"
  }
  common_env_vars = {
    "MICRO_LOG_LEVEL"         = "debug"
    "MICRO_BROKER"            = "nats"
    "MICRO_BROKER_ADDRESS"    = "nats-cluster.${var.resource_namespace}.svc"
    "MICRO_REGISTRY"          = "etcd"
    "MICRO_REGISTRY_ADDRESS"  = "etcd-cluster.${var.resource_namespace}.svc"
    "MICRO_REGISTER_TTL"      = "60"
    "MICRO_REGISTER_INTERVAL" = "30"
    "MICRO_STORE"             = "cockroach"
    "MICRO_STORE_ADDRESS"     = "postgres://root@cockroachdb-public.${var.resource_namespace}.svc:26257/?sslmode=disable"
  }
}

resource "kubernetes_namespace" "platform" {
  metadata {
    name = kubernetes_namespace.platform.id
  }
}

resource "kubernetes_secret" "cloudflare_credentals" {
  metadata {
    name        = "cloudflare-credentials"
    namespace   = kubernetes_namespace.platform.id
    labels      = local.common_labels
    annotations = local.common_annotations
  }
  data = {
    "CF_API_TOKEN" = var.cf_api_token
  }
}

resource "kubernetes_secret" "micro_keypair" {
  metadata {
    name        = "micro-keypair"
    namespace   = kubernetes_namespace.platform.id
    labels      = local.common_labels
    annotations = local.common_annotations
  }
  data = {
    private = var.micro_auth_private
    public  = var.micro_auth_public
  }
}

resource "kubernetes_secret" "platform_ca" {
  metadata {
    name        = "platform-ca"
    namespace   = kubernetes_namespace.platform.id
    labels      = local.common_labels
    annotations = local.common_annotations
  }
  data = {
    "ca.pem" = tls_self_signed_cert.platform_ca_cert.cert_pem
  }
}

resource "kubernetes_secret" "slack_token" {
  metadata {
    name        = "micro-slack-token"
    namespace   = kubernetes_namespace.platform.id
    labels      = local.common_labels
    annotations = local.common_annotations
  }
  data = {
    "token" = var.micro_slack_token
  }
}

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
