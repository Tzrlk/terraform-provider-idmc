resource "idmc_user" "example" {

  # Required
  name       = var.user_name
  password   = var.user_password
  first_name = var.user_first_name
  last_name  = var.user_last_name
  title      = var.user_title
  phone      = var.user_phone
  roles      = var.user_roles

  # Optional
  org_id                = var.user_org_id
  salesforce_username   = var.user_salesforce_username
  email                 = var.user_email
  description           = var.user_description
  timezone              = var.user_timezone
  security_question     = var.user_security_question
  security_answer       = var.user_security_answer
  force_change_password = var.user_force_change_password

}

# Inputs
variable "user_name" {
  type = string
}
variable "user_password" {
  type = string
  sensitive = true
}
variable "user_first_name" {
  type = string
}
variable "user_last_name" {
  type = string
}
variable "user_title" {
  type = string
}
variable "user_phone" {
  type = string
}
variable "user_roles" {
  type = list(object{
    name        = string
    description = optional(string)
  })
}
variable "user_org_id" {
  type    = string
  default = null
}
variable "user_salesforce_username" {
  type    = string
  default = null
}
variable "user_description" {
  type    = string
  default = null
}
variable "user_timezone" {
  type    = string
  default = null
}
variable "user_security_question" {
  type    = string
  default = null
}
variable "user_security_answer" {
  type    = string
  default = null
}
variable "user_force_change_password" {
  type    = bool
  default = null
}
