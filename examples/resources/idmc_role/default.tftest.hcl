variables {
  role_name        = "test_role"
  role_description = "Role specifically for testing the IDMC terraform provider"
  role_privileges  = [
    "5ElQYfz8pekjksNRoMY3Kj",
    "???",
  ]
}

run "create" {

  assert {
    error_message = "Resulting name should be as configured."
    condition     = idmc_role.example.name == "test_role"
  }

  assert {
    error_message = "Resulting desc should be as configured."
    condition     = idmc_role.example.description == var.role_description
  }

}

run "remove_privilege" {
  variables {
    role_privileges = [var.role_privileges[0]]
  }

  assert {
    error_message = "Resource should not be re-created."
    condition     = idmc_role.example.id == run.create.example.id
  }

}

run "add_privilege" {
  variables {
    role_privileges = var.role_privileges
  }

  assert {
    error_message = "Resource should not be re-created."
    condition     = idmc_role.example.id == run.remove_privilege.example.id
  }

}

run "force_recreate" {
  variables {
    description = format("%s with a changed description", var.role_description)
  }

  assert {
    error_message = "Resulting description should have changed."
    condition     = idmc_role.example.description != run.add_privilege.example.description
  }

  assert {
    error_message = "Resource should be re-created."
    condition     = idmc_role.example.id != run.add_privilege.example.id
  }

}
