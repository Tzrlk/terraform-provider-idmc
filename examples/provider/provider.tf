provider "idmc" {
  auth_host = var.auth_host
  auth_user = var.auth_user
  auth_pass = var.auth_pass
}

# Inputs
variable "auth_host" {
  type = string
}
variable "auth_user" {
  type = string
}
variable "auth_pass" {
  type      = string
  sensitive = true
}
