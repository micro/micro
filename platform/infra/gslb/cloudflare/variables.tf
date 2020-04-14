variable "micro_ingress_address" {
  description = <<-EOD
  An object describing the Kubernetes ingress service address.
  For Example: 
    {
      "ip4" = "10.20.30.40"
      "ip6" = "2001:db8::1"
      "hostname" = "aws-ingress.amazonaws.com"
    }
  At least one attribute must be set, setting multiple creates more records
  EOD
  type = object({
    ip4      = string,
    ip6      = string,
    hostname = string,
  })
  default = {
    // ip4   = "10.20.30.40",
    // ip6   = "2001:db8::1",
    ip4      = ""
    ip6      = ""
    hostname = "example.com",
  }
}

variable "global_domain_name" {
  description = "Domain name under which to create global load balancers"
  type        = string
  default     = "mu.network"
}

variable "regional_subdomain" {
  description = "Domain name under which to create regional DNS entries"
  type        = string
  default     = "cloud"
}

variable "region" {
  description = "Human-readable region code (usually set to the same as whatever cloud provider you're using)"
  type        = string
  default     = "lon1"
}

variable "nearest_cloudflare_region" {
  description = "Nearest cloudflare region for geo-steering https://developers.cloudflare.com/load-balancing/understand-basics/traffic-steering/"
  type        = string
  default     = "WEU"
}
