terraform {
  required_version = ">= 0.12.0"
}

provider "digitalocean" {
  version = "~> 1.13"
}

provider "local" {
  version = "~> 1.4"
}

provider "random" {
  version = "~> 2.2"
}
