terraform {
  required_version = ">= 0.12.0"
}

provider "aws" {
  region  = var.region
  version = "~> 2.45"
}

provider "random" {
  version = "~> 2.2"
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  token                  = data.aws_eks_cluster_auth.cluster.token
  load_config_file       = false
}
