variable "project" {
  type    = string
  default = "voltstream"
}

variable "environment" {
  type = string
}

variable "kinesis_stream_arn" {
  description = "ARN of the Kinesis stream the Lambda role needs read access to."
  type        = string
}
