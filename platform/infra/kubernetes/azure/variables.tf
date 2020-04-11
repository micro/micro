variable "name" {
  type        = string
  description = "Cluster name"
  default     = "micro"
}

variable "region" {
  type        = string
  description = "Region Code"
  default     = "westeurope"
}

variable "vm_size" {
  type        = string
  description = "Azure VM size"
  default     = "Standard_A2_v2"
}

variable "instance_count" {
  type        = number
  description = "Instance count in default node pool"
  default     = 3
}