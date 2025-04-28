resource "landb_set_attach" "example" {
  set_name    = landb_set.example.name
  name        = "attachment-01"
  ipv4        = "192.168.100.100"
  ipv6        = "2001:db8::100"
  description = "example interface"
}