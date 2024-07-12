---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "idmc_role Resource - idmc"
subcategory: ""
description: |-
  TODO
---

# idmc_role (Resource)

TODO

## Example Usage

```terraform
resource "idmc_rbac_role" "example" {
  name        = var.role_name
  description = var.role_description
  privileges = [
    var.role_privilege_id,
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the role.
- `privileges` (Set of String) The privileges assigned to the role.

### Optional

- `description` (String) Description of the role.

### Read-Only

- `created_by` (String) User who created the role.
- `created_time` (String) Date and time the role was created.
- `display_description` (String) Description displayed in the user interface.
- `display_name` (String) Role name displayed in the user interface.
- `id` (String) Service generated identifier for the role.
- `org_id` (String) ID of the organization the role belongs to.
- `status` (String) Whether the organization's license to use the role is valid or has expired.
- `system_role` (Boolean) Whether the role is a system-defined role.
- `updated_by` (String) User who last updated the role.
- `updated_time` (String) Date and time the role was last updated.