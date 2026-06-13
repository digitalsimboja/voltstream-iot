variable "project" {
  description = "Project name used in resource tags and naming."
  type        = string
  default     = "voltstream"
}

variable "environment" {
  description = "Deployment environment (dev, prod)."
  type        = string
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC."
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "List of AZs to create subnets in."
  type        = list(string)
  default     = ["eu-north-1a", "eu-north-1b"]
}
