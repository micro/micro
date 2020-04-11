# This file spins up micro in digitalocean, as a test of the infrastructure as code.

# Instatiate a specific digitalocean provider per module.
provider "digitalocean" {
  alias   = "e2etest"
  token   = var.do_token
  version = "~> 1.13"
}

module "e2etest_k8s" {
  source = "./infrastructure/kubernetes/do"
  providers = {
    digitalocean = digitalocean.e2etest
  }

  region      = "lon1"
  k8s_version = "1.16"
  node_count  = 3
  node_cpu    = [2]
  node_memory = [4096]
}

data "digitalocean_kubernetes_cluster" "e2etest" {
  provider = digitalocean.e2etest
  name     = module.e2etest_k8s.cluster_name
}

provider "kubernetes" {
  alias            = "e2etest"
  load_config_file = false
  host             = data.digitalocean_kubernetes_cluster.e2etest.endpoint
  token            = var.do_token
  #token            = data.digitalocean_kubernetes_cluster.e2etest.kube_config[0].token
  cluster_ca_certificate = base64decode(
    data.digitalocean_kubernetes_cluster.e2etest.kube_config[0].cluster_ca_certificate
  )
}

module "e2etest_resource" {
  source = "./resource"
  providers = {
    kubernetes = kubernetes.e2etest
  }

  resource_namespace = "e2eresource"
}

module "e2etest_control" {
  source = "./control"
  providers = {
    kubernetes = kubernetes.e2etest
  }

  control_namespace  = "e2econtrol"
  resource_namespace = "e2eresource"
  slack_token        = var.slack_token
}

module "e2etest_network" {
  source = "./network"
  providers = {
    kubernetes = kubernetes.e2etest
  }
  network_namespace  = "e2enetwork"
  resource_namespace = "e2eresource"

  cloudflare_account_id              = var.cloudflare_account_id
  cloudflare_api_token               = var.cloudflare_api_token
  cloudflare_dns_zone_id             = var.cloudflare_dns_zone_id
  cloudflare_kv_namespace_id         = var.cloudflare_kv_namespace_id
  cloudflare_kv_namespace_id_runtime = var.cloudflare_kv_namespace_id_runtime
}

resource "local_file" "e2etest_kubeconfig" {
  content  = data.digitalocean_kubernetes_cluster.e2etest.kube_config[0].raw_config
  filename = "/tmp/e2ekubeconfig"
}
