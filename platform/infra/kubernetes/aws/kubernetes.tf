locals {
  cluster_name = "micro-${var.region}-${random_id.k8s_name.hex}"
}

data "aws_availability_zones" "available" {}

resource "random_id" "k8s_name" {
  byte_length = 4
}

module vpc {
  source  = "terraform-aws-modules/vpc/aws"
  version = "2.6.0"

  name                 = "micro-${random_id.k8s_name.hex}"
  cidr                 = "10.0.0.0/16"
  azs                  = data.aws_availability_zones.available.names
  private_subnets      = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets       = ["10.0.4.0/24", "10.0.5.0/24", "10.0.6.0/24"]
  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true

  tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
  }

  private_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
    "kubernetes.io/role/internal-elb"             = "1"
  }

  public_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                      = "1"
  }
}

resource "aws_security_group" "nodes_mgmt" {
  name_prefix = "micro-mgmt-${random_id.k8s_name.hex}"
  vpc_id      = module.vpc.vpc_id
  ingress {
    from_port = 22
    to_port   = 22
    protocol  = "tcp"

    cidr_blocks = [
      "10.0.0.0/8",
      "172.16.0.0/12",
      "192.168.0.0/16",
    ]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }
}

data "aws_eks_cluster" "cluster" {
  name = module.eks.cluster_id
}

data "aws_eks_cluster_auth" "cluster" {
  name = module.eks.cluster_id
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "8.1.0"

  cluster_name    = local.cluster_name
  cluster_version = var.k8s_version
  subnets         = module.vpc.private_subnets
  vpc_id          = module.vpc.vpc_id

  tags = {
    Environment = "micro-aws-${var.region}"
  }
  worker_groups = [{
    name                          = "nodes"
    instance_type                 = var.node_flavor
    asg_desired_capacity          = var.node_count
    additional_security_group_ids = [aws_security_group.nodes_mgmt.id]
  }]
}

# Cluster ID is output for later use to configure a kubernetes provider
output "cluster_name" {
  value = module.eks.cluster_id
}
