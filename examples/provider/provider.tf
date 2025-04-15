terraform {
  required_providers {
    certmgr = {
      source  = "barnes-c/landb"
      version = "1.0.0"
    }
  }
}

provider "landb" {
  endpoint      = "<YOUR-LANDB-SERVER>"
  client_id     = "<YOUR-CLIENT-id>"
  client_secret = "<YOUR-CLIENT-SECRET>"
  audience      = "<YOUR-AUDIENCE>"
}
