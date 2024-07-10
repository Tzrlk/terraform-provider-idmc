
# Full role list
data "idmc_rbac_roles" "full" {
}

# Full role list with expanded privileges
data "idmc_rbac_roles" "full_expanded" {
  filter {
    expand_privileges = true
  }
}

# Specific role by id
data "idmc_rbac_roles" "role_id" {
  filter {
    role_id = var.role_id
  }
}

# Specific role by id with expanded privileges
data "idmc_rbac_roles" "role_id_expanded" {
  filter {
    role_id           = var.role_id
    expand_privileges = true
  }
}

# Specific role by name
data "idmc_rbac_roles" "role_name" {
  filter {
    role_name = var.role_name
  }
}

# Specific role by name with expanded privileges
data "idmc_rbac_roles" "role_name_expanded" {
  filter {
    role_name         = var.role_name
    expand_privileges = true
  }
}
