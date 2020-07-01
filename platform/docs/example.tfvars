# k8s namespace to create to contain the shared resources (etcd, cockroach, nats)
resource_namespace = "example-shared"

# k8s namesapce to deploy the m3o platform 
platform_namespace = "example-platform"

# Per-service env overrides:
per_service_overrides = {
  "runtime" = { "MICRO_RUNTIME" = "kubernetes" },
}

## Everything else has defaults that can be overridden by uncommenting them below:

## M3o services configuration

# which services to deploy and their ports:
# services = {
#   "config"   = 8080,
#   "auth"     = 8010,
#   "network"  = 8085,
#   "runtime"  = 8088,
#   "registry" = 8000,
#   "broker"   = 8001,
#   "store"    = 8002,
#   "router"   = 8084,
#   "debug"    = 8080,
#   "proxy"    = 443,
#   "api"      = 443,
#   "web"      = 443,
# }

# which services to expose
# external_services = ["api", "proxy", "web"]

# How to expose them (NodePort / LoadBalancer)
# external_service_type = "NodePort"

## Docker image configuration

# k8s image pull policy
# image_pull_policy = "Always"

# docker config json used to pull private images
# image_pull_secret = "{}"

# Micro docker image
# micro_image  = "micro/micro:latest"

# etcd docker image
# etcd_image = "gcr.io/etcd-development/etcd:v3.3.18"

# nats docker image
# nats_image = "nats:2.1.0-alpine3.10"

# cockroachdb image
# cockroachdb_image = "cockroachdb/cockroach:v19.2.1"

# cockroachdb PV request size
# cockroachdb_storage = "10Gi"

# jaeger image
# jaeger_image = "jaegertracing/all-in-one:latest"

## TLS configuration

# acme_hosts = "*.m3o.dev,m3o.dev"
# private_key_alg = "ECDSA"

## Shared resources replicas
# etcd_replicas = 3
# cockroach_replicas = 3
# nats_replicas = 3
