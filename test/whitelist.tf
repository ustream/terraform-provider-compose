provider "compose" {
  region = "eu-de"
}

resource "compose_whitelist" "office_ip" {
  ip            = "195.56.66.6/32"
  description   = "Allow connection from Ustream office only"
  deployment_id = "bmix-eude-yp-dacd993c-8989-47c8-96a5-01a8ea4a99f4"
}
