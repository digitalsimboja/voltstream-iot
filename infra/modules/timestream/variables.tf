variable "project" {
  type    = string
  default = "voltstream"
}

variable "environment" {
  type = string
}

variable "memory_retention_hours" {
  description = "Hours to keep data in the memory store (fast queries)."
  type        = number
  default     = 24
}

variable "magnetic_retention_days" {
  description = "Days to retain data in the magnetic store (long-term analytics)."
  type        = number
  default     = 90
}
