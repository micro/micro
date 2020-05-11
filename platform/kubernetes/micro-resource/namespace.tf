resource "kubernetes_namespace" "resource_namespace" {
  metadata {
    name = var.resource_namespace
  }
}