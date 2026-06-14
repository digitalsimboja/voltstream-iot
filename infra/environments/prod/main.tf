provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "voltstream"
      Environment = "prod"
      ManagedBy   = "terraform"
    }
  }
}

module "vpc" {
  source             = "../../modules/vpc"
  project            = "voltstream"
  environment        = "prod"
  vpc_cidr           = var.vpc_cidr
  availability_zones = var.availability_zones
}

module "kinesis" {
  source          = "../../modules/kinesis"
  project         = "voltstream"
  environment     = "prod"
  shard_count     = var.kinesis_shard_count
  retention_hours = 48
}

module "iam" {
  source             = "../../modules/iam"
  project            = "voltstream"
  environment        = "prod"
  kinesis_stream_arn = module.kinesis.stream_arn
}

module "lambda" {
  source             = "../../modules/lambda"
  project            = "voltstream"
  environment        = "prod"
  execution_role_arn = module.iam.lambda_exec_role_arn
  kinesis_stream_arn = module.kinesis.stream_arn
}

module "timestream" {
  source                  = "../../modules/timestream"
  project                 = "voltstream"
  environment             = "prod"
  memory_retention_hours  = 24
  magnetic_retention_days = 365
}
