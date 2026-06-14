terraform {
  required_version = ">= 1.7"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  # TODO: configure remote state backend (S3 + DynamoDB lock)
  # backend "s3" {
  #   bucket         = "voltstream-tfstate"
  #   key            = "dev/terraform.tfstate"
  #   region         = var.aws_region
  #   dynamodb_table = "voltstream-tfstate-lock"
  # }
}