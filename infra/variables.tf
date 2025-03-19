locals {
  default_ns_labels = {
    monitoring = "enabled"
  }
}

variable "telegram_api_key" {
  type      = string
  sensitive = true
}

variable "telegram_chat_id" {
  type = number
}
