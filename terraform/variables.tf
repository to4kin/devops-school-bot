variable "additional_tags" {
  default = {
    Project    = "DevOps-School-bot"
    Repository = "https://github.com/to4kin/devops-school-bot"
  }
  description = "Additional resource tags"
  type        = map(string)
}

variable "security_class_tags" {
  default = {
    "dtit:sec:InfoSecClass" = "open"
  }
  description = "CCOE Security Class tags"
  type        = map(string)
}

variable "schedule_tags" {
  default = {
    Schedule = "DTIT:no-schedule"
  }
  description = "CCOE Schedule tags"
  type        = map(string)
}

variable "telegram_bot_token" {
  description = "Token to access the Telegram Bot"
  type        = string
  default     = ""
}
