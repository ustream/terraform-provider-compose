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

resource "compose_whitelist" "test_ip" {
  ip            = "1.2.3.4/32"
  description   = "whitelist test"
  deployment_id = "bmix-dal-yp-d02b2b39-f7a2-4c99-97f5-3f34d56db2b7"
}
```

### Import:
```
terraform import compose_whitelist.test_ip bmix-dal-yp-d02b2b39-f7a2-4c99-97f5-3f34d56db2b7@1.2.3.4/32
```

### Installation:
```
go get github.com/ustream/terraform-provider-compose

# In ~/.terraformrc

providers {
    compose = "${GOPATH}/bin/terraform-provider-compose"
}
```
