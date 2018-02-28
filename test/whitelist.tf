variable "test_deployment_id" {}
variable "region" {}

provider "compose" {
  region = "${var.region}"
}

resource "compose_whitelist" "test_ip" {
  ip            = "1.2.3.4/32"
  description   = "test whitelist"
  deployment_id = "${var.test_deployment_id}"
}
