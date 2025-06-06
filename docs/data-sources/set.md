---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "landb_set Data Source - landb"
subcategory: ""
description: |-
  Data source for retrieving an existing set
---

# landb_set (Data Source)

Data source for retrieving an existing set



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)

### Optional

- `responsible` (Attributes) Responsible entity for the set (see [below for nested schema](#nestedatt--responsible))

### Read-Only

- `description` (String)
- `id` (String) The ID of this resource.
- `last_updated` (String)
- `network_domain` (String)
- `project_url` (String)
- `receive_notifications` (Boolean)
- `type` (String)
- `version` (Number)

<a id="nestedatt--responsible"></a>
### Nested Schema for `responsible`

Optional:

- `egroup` (Attributes) Details if type == EGROUP (see [below for nested schema](#nestedatt--responsible--egroup))
- `person` (Attributes) Details if type == PERSON (see [below for nested schema](#nestedatt--responsible--person))
- `reserved` (Attributes) Details if type == RESERVED (see [below for nested schema](#nestedatt--responsible--reserved))
- `type` (String) One of PERSON, EGROUP, or RESERVED

<a id="nestedatt--responsible--egroup"></a>
### Nested Schema for `responsible.egroup`

Optional:

- `email` (String)
- `name` (String)


<a id="nestedatt--responsible--person"></a>
### Nested Schema for `responsible.person`

Optional:

- `department` (String)
- `email` (String)
- `first_name` (String)
- `group` (String)
- `last_name` (String)
- `username` (String)


<a id="nestedatt--responsible--reserved"></a>
### Nested Schema for `responsible.reserved`

Optional:

- `first_name` (String)
- `last_name` (String)
