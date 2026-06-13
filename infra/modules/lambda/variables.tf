variable "project" {
  type    = string
  default = "voltstream"
}

variable "environment" {
  type = string
}

variable "execution_role_arn" {
  description = "IAM role ARN the Lambda function assumes."
  type        = string
}

variable "kinesis_stream_arn" {
  description = "Kinesis stream ARN to consume from."
  type        = string
}

variable "lambda_zip_path" {
  description = "Path to the compiled Lambda zip artifact."
  type        = string
  default     = "../../../lambda/anomaly-detector/anomaly-detector.zip"
}
