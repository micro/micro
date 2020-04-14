provider "cloudflare" {
  version = "~> 2.3"
}

data "cloudflare_zones" "micro" {
  filter {
    name   = var.global_domain_name
    status = "active"
  }
}

resource "cloudflare_record" "regional_api_a" {
  count   = length(var.micro_ingress_address.ip4) > 0 ? 1 : 0
  zone_id = data.cloudflare_zones.micro.zones[0].id
  name    = "api-${var.region}.${var.regional_subdomain}"
  value   = var.micro_ingress_address.ip4
  type    = "A"
  ttl     = 1
  proxied = false
}

resource "cloudflare_record" "regional_api_aaaa" {
  count   = length(var.micro_ingress_address.ip6) > 0 ? 1 : 0
  zone_id = data.cloudflare_zones.micro.zones[0].id
  name    = "api-${var.region}.${var.regional_subdomain}"
  value   = var.micro_ingress_address.ip6
  type    = "AAAA"
  ttl     = 1
  proxied = false
}

resource "cloudflare_record" "regional_api_cname" {
  count   = length(var.micro_ingress_address.hostname) > 0 ? 1 : 0
  zone_id = data.cloudflare_zones.micro.zones[0].id
  name    = "api-${var.region}.${var.regional_subdomain}"
  value   = var.micro_ingress_address.hostname
  type    = "CNAME"
  ttl     = 1
  proxied = false
}

resource "cloudflare_record" "regional_web_a" {
  count   = length(var.micro_ingress_address.ip4) > 0 ? 1 : 0
  zone_id = data.cloudflare_zones.micro.zones[0].id
  name    = "web-${var.region}.${var.regional_subdomain}"
  value   = var.micro_ingress_address.ip4
  type    = "A"
  ttl     = 1
  proxied = false
}

resource "cloudflare_record" "regional_web_aaaa" {
  count   = length(var.micro_ingress_address.ip6) > 0 ? 1 : 0
  zone_id = data.cloudflare_zones.micro.zones[0].id
  name    = "web-${var.region}.${var.regional_subdomain}"
  value   = var.micro_ingress_address.ip6
  type    = "AAAA"
  ttl     = 1
  proxied = false
}

resource "cloudflare_record" "regional_web_cname" {
  count   = length(var.micro_ingress_address.hostname) > 0 ? 1 : 0
  zone_id = data.cloudflare_zones.micro.zones[0].id
  name    = "web-${var.region}.${var.regional_subdomain}"
  value   = var.micro_ingress_address.hostname
  type    = "CNAME"
  ttl     = 1
  proxied = false
}
