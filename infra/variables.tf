variable "cp_host" {
  type = string
}

variable "cp_username" {
  type      = string
  sensitive = true
}

variable "cp_password" {
  type      = string
  sensitive = true
}
