resource "random_id" "kv_identifier" {
  byte_length = 4
}

resource "cloudflare_workers_kv_namespace" "micro" {
  title = "micro-${random_id.kv_identifier.hex}"
}

resource "cloudflare_workers_kv_namespace" "runtime" {
  title = "micro-runtime-${random_id.kv_identifier.hex}"
}

output "kv_namespace_id" {
  value = cloudflare_workers_kv_namespace.micro.id
}

output "kv_namespace_id_runtime" {
  value = cloudflare_workers_kv_namespace.runtime.id
}
