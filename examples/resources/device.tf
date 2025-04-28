resource "landb_device" "my_device" {
  name                     = "myhostname.cern.ch"
  zone                     = "ZONE1"
  dhcp_response            = "ALWAYS"
  ipv4_in_dns_and_firewall = true
  ipv6_in_dns_and_firewall = true
  manager_lock             = "NO_LOCK"
  ownership                = "CERN"
  type                     = "COMPUTER"

  manager = {
    type = "EGROUP"
    egroup = {
      name  = "terraform-provider-landb"
      email = "terraform-provider-landb@cern.ch"
    }
  }

  responsible = {
    type = "PERSON"
    person = {
      first_name = "Name"
      last_name  = "LastName"
      email      = "Name@cern.ch"
      username   = "user"
      department = "IT"
      group      = "CD"
    }
  }

  user = {
    type = "PERSON"
    person = {
      first_name = "Name"
      last_name  = "LastName"
      email      = "Name@cern.ch"
      username   = "user"
      department = "IT"
      group      = "CD"
    }
  }

  serial_number    = "SN123456"
  inventory_number = "INV987654"
  tag              = "TEST"
  description      = "Terraform-managed test device"
  parent           = "some-parent-device"
  manufacturer     = "Cisco"
  model            = "ABC123"

  location = {
    building = "0000"
    floor    = "0"
    room     = "0000"
  }

  operating_system {
    family  = "LINUX"
    version = "7"
  }
}
