variables {
  env_name   = "test_example"
  env_shared = false
}

run "create" {
}

run "change_name" {
  variables {
    env_name = "test_example_changed"
  }

  assert {
    error_message = "Resource should not be re-created."
    condition     = idmc_runtime_environment.example.id == run.create.example.id
  }

}

run "change_shared" {
  variables {
    env_name = run.change_name.example.name
    shared   = true
  }

  assert {
    error_message = "Resource should not be re-created."
    condition     = idmc_runtime_environment.example.id == run.change_name.example.id
  }

}
