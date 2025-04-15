resource "landb_set" "my_set" {
  name                  = "TF-TEST-SET-12345"
  type                  = "INTERDOMAIN"
  network_domain        = "IT-COMPUTING-NETWORK"
  description           = "Terraform-managed test set"
  project_url           = "https://example.com"
  receive_notifications = true

  responsible {
    type = "PERSON"
    person {
      first_name = "Christopher"
      last_name  = "Barnes"
      email      = "christopher.barnes@cern.ch"
      username   = "chbarnes"
      department = "IT"
      group      = "CD"
    }
  }
}
