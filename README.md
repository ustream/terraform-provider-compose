# terraform-provider-compose
Terraform provider for IBM Cloud (formerly Bluemix) based Compose databases

### Based on:
- https://www.compose.com/articles/p/078ca17a-ffe2-43cb-bd35-eb69593d2bff/
- https://apidocs.compose.com/v1.0/reference#the-compose-api

### Usage:
```hcl
provider "compose" {
  bluemix_api_key = "${var.ibm_bmx_api_key}" # if not set, the BM_API_KEY environment variable is read
  region          = "${var.ibm_bmx_region}"  # if not set, the BM_REGION environment variable is read (default: us-south)
}

resource "compose_whitelist" "office_ip" {
  ip            = "195.56.66.6/32"
  description   = "Allow connection from Ustream office only"
  deployment_id = "bmix-eude-yp-dacd993c-8989-47c8-96a5-01a8ea4a99f4"
}
```

### Installation:
```
go get github.com/ustream/terraform-provider-compose

# In ~/.terraformrc

providers {
    compose = "${GOPATH}/bin/terraform-provider-compose"
}
```
