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
