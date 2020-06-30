// Platform namespace
resource "kubernetes_namespace" "platform" {
  metadata {
    name = var.platform_namespace
  }
}

// RSA key pair for signing / verifying JWTs
resource "tls_private_key" "platform_jwt_key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

// Image pull secret for private docker registries
resource "kubernetes_secret" "image_pull_secret" {
  count = length(var.image_pull_secret) > 0 ? 1 : 0
  type  = "kubernetes.io/dockerconfigjson"
  metadata {
    name      = "micro-image-pull-secret"
    namespace = kubernetes_namespace.platform.metadata[0].name
  }
  data = {
    ".dockerconfigjson" = var.image_pull_secret
  }
}