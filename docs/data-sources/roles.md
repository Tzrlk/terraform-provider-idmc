---
# Had to mess with description rendering because it was failing with a weird error.
page_title: "idmc_roles Data Source - idmc"
subcategory: ""
description: |-
  https://docs.informatica.com/integration-cloud/data-integration/current-version/rest-api-reference/platform-rest-api-version-3-resources/roles/getting-role-details.html
---

# idmc_roles (Data Source)

https://docs.informatica.com/integration-cloud/data-integration/current-version/rest-api-reference/platform-rest-api-version-3-resources/roles/getting-role-details.html

## Example Usage

```terraform
# Full role list
data "idmc_roles" "full" {
}

# Full role list with expanded privileges
data "idmc_roles" "full_expanded" {
  filter {
    expand_privileges = true
  }
}

# Specific role by id
data "idmc_roles" "role_id" {
  filter {
    role_id = var.role_id
  }
}

# Specific role by id with expanded privileges
data "idmc_roles" "role_id_expanded" {
  filter {
    role_id           = var.role_id
    expand_privileges = true
  }
}

# Specific role by name
data "idmc_roles" "role_name" {
  filter {
    role_name = var.role_name
  }
}

# Specific role by name with expanded privileges
data "idmc_roles" "role_name_expanded" {
  filter {
    role_name         = var.role_name
    expand_privileges = true
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block, Optional) Allows for results to be narrowed. (see [below for nested schema](#nestedblock--filter))

### Read-Only

- `roles` (Attributes Map) The query results (see [below for nested schema](#nestedatt--roles))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Optional:

- `expand_privileges` (Boolean) Returns the privileges associated with the role specified in the query filter.
- `role_id` (String) Unique identifier for the role.
- `role_name` (String) Name of the role.


<a id="nestedatt--roles"></a>
### Nested Schema for `roles`

Read-Only:

- `created_by` (String) User who created the role.
- `created_time` (String) Date and time the role was created.
- `description` (String) Description of the role.
- `display_description` (String) Description displayed in the user interface.
- `display_name` (String) Role name displayed in the user interface.
- `name` (String) Name of the role.
- `org_id` (String) ID of the organization the role belongs to.
- `privileges` (Attributes Map) The privileges assigned to the role. (see [below for nested schema](#nestedatt--roles--privileges))
- `status` (String) Whether the organization's license to use the role is valid or has expired.
- `system_role` (Boolean) Whether the role is a system-defined role.
- `updated_by` (String) User who last updated the role.
- `updated_time` (String) Date and time the role was last updated.

<a id="nestedatt--roles--privileges"></a>
### Nested Schema for `roles.privileges`

Read-Only:

- `description` (String) Description of the privilege.
- `name` (String) Name of the privilege.
- `service` (String) Service the privilege applies to.
- `status` (String) Status of the privilege (Enabled/Disabled).
