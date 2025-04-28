resource "landb_set_attach" "set_attach" {
  set_id  = landb_set.set.id
  host_ip = local.ipv4
}
