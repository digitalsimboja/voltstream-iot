variable "project" {
  type    = string
  default = "voltstream"
}

variable "environment" {
  type = string
}

variable "shard_count" {
  description = "Number of shards. Each shard handles 1 MB/s ingest, 2 MB/s read."
  type        = number
  default     = 2
}

variable "retention_hours" {
  description = "Data retention window in hours (24–8760)."
  type        = number
  default     = 24
}
