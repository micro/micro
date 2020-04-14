provider "azurerm" {
  version = "~> 2.2"
  features {}
}

provider "azuread" {
  version = "~> 0.8"
}

provider "random" {
  version = "~> 2.2"
}

resource "random_id" "k8s_name" {
  byte_length = 2
}

resource "azuread_application" "k8s" {
  name = "${var.name}-${var.region}-k8s-${random_id.k8s_name.hex}"

  // hack for consistency
  provisioner "local-exec" {
    command = "sleep 10"
  }
}

resource "azuread_service_principal" "k8s" {
  application_id = azuread_application.k8s.application_id

  // hack for consistency
  provisioner "local-exec" {
    command = "sleep 10"
  }
}

resource "random_password" "service_principal_secret" {
  length  = 32
  special = false
}

resource "azuread_service_principal_password" "k8s" {
  service_principal_id = azuread_service_principal.k8s.id
  value                = random_password.service_principal_secret.result
  end_date_relative    = "87600h"

  // hack for consistency
  provisioner "local-exec" {
    command = "sleep 10"
  }
}

resource "azurerm_resource_group" "k8s" {
  name     = "${var.name}-${var.region}-${random_id.k8s_name.hex}"
  location = var.region
}

resource "azurerm_kubernetes_cluster" "k8s_cluster" {
  name                = "${var.name}-${var.region}-${random_id.k8s_name.hex}"
  location            = azurerm_resource_group.k8s.location
  resource_group_name = azurerm_resource_group.k8s.name
  dns_prefix          = "${var.name}-${var.region}-${random_id.k8s_name.hex}"

  addon_profile {
    kube_dashboard {
      enabled = false
    }
  }

  default_node_pool {
    name       = "default${random_id.k8s_name.dec}"
    node_count = var.instance_count
    vm_size    = var.vm_size
  }

  service_principal {
    client_id     = azuread_service_principal.k8s.application_id
    client_secret = azuread_service_principal_password.k8s.value
  }

}

output "cluster_name" {
  value = azurerm_kubernetes_cluster.k8s_cluster.name
}

output "kubeconfig" {
  value     = azurerm_kubernetes_cluster.k8s_cluster.kube_config_raw
  sensitive = true
}
