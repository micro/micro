variable "do_token" {
  description = "DigitalOcean Personal Access Token"
}

variable "slack_token" {
  description = "Slack token for Micro Bot"
}

variable "cloudflare_account_id" {
  description = "Cloudflare Account ID (For connecting to the Cloudflare API)"
  type        = string
}

variable "cloudflare_api_token" {
  description = "Cloudflare API token (For connecting to the Cloudflare API)"
  type        = string
}

variable "cloudflare_kv_namespace_id" {
  description = "Cloudflare workers KV namespace ID"
  type        = string
}

variable "cloudflare_kv_namespace_id_runtime" {
  description = "Cloudflare workers KV namespace ID (runtime)"
  type        = string
}

variable "cloudflare_dns_zone_id" {
  description = "Cloudflare DNS Zone ID"
  type        = string
}
