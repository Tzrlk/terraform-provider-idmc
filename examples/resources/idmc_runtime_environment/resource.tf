resource "idmc_runtime_environment" "example" {
  name   = var.env_name
  shared = false
}
resource "idmc_runtime_environment" "example_shared" {
  name   = var.env_shared_name
  shared = true
}
