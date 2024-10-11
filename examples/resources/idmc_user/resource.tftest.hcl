variables {
  user_name        = "test_user"
  user_first_name  = "Test"
  user_last_name   = "User"
  user_title       = "Lord"
  user_roles       = [
    { name = "Service Consumer" }
    { name = "Designer" }
  ]
}

run "create" {

  assert {
    error_message = "Resulting name shoudl be as configured."
    condition     = idmc_user.example.name == "test_user"
  }

  assert {
    error_message = "Resulting description should be as configured."
    condition     = idmc_user.example.description == var.user_description
  }

}

run "remove_role" {
}

run "add_role" {
}
