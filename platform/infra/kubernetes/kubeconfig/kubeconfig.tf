variable "kubernetes" {
  type        = string
  description = "The name of the Kubernetes module that was used to instantiate kube"
}

variable "args" {
  type        = list(string)
  description = <<-EOD
    arg 0: Cluster remote state key
    arg 1: Cluster remote state region
    With the name of the module and args, the provider should output a kubeconfig file to the module path
  EOD
}

variable "output_path" {
  type        = string
  description = "File to write kubeconfig to (blank for module path)"
  default     = ""
}

provider "local" {
  version = "~> 1.4"
}

provider "digitalocean" {
  version = "~> 1.15"
}

data "digitalocean_kubernetes_cluster" "do_k8s" {
  count = var.kubernetes == "do" ? 1 : 0
  name  = data.terraform_remote_state.k8s.outputs.cluster_name
}

resource "local_file" "do_kubeconfig" {
  count             = var.kubernetes == "do" ? 1 : 0
  sensitive_content = data.digitalocean_kubernetes_cluster.do_k8s[count.index].kube_config.0.raw_config
  filename          = length(var.output_path) > 0 ? var.output_path : "${path.module}/kubeconfig"
  file_permission   = "0600"
}


provider "azurerm" {
  version = "~>2.2"
  features {}
}

data "azurerm_kubernetes_cluster" "aks" {
  count               = var.kubernetes == "azure" ? 1 : 0
  resource_group_name = data.terraform_remote_state.k8s.outputs.cluster_name
  name                = data.terraform_remote_state.k8s.outputs.cluster_name
}

resource "local_file" "aks_kubeconfig" {
  count             = var.kubernetes == "azure" ? 1 : 0
  sensitive_content = data.azurerm_kubernetes_cluster.aks[count.index].kube_config_raw
  filename          = length(var.output_path) > 0 ? var.output_path : "${path.module}/kubeconfig"
  file_permission   = "0600"
}
