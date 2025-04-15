resource "landb_device" "my_device" {
  name                     = "myhostname.cern.ch"
  zone                     = "ZONE1"
  dhcp_response            = "ALWAYS"
  ipv4_in_dns_and_firewall = true
  ipv6_in_dns_and_firewall = true
  manager_lock             = "NO_LOCK"
  ownership                = "CERN"
  type                     = "COMPUTER"

  manager {
    type = "EGROUP"
    egroup {
      name  = "ai-playground"
      email = "ai-playground-admins@cern.ch"
    }
  }

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

  user {
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

  serial_number    = "SN123456"
  inventory_number = "INV987654"
  tag              = "TEST"
  description      = "Terraform-managed test device"
  parent           = "some-parent-device"
  manufacturer     = "Cisco"
  model            = "ABC123"

  location {
    building = "31"
    floor    = "1"
    room     = "006"
  }

  operating_system {
    family  = "LINUX"
    version = "7"
  }
}
