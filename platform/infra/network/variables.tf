variable "cloudflare_account_id" {
  description = "Cloudflare Account ID (For connecting to the Cloudflare API)"
  type        = string
}

variable "cloudflare_api_token" {
  description = "Cloudflare API token (For connecting to the Cloudflare API)"
  type        = string
}

variable "cloudflare_dns_zone_id" {
  description = "Cloudflare DNS Zone ID"
  type        = string
}

variable "domain_name" {
  description = "platform domain name"
  type        = string
}

variable "region_slug" {
  description = "e.g. lon1-do would make api-lon1-do.cloud.(domain)"
  type = string
}
