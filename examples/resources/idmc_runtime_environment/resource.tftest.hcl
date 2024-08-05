variables {
  name   = "test_example"
  shared = false
}

run "create" {
}

run "change_name" {
  variables {
    name = "test_example_changed"
  }

  assert {
    error_message = "Resource should not be re-created."
    condition     = idmc_runtime_environment.example.id == run.create.example.id
  }

}

run "change_shared" {
  variables {
    name   = run.change_name.example.name
    shared = true
  }

  assert {
    error_message = "Resource should not be re-created."
    condition     = idmc_runtime_environment.example.id == run.change_name.example.id
  }

}
